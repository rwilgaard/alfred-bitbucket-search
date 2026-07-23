package cmd

import (
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/rwilgaard/alfred-bitbucket-search/src/internal/util"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var commitsCmd = &cobra.Command{
	Use:   "commits [query]",
	Short: "list commits",
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

		commits, err := client.GetCommits(cfg.ProjectKey, cfg.RepoSlug)
		if err != nil {
			wf.FatalError(err)
		}

		for _, c := range commits.Values {
			t := time.UnixMilli(c.CommitterTimestamp).Format("02-01-2006 15:04")
			i := wf.NewItem(c.Message).
				Subtitle(fmt.Sprintf("%s  |  %s  |  %s", c.DisplayID, c.Committer.Name, t)).
				Icon(util.GetIcon("commit")).
				Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/commits/%s", cfg.URL, cfg.ProjectKey, cfg.RepoSlug, c.ID)).
				Arg("open").
				Valid(true)

			i.NewModifier(aw.ModCmd).
				Subtitle("Show full commit message.").
				Arg("commit").
				Var("message", c.Message).
				Valid(true)
		}

		if len(args) > 0 {
			wf.Filter(args[0])
		}

		alfredutils.HandleFeedback(wf)
	},
}

func init() {
	rootCmd.AddCommand(commitsCmd)
}
