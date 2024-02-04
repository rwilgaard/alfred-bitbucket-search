package cmd

import (
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/rwilgaard/alfred-bitbucket-search/src/pkg/util"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var (
    projectKeyFlag string
    repoSlugFlag   string
    commitsCmd     = &cobra.Command{
        Use:   "commits",
        Short: "List commits for repository",
        Run: func(cmd *cobra.Command, args []string) {
            wf.Configure(aw.SuppressUIDs(true))
            commits, err := bs.GetCommits(projectKeyFlag, repoSlugFlag)
            if err != nil {
                wf.FatalError(err)
            }

            wf.NewItem("Go back").
                Icon(util.BackIcon).
                Arg("go-back").
                Valid(true)

            icon := aw.Icon{Value: fmt.Sprintf("%s/icons/commit.png", wf.Dir())}
            for _, c := range commits.Values {
                t := time.UnixMilli(c.CommitterTimestamp).Format("02-01-2006 15:04")
                i := wf.NewItem(c.Message).
                    Subtitle(fmt.Sprintf("%s  •  %s  •  %s", c.DisplayID, c.Committer.Name, t)).
                    Icon(&icon).
                    Var("link", fmt.Sprintf("%s/projects/%s/repos/%s/commits/%s", cfg.URL, projectKeyFlag, repoSlugFlag, c.ID)).
                    Valid(true)

                i.NewModifier(aw.ModCmd).
                    Subtitle("Show full commit message.").
                    Arg("commit").
                    Var("message", c.Message).
                    Valid(true)
            }

            alfredutils.HandleFeedback(wf)
        },
    }
)

func init() {
    commitsCmd.Flags().StringVarP(&projectKeyFlag, "project", "p", "", "project key")
    commitsCmd.Flags().StringVarP(&repoSlugFlag, "repo", "r", "", "repo slug")
    rootCmd.AddCommand(commitsCmd)
}
