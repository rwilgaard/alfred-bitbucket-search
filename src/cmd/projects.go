package cmd

import (
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
	bb "github.com/rwilgaard/bitbucket-go-api"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects [query]",
	Short: "list bitbucket projects",
	Args:  cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
			return
		}

		var projects *bb.ProjectList
		if err := alfredutils.LoadCache(wf, projectCacheName, &projects); err != nil {
			wf.FatalError(err)
		}

		maxAge := time.Duration(cfg.CacheAge) * time.Minute
		if err := alfredutils.RefreshCache(wf, projectCacheName, maxAge, []string{"cache"}); err != nil {
			wf.FatalError(err)
		}

		if projects == nil {
			alfredutils.HandleFeedback(wf)
			return
		}

		prevQuery := cfg.ListQuery

		for _, p := range projects.Values {
			item := wf.NewItem(p.Key).
				Match(fmt.Sprintf("%s %s", p.Key, p.Name)).
				UID(p.Key).
				Subtitle(p.Name).
				Arg(prevQuery+p.Key+" ").
				Var("project", p.Key).
				Valid(true)

			item.NewModifier(aw.ModCmd).
				Subtitle("Cancel").
				Arg("cancel")
		}

		alfredutils.HandleFeedback(wf)
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
