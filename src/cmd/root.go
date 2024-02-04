package cmd

import (
    aw "github.com/deanishe/awgo"
    "github.com/deanishe/awgo/update"
    "github.com/rwilgaard/alfred-bitbucket-search/src/pkg/bitbucket"
    "github.com/rwilgaard/go-alfredutils/alfredutils"
    "github.com/spf13/cobra"
    "go.deanishe.net/fuzzy"
)

type workflowConfig struct {
    URL      string `env:"bitbucket_url"`
    Username string `env:"username"`
    CacheAge int    `env:"cache_age"`
}

const (
    repo             = "rwilgaard/alfred-bitbucket-search"
    keychainAccount  = "alfred-bitbucket-search"
    repoCacheName    = "repositories.json"
    projectCacheName = "projects.json"
)

var (
    wf      *aw.Workflow
    bs      *bitbucket.BitbucketService
    cfg     = &workflowConfig{}
    rootCmd = &cobra.Command{
        Use:   "bitbucket",
        Short: "bitbucket is a CLI to be used by Alfred for searching Bitbucket repositories",
    }
)

func Execute() {
    wf.Run(run)
}

func run() {
    alfredutils.AddClearAuthMagic(wf, keychainAccount)

    if err := alfredutils.InitWorkflow(wf, cfg); err != nil {
        wf.FatalError(err)
    }

    if err := alfredutils.CheckForUpdates(wf); err != nil {
        wf.FatalError(err)
    }

    alfredutils.HandleAuthentication(wf, keychainAccount)
    token, err := wf.Keychain.Get(keychainAccount)
    if err != nil {
        wf.FatalError(err)
    }

    bs, err = bitbucket.NewBitbucketService(cfg.URL, cfg.Username, token)
    if err != nil {
        wf.FatalError(err)
    }

    if err := rootCmd.Execute(); err != nil {
        wf.FatalError(err)
    }
}

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
}
