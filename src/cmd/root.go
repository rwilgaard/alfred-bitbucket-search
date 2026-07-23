package cmd

import (
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/rwilgaard/alfred-bitbucket-search/src/internal/bitbucket"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
	"go.deanishe.net/fuzzy"
)

type workflowConfig struct {
	URL           string `env:"bitbucket_url"`
	Username      string `env:"username"`
	CacheAge      int    `env:"cache_age"`
	RepoSlug      string `env:"repo_slug"`
	ProjectKey    string `env:"project_key"`
	RepoName      string `env:"repo_name"`
	RepoHTMLURL   string `env:"repo_html_url"`
	RepoCloneSSH  string `env:"repo_clone_ssh"`
	RepoCloneHTTP string `env:"repo_clone_http"`
	ListQuery     string `env:"list_query"`
}

const (
	repo             = "rwilgaard/alfred-bitbucket-search"
	keychainAccount  = "alfred-bitbucket-search"
	repoCacheName    = "repositories.json"
	projectCacheName = "projects.json"
)

var (
	wf      *aw.Workflow
	cfg     = &workflowConfig{}
	rootCmd = &cobra.Command{
		Use:           "alfred-bitbucket-search",
		Short:         "Bitbucket search for Alfred",
		SilenceUsage:  true,
		SilenceErrors: true,
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
	if err := rootCmd.Execute(); err != nil {
		wf.FatalError(err)
	}
}

func setupClient() (*bitbucket.Client, error) {
	if cfg.URL == "" || cfg.Username == "" {
		return nil, fmt.Errorf("bitbucket_url and username must be configured")
	}
	token, err := wf.Keychain.Get(keychainAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from keychain: %w", err)
	}
	return bitbucket.NewClient(cfg.URL, cfg.Username, token)
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
