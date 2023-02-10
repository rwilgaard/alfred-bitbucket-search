package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"time"

	aw "github.com/deanishe/awgo"
	bb "github.com/rwilgaard/bitbucket-go-api"
)

type workflowConfig struct {
    URL      string `env:"URL"`
    Username string `env:"USERNAME"`
    APIToken string
}

var (
    wf          *aw.Workflow
    authFlag    string
    cacheFlag   bool
    cacheName   = "repositories.json"
    maxCacheAge = 1 * time.Hour
)

func init() {
    wf = aw.New()
    flag.StringVar(&authFlag, "auth", "", "authentication")
    flag.BoolVar(&cacheFlag, "cache", false, "cache repositories")
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

func run() {
    wf.Args()
    flag.Parse()
    query := flag.Arg(0)

    if authFlag != "" {
        if err := wf.Keychain.Set(authFlag, os.Getenv("ALFRED_AUTHCONFIG_PASSWORD")); err != nil {
            wf.FatalError(err)
        }
        return
    }

    authType := "API Token"
    token, err := wf.Keychain.Get(authType)
    if err != nil {
        log.Println(err)
        wf.NewItem("Credentials not configured...").
            Subtitle("Press ⏎ to configure").
            Icon(aw.IconInfo).
            Arg("auth").
            Var("auth_type", authType).
            Valid(true)
        wf.SendFeedback()
        return
    }

    cfg := workflowConfig{
        APIToken: token,
    }


    if cacheFlag {
        wf.Configure(aw.TextErrors(true))
        log.Printf("[cache] fetching repositories...")

        api, err := bb.NewAPI(cfg.URL, cfg.Username, cfg.APIToken)
        if err != nil {
            wf.FatalError(err)
        }

        repos, err := getAllRepositories(api)
        if err != nil {
            wf.FatalError(err)
        }

        if err := wf.Cache.StoreJSON(cacheName, repos); err != nil {
            wf.FatalError(err)
        }

        log.Printf("[cache] repositories fetched")
        return
    }

    var repos []*bb.RepositoryList
    if wf.Cache.Exists(cacheName) {
        if err := wf.Cache.LoadJSON(cacheName, &repos); err != nil {
            wf.FatalError(err)
        }
    }

    if wf.Cache.Expired(cacheName, maxCacheAge) {
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
            wf.NewItem(repo.Name).
                Subtitle(repo.Project.Name)
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
