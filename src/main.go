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
)

type workflowConfig struct {
    URL      string `env:"bitbucket_url"`
    Username string `env:"username"`
    APIToken string `env:"apitoken"`
}

const (
    repo          = "rwilgaard/alfred-bitbucket-search"
    updateJobName = "checkForUpdates"
)

var (
    wf            *aw.Workflow
    authFlag      string
    cacheFlag     bool
    updateFlag    bool
    commitFlag    bool
    tagFlag       bool
    prFlag        bool
    repoCacheName = "repositories.json"
    maxCacheAge   = 1 * time.Hour
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
    flag.BoolVar(&cacheFlag, "cache", false, "cache repositories")
    flag.BoolVar(&commitFlag, "commits", false, "show commits for repository")
    flag.BoolVar(&tagFlag, "tags", false, "show tags for repository")
    flag.BoolVar(&prFlag, "pullrequests", false, "show pull requests for repository")
    flag.BoolVar(&updateFlag, "update", false, "check for updates")
}

func getAllRepositories(api *bb.API) ([]*bb.RepositoryList, error) {
    query := bb.RepositoriesQuery{
        Limit: 9999,
    }

    repos, err := api.GetRepositories(query)
    if err != nil {
        return nil, err
    }

    var results []*bb.RepositoryList
    results = append(results, repos)

    for !repos.IsLastPage {
        query := bb.RepositoriesQuery{
            Limit: 9999,
            Start: int(repos.NextPageStart),
        }
        repos, err = api.GetRepositories(query)
        if err != nil {
            return nil, err
        }
        results = append(results, repos)
    }

    return results, nil
}

func getCommits(api *bb.API, projectKey string, repoSlug string) (*bb.CommitList, error) {
    query := bb.CommitsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    commits, err := api.GetCommits(query)
    if err != nil {
        return nil, err
    }

    return commits, nil
}

func getTags(api *bb.API, projectKey string, repoSlug string) (*bb.TagList, error) {
    query := bb.TagsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
        OrderBy:        "MODIFICATION",
    }

    tags, err := api.GetTags(query)
    if err != nil {
        return nil, err
    }

    return tags, nil
}

func getPullRequests(api *bb.API, projectKey string, repoSlug string) (*bb.PullRequestList, error) {
    query := bb.PullRequestsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    pr, err := api.GetPullRequests(query)
    if err != nil {
        return nil, err
    }

    return pr, nil
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
            wf.NewItem(c.Message).
                Subtitle(fmt.Sprintf("%s  |  %s  |  %s", c.DisplayID, c.Committer.Name, t)).
                Icon(&icon).
                Var("message", c.Message).
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
                Valid(true)
        }

        wf.SendFeedback()
        return
    }

    if cacheFlag {
        wf.Configure(aw.TextErrors(true))
        log.Printf("[cache] fetching repositories...")

        repos, err := getAllRepositories(api)
        if err != nil {
            wf.FatalError(err)
        }

        if err := wf.Cache.StoreJSON(repoCacheName, repos); err != nil {
            wf.FatalError(err)
        }

        log.Printf("[cache] repositories fetched")
        return
    }

    var repos []*bb.RepositoryList
    if wf.Cache.Exists(repoCacheName) {
        if err := wf.Cache.LoadJSON(repoCacheName, &repos); err != nil {
            wf.FatalError(err)
        }
    }

    if wf.Cache.Expired(repoCacheName, maxCacheAge) {
        wf.Rerun(2)
        if !wf.IsRunning("cache") {
            log.Printf("[cache] starting cache job")
            cmd := exec.Command(os.Args[0], "-cache")
            if err := wf.RunInBackground("cache", cmd); err != nil {
                wf.FatalError(err)
            }
        } else {
            log.Printf("[cache] cache job already running.")
        }

        if len(repos) == 0 {
            wf.NewItem("Fetching repositories...").
                Icon(aw.IconInfo)
            wf.SendFeedback()
            return
        }
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
