package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "time"

    aw "github.com/deanishe/awgo"
    "github.com/deanishe/awgo/update"
    bb "github.com/rwilgaard/bitbucket-go-api"
    "go.deanishe.net/fuzzy"
    "golang.org/x/exp/slices"
)

type workflowConfig struct {
    URL      string `env:"bitbucket_url"`
    Username string `env:"username"`
    APIToken string `env:"apitoken"`
    CacheAge int    `env:"cache_age"`
}

const (
    repo          = "rwilgaard/alfred-bitbucket-search"
    updateJobName = "checkForUpdates"
)

var (
    wf            *aw.Workflow
    authFlag      string
    updateFlag    bool
    commitFlag    bool
    branchFlag    bool
    tagFlag       bool
    prFlag        bool
    repoCacheName = "repositories.json"
)

func init() {
    sopts := []fuzzy.Option{
        fuzzy.AdjacencyBonus(10.0),
        fuzzy.LeadingLetterPenalty(-0.1),
        fuzzy.MaxLeadingLetterPenalty(-3.0),
        fuzzy.UnmatchedLetterPenalty(-0.5),
    }
    wf = aw.New(
        aw.SortOptions(sopts...),
        update.GitHub(repo),
    )
    flag.StringVar(&authFlag, "auth", "", "authentication")
    flag.BoolVar(&commitFlag, "commits", false, "show commits for repository")
    flag.BoolVar(&branchFlag, "branches", false, "show branches for repository")
    flag.BoolVar(&tagFlag, "tags", false, "show tags for repository")
    flag.BoolVar(&prFlag, "pullrequests", false, "show pull requests for repository")
    flag.BoolVar(&updateFlag, "update", false, "check for updates")
}

func cacheRepositories(api *bb.API) error {
    log.Printf("[cache] fetching repositories...")

    repos, err := getAllRepositories(api)
    if err != nil {
        return err
    }

    if err := wf.Cache.StoreJSON(repoCacheName, repos); err != nil {
        return err
    }

    log.Printf("[cache] repositories fetched")
    return nil
}

func run() {
    wf.Args()
    flag.Parse()
    query := flag.Arg(0)

    if updateFlag {
        wf.Configure(aw.TextErrors(true))
        log.Println("Checking for updates...")
        if err := wf.CheckForUpdate(); err != nil {
            wf.FatalError(err)
        }
        return
    }

    if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
        log.Println("Running update check in background...")
        cmd := exec.Command(os.Args[0], "-update")
        if err := wf.RunInBackground(updateJobName, cmd); err != nil {
            log.Printf("Error starting update check: %s", err)
        }
    }

    if wf.UpdateAvailable() {
        wf.Configure(aw.SuppressUIDs(true))
        wf.NewItem("Update Available!").
            Subtitle("Press ⏎ to install").
            Autocomplete("workflow:update").
            Valid(false).
            Icon(aw.IconInfo)
    }

    cfg := &workflowConfig{}
    if err := wf.Config.To(cfg); err != nil {
        wf.FatalError(err)
    }

    backIcon := aw.Icon{Value: fmt.Sprintf("%s/icons/go-back.png", wf.Dir())}

    api, err := bb.NewAPI(cfg.URL, cfg.Username, cfg.APIToken)
    if err != nil {
        wf.FatalError(err)
    }

    if commitFlag {
        wf.Configure(aw.SuppressUIDs(true))
        repoSlug := os.Getenv("repoSlug")
        projectKey := os.Getenv("projectKey")
        commits, err := getCommits(api, projectKey, repoSlug)
        if err != nil {
            wf.FatalError(err)
        }

        wf.NewItem("Go back").
            Icon(&backIcon).
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

        wf.SendFeedback()
        return
    }

    if branchFlag {
        wf.Configure(aw.SuppressUIDs(true))
        repoSlug := os.Getenv("repoSlug")
        projectKey := os.Getenv("projectKey")
        branches, err := getBranches(api, projectKey, repoSlug)
        if err != nil {
            wf.FatalError(err)
        }

        wf.NewItem("Go back").
            Icon(&backIcon).
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

        wf.SendFeedback()
        return
    }

    if tagFlag {
        wf.Configure(aw.SuppressUIDs(true))
        repoSlug := os.Getenv("repoSlug")
        projectKey := os.Getenv("projectKey")
        tags, err := getTags(api, projectKey, repoSlug)
        if err != nil {
            wf.FatalError(err)
        }

        wf.NewItem("Go back").
            Icon(&backIcon).
            Arg("go-back").
            Valid(true)

        icon := aw.Icon{Value: fmt.Sprintf("%s/icons/tag.png", wf.Dir())}
        for _, t := range tags.Values {
            wf.NewItem(t.DisplayID).
                Subtitle(fmt.Sprintf("Commit: %s", t.LatestCommit[0:10])).
                Icon(&icon).
                Valid(true)
        }

        wf.SendFeedback()
        return
    }

    if prFlag {
        wf.Configure(aw.SuppressUIDs(true))
        repoSlug := os.Getenv("repoSlug")
        projectKey := os.Getenv("projectKey")
        tags, err := getPullRequests(api, projectKey, repoSlug)
        if err != nil {
            wf.FatalError(err)
        }

        wf.NewItem("Go back").
            Icon(&backIcon).
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

        wf.SendFeedback()
        return
    }

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
                Var("lastQuery", query).
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

    wf.Filter(query)

    if wf.IsEmpty() {
        wf.NewItem("No results found...").
            Subtitle("Try a different query?").
            Icon(aw.IconInfo)
    }
    wf.SendFeedback()
}

func main() {
    wf.Run(run)
}
