package cmd

import (
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
	bb "github.com/rwilgaard/bitbucket-go-api"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var (
    projectsCmd = &cobra.Command{
        Use:   "projects",
        Short: "List projects",
        Run: func(cmd *cobra.Command, args []string) {
            var projects *bb.ProjectList
            if err := alfredutils.LoadCache(wf, projectCacheName, &projects); err != nil {
                wf.FatalError(err)
            }

            maxCacheAge := time.Duration(cfg.CacheAge * int(time.Minute))
            if err := alfredutils.RefreshCache(wf, projectCacheName, maxCacheAge); err != nil {
                wf.FatalError(err)
            }

            prevQuery, _ := wf.Config.Env.Lookup("prev_query")
            for _, project := range projects.Values {
                i := wf.NewItem(project.Key).
                    Match(fmt.Sprintf("%s %s", project.Key, project.Name)).
                    UID(project.Key).
                    Subtitle(project.Name).
                    Arg(prevQuery+project.Key+" ").
                    Var("project", project.Key).
                    Valid(true)

                i.NewModifier(aw.ModCmd).
                    Subtitle("Cancel").
                    Arg("cancel")

            }

            alfredutils.HandleFeedback(wf)
        },
    }
)

func init() {
    rootCmd.AddCommand(projectsCmd)
}
