package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	aw "github.com/deanishe/awgo"
	bb "github.com/rwilgaard/bitbucket-go-api"
	"go.deanishe.net/fuzzy"
)

type workflowConfig struct {
    URL      string `env:"bitbucket_url"`
    Username string `env:"username"`
    APIToken string `env:"apitoken"`
}

var (
    wf            *aw.Workflow
    authFlag      string
    cacheFlag     bool
    commitFlag    bool
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
    )
    flag.StringVar(&authFlag, "auth", "", "authentication")
    flag.BoolVar(&cacheFlag, "cache", false, "cache repositories")
    flag.BoolVar(&commitFlag, "commits", false, "show commits for repository")
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
        log.Println(repos.NextPageStart)
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

func run() {
    wf.Args()
    flag.Parse()
    query := flag.Arg(0)

    cfg := &workflowConfig{}
    if err := wf.Config.To(cfg); err != nil {
        wf.FatalError(err)
    }

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
            Icon(aw.IconHome).
            Valid(true)

        for _, c := range commits.Values {
            t := time.UnixMilli(c.CommitterTimestamp).Format("02-01-2006 15:04")
            wf.NewItem(c.Message).
                Subtitle(fmt.Sprintf("%s  |  %s  |  %s", c.DisplayID, c.Committer.Name, t)).
                Var("message", c.Message).
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
        wf.Rerun(0.3)
        if !wf.IsRunning("cache") {
            cmd := exec.Command(os.Args[0], "-cache")
            if err := wf.RunInBackground("cache", cmd); err != nil {
                wf.FatalError(err)
            }
        } else {
            log.Printf("cache job already running.")
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
                Var("query", query).
                Valid(true)

            it.NewModifier(aw.ModOpt).
                Subtitle("Show commits").
                Arg("commits").
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
