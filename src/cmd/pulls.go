package cmd

import (
	"fmt"

	"github.com/rwilgaard/alfred-bitbucket-search/src/internal/util"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var pullsCmd = &cobra.Command{
	Use:     "pulls [query]",
	Aliases: []string{"pullrequests"},
	Short:   "list pull requests",
	Args:    cobra.MaximumNArgs(1),
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

		prs, err := client.GetPullRequests(cfg.ProjectKey, cfg.RepoSlug)
		if err != nil {
			wf.FatalError(err)
		}

		for _, p := range prs.Values {
			wf.NewItem(p.Title).
				Subtitle(fmt.Sprintf("%s ➔ %s", p.FromRef.DisplayID, p.ToRef.DisplayID)).
				Icon(util.GetIcon("pull-request")).
				Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/overview", cfg.URL, cfg.ProjectKey, cfg.RepoSlug, p.ID)).
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
	rootCmd.AddCommand(pullsCmd)
}
