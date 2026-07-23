package cmd

import (
	"fmt"

	"github.com/rwilgaard/alfred-bitbucket-search/src/internal/util"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var branchesCmd = &cobra.Command{
	Use:   "branches [query]",
	Short: "list branches",
	Args:  cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if cfg.ProjectKey == "" || cfg.RepoSlug == "" {
			wf.FatalError(fmt.Errorf("project_key and repo_slug variables must be set"))
		}

		if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
			return
		}

		wf.NewItem("Back to Actions").
			Subtitle("Return to the actions menu for " + cfg.RepoName).
			Icon(util.GetIcon("go-back")).
			Arg("details").
			Valid(true)

		client, err := setupClient()
		if err != nil {
			wf.FatalError(err)
		}

		branches, err := client.GetBranches(cfg.ProjectKey, cfg.RepoSlug)
		if err != nil {
			wf.FatalError(err)
		}

		for _, b := range branches.Values {
			commit := b.LatestCommit
			if len(commit) > 10 {
				commit = commit[0:10]
			}
			wf.NewItem(b.DisplayID).
				Subtitle(fmt.Sprintf("Latest commit: %s", commit)).
				Icon(util.GetIcon("branch")).
				Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/browse?at=%s", cfg.URL, cfg.ProjectKey, cfg.RepoSlug, b.ID)).
				Arg("open").
				Valid(true)
		}

		if len(args) > 0 {
			wf.Filter(args[0])
		}

		alfredutils.HandleFeedback(wf)
	},
}

func init() {
	rootCmd.AddCommand(branchesCmd)
}
