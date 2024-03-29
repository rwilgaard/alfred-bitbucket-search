package main

import (
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
    CacheAge int    `env:"cache_age"`
    APIToken string
}

const (
    repo             = "rwilgaard/alfred-bitbucket-search"
    updateJobName    = "checkForUpdates"
    repoCacheName    = "repositories.json"
    projectCacheName = "projects.json"
    keychainAccount  = "alfred-bitbucket-search"
)

var (
    cfg *workflowConfig
    wf  *aw.Workflow
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
        aw.AddMagic(magicAuth{wf}),
    )
}

func refreshCache(api *bb.API) error {
    log.Printf("[cache] fetching repositories...")
    repos, err := getAllRepositories(api)
    if err != nil {
        return err
    }

    if err := wf.Cache.StoreJSON(repoCacheName, repos); err != nil {
        return err
    }
    log.Printf("[cache] repositories fetched")

    log.Printf("[cache] fetching projects...")
    projects, err := getAllProjects(api)
    if err != nil {
        return err
    }

    if err := wf.Cache.StoreJSON(projectCacheName, projects); err != nil {
        return err
    }
    log.Printf("[cache] projects fetched")

    return nil
}

func run() {
    wf.Args()
    if err := cli.Parse(wf.Args()); err != nil {
        wf.FatalError(err)
    }
    opts.Query = cli.Arg(0)

    if opts.Update {
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

    cfg = &workflowConfig{}
    if err := wf.Config.To(cfg); err != nil {
        wf.FatalError(err)
    }

    if opts.Auth {
        runAuth()
    }

    token, err := wf.Keychain.Get(keychainAccount)
    if err != nil {
        wf.NewItem("You're not logged in.").
            Subtitle("Press ⏎ to authenticate").
            Icon(aw.IconInfo).
            Arg("auth").
            Valid(true)
        wf.SendFeedback()
        return
    }

    cfg.APIToken = token

    api, err := bb.NewAPI(cfg.URL, cfg.Username, cfg.APIToken)
    if err != nil {
        wf.FatalError(err)
    }

    maxCacheAge := time.Duration(cfg.CacheAge * int(time.Minute))
    if wf.Cache.Expired(repoCacheName, maxCacheAge) || wf.Cache.Expired(projectCacheName, maxCacheAge) {
        if err := refreshCache(api); err != nil {
            wf.FatalError(err)
        }
        wf.Rerun(0.3)
    }

    if a := autocomplete(opts.Query); a != "" {
        if err := wf.Alfred.RunTrigger(a, opts.Query); err != nil {
            wf.FatalError(err)
        }
        return
    }

    if opts.Commits {
        runCommits(api)
        wf.SendFeedback()
        return
    }

    if opts.Branches {
        runBranches(api)
        wf.SendFeedback()
        return
    }

    if opts.Tags {
        runTags(api)
        wf.SendFeedback()
        return
    }

    if opts.PullRequests {
        runPullRequests(api)
        wf.SendFeedback()
        return
    }

    if opts.Projects {
        runProjects()
        wf.SendFeedback()
        return
    }

    parsedQuery := parseQuery(opts.Query)
    runSearch(parsedQuery)

    wf.Filter(parsedQuery.Text)

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
