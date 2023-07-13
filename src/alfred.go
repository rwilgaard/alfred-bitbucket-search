package main

import (
    "fmt"
    "os"
    "time"

    aw "github.com/deanishe/awgo"
    bb "github.com/rwilgaard/bitbucket-go-api"
    "golang.org/x/exp/slices"
)

func runCommits(api *bb.API) {
    wf.Configure(aw.SuppressUIDs(true))
    repoSlug := os.Getenv("repoSlug")
    projectKey := os.Getenv("projectKey")
    commits, err := getCommits(api, projectKey, repoSlug)
    if err != nil {
        wf.FatalError(err)
    }

    wf.NewItem("Go back").
        Icon(backIcon).
        Arg("go-back").
        Valid(true)

    icon := aw.Icon{Value: fmt.Sprintf("%s/icons/commit.png", wf.Dir())}
    for _, c := range commits.Values {
        t := time.UnixMilli(c.CommitterTimestamp).Format("02-01-2006 15:04")
        i := wf.NewItem(c.Message).
            Subtitle(fmt.Sprintf("%s  |  %s  |  %s", c.DisplayID, c.Committer.Name, t)).
            Icon(&icon).
            Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/commits/%s", cfg.URL, projectKey, repoSlug, c.ID)).
            Valid(true)

        i.NewModifier(aw.ModCmd).
            Subtitle("Show full commit message.").
            Arg("commit").
            Var("message", c.Message).
            Valid(true)
    }
}

func runBranches(api *bb.API) {
    wf.Configure(aw.SuppressUIDs(true))
    repoSlug := os.Getenv("repoSlug")
    projectKey := os.Getenv("projectKey")
    branches, err := getBranches(api, projectKey, repoSlug)
    if err != nil {
        wf.FatalError(err)
    }

    wf.NewItem("Go back").
        Icon(backIcon).
        Arg("go-back").
        Valid(true)

    icon := aw.Icon{Value: fmt.Sprintf("%s/icons/branch.png", wf.Dir())}
    for _, b := range branches.Values {
        wf.NewItem(b.DisplayID).
            Subtitle(fmt.Sprintf("Latest commit: %s", b.LatestCommit[0:10])).
            Icon(&icon).
            Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/browse?at=%s", cfg.URL, projectKey, repoSlug, b.ID)).
            Valid(true)
    }
}

func runTags(api *bb.API) {
    wf.Configure(aw.SuppressUIDs(true))
    repoSlug := os.Getenv("repoSlug")
    projectKey := os.Getenv("projectKey")
    tags, err := getTags(api, projectKey, repoSlug)
    if err != nil {
        wf.FatalError(err)
    }

    wf.NewItem("Go back").
        Icon(backIcon).
        Arg("go-back").
        Valid(true)

    icon := aw.Icon{Value: fmt.Sprintf("%s/icons/tag.png", wf.Dir())}
    for _, t := range tags.Values {
        wf.NewItem(t.DisplayID).
            Subtitle(fmt.Sprintf("Commit: %s", t.LatestCommit[0:10])).
            Icon(&icon).
            Valid(true)
    }
}

func runPullRequests(api *bb.API) {
    wf.Configure(aw.SuppressUIDs(true))
    repoSlug := os.Getenv("repoSlug")
    projectKey := os.Getenv("projectKey")
    tags, err := getPullRequests(api, projectKey, repoSlug)
    if err != nil {
        wf.FatalError(err)
    }

    wf.NewItem("Go back").
        Icon(backIcon).
        Arg("go-back").
        Valid(true)

    icon := aw.Icon{Value: fmt.Sprintf("%s/icons/pull-request.png", wf.Dir())}
    for _, p := range tags.Values {
        wf.NewItem(p.Title).
            Subtitle(fmt.Sprintf("%s ➔ %s", p.FromRef.DisplayID, p.ToRef.DisplayID)).
            Icon(&icon).
            Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/overview", cfg.URL, projectKey, repoSlug, p.ID)).
            Valid(true)
    }
}

func runSearch(api *bb.API) {
    var repos []*bb.RepositoryList
    if wf.Cache.Exists(repoCacheName) {
        if err := wf.Cache.LoadJSON(repoCacheName, &repos); err != nil {
            wf.FatalError(err)
        }
    }

    maxCacheAge := cfg.CacheAge * int(time.Minute)
    if wf.Cache.Expired(repoCacheName, time.Duration(maxCacheAge)) {
        if err := cacheRepositories(api); err != nil {
            wf.FatalError(err)
        }
        wf.Rerun(0.3)
    }

    for _, list := range repos {
        for _, repo := range list.Values {
            it := wf.NewItem(repo.Name).
                Subtitle(repo.Project.Name).
                Match(fmt.Sprintf("%s %s %s %s", repo.Name, repo.Slug, repo.Project.Name, repo.Project.Key)).
                Var("projectKey", repo.Project.Key).
                Var("repoSlug", repo.Slug).
                Var("link", repo.Links["self"][0].Href).
                Var("lastQuery", opts.Query).
                Valid(true)

            it.NewModifier(aw.ModCmd).
                Subtitle("Show Commits").
                Arg("commits").
                Valid(true)

            it.NewModifier(aw.ModCtrl).
                Subtitle("Show Tags").
                Arg("tags").
                Valid(true)

            it.NewModifier(aw.ModOpt).
                Subtitle("Show Pull Requests").
                Arg("pullrequests").
                Valid(true)

            it.NewModifier(aw.ModShift).
                Subtitle("Show Branches").
                Arg("branches").
                Valid(true)

            sshIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "ssh" })
            httpIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "http" })
            it.NewModifier(aw.ModOpt, aw.ModShift).
                Subtitle("Copy HTTP clone URL").
                Arg("copy").
                Var("copy_value", repo.Links["clone"][httpIdx].Href).
                Valid(true)

            it.NewModifier(aw.ModCmd, aw.ModShift).
                Subtitle("Copy SSH clone URL").
                Arg("copy").
                Var("copy_value", repo.Links["clone"][sshIdx].Href).
                Valid(true)
        }
    }
}
