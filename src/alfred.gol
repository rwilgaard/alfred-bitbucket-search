package main

import (
    "fmt"
    "os"
    "regexp"
    "slices"
    "strings"
    "time"

    aw "github.com/deanishe/awgo"
    "github.com/ncruces/zenity"
    bb "github.com/rwilgaard/bitbucket-go-api"
)

type parsedQuery struct {
    Text     string
    Projects []string
}

type magicAuth struct {
    wf *aw.Workflow
}

func (a magicAuth) Keyword() string     { return "clearauth" }
func (a magicAuth) Description() string { return "Clear credentials for Bitbucket." }
func (a magicAuth) RunText() string     { return "Credentials cleared!" }
func (a magicAuth) Run() error          { return clearAuth() }

func parseQuery(query string) *parsedQuery {
    q := new(parsedQuery)
    projectRegex := regexp.MustCompile(`^@\w+`)

    for _, w := range strings.Split(query, " ") {
        switch {
        case projectRegex.MatchString(w):
            q.Projects = append(q.Projects, w[1:])
        default:
            q.Text = q.Text + w + " "
        }
    }

    return q
}

func autocomplete(query string) string {
    for _, w := range strings.Split(query, " ") {
        switch w {
        case "@":
            return "projects"
        }
    }
    return ""
}

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

func runProjects() {
    var projects *bb.ProjectList
    if err := wf.Cache.LoadJSON(projectCacheName, &projects); err != nil {
        wf.FatalError(err)
    }

    prevQuery, _ := wf.Config.Env.Lookup("prev_query")

    for _, p := range projects.Values {
        i := wf.NewItem(p.Key).
            Match(fmt.Sprintf("%s %s", p.Key, p.Name)).
            UID(p.Key).
            Subtitle(p.Name).
            Arg(prevQuery+p.Key+" ").
            Var("project", p.Key).
            Valid(true)

        i.NewModifier(aw.ModCmd).
            Subtitle("Cancel").
            Arg("cancel")

    }
}

func runSearch(query *parsedQuery) {
    var repos []*bb.RepositoryList
    if err := wf.Cache.LoadJSON(repoCacheName, &repos); err != nil {
        wf.FatalError(err)
    }

    for _, list := range repos {
        for _, repo := range list.Values {
            if len(query.Projects) > 0 && !slices.Contains(query.Projects, repo.Project.Key) {
                continue
            }

            i := wf.NewItem(repo.Name).
                Subtitle(repo.Project.Name).
                Match(fmt.Sprintf("%s %s %s %s", repo.Name, repo.Slug, repo.Project.Name, repo.Project.Key)).
                UID(fmt.Sprintf("%s/%s", repo.Project.Key, repo.Slug)).
                Var("projectKey", repo.Project.Key).
                Var("repoSlug", repo.Slug).
                Var("link", repo.Links["self"][0].Href).
                Var("lastQuery", opts.Query).
                Valid(true)

            i.NewModifier(aw.ModCmd).
                Subtitle("Show Commits").
                Arg("commits").
                Valid(true)

            i.NewModifier(aw.ModCtrl).
                Subtitle("Show Tags").
                Arg("tags").
                Valid(true)

            i.NewModifier(aw.ModOpt).
                Subtitle("Show Pull Requests").
                Arg("pullrequests").
                Valid(true)

            i.NewModifier(aw.ModShift).
                Subtitle("Show Branches").
                Arg("branches").
                Valid(true)

            sshIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "ssh" })
            httpIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "http" })
            i.NewModifier(aw.ModOpt, aw.ModShift).
                Subtitle("Copy HTTP clone URL").
                Arg("copy").
                Var("copy_value", repo.Links["clone"][httpIdx].Href).
                Valid(true)

            i.NewModifier(aw.ModCmd, aw.ModShift).
                Subtitle("Copy SSH clone URL").
                Arg("copy").
                Var("copy_value", repo.Links["clone"][sshIdx].Href).
                Valid(true)
        }
    }
}

func runAuth() {
    _, pwd, err := zenity.Password(
        zenity.Title(fmt.Sprintf("Enter API Token for %s", cfg.Username)),
    )
    if err != nil {
        wf.FatalError(err)
    }

    api, err := bb.NewAPI(cfg.URL, cfg.Username, pwd)
    if err != nil {
        wf.FatalError(err)
    }

    sc, err := testAuthentication(api)
    if err != nil {
        zerr := zenity.Error(
            fmt.Sprintf("Error authenticating: HTTP %d", sc),
            zenity.ErrorIcon,
        )
        if zerr != nil {
            wf.FatalError(err)
        }
        wf.FatalError(err)
    }

    if err := wf.Keychain.Set(keychainAccount, pwd); err != nil {
        wf.FatalError(err)
    }
}

func clearAuth() error {
    if err := wf.Keychain.Delete(keychainAccount); err != nil {
        return err
    }
    return nil
}
