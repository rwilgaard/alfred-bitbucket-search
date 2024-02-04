package cmd

import (
    "fmt"
    "slices"
    "time"

    aw "github.com/deanishe/awgo"
    "github.com/rwilgaard/alfred-bitbucket-search/src/pkg/util"
    bb "github.com/rwilgaard/bitbucket-go-api"
    "github.com/rwilgaard/go-alfredutils/alfredutils"
    "github.com/spf13/cobra"
)

var (
    reposCmd = &cobra.Command{
        Use:   "repos",
        Short: "List repositories",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            query := args[0]

            if comp := util.CheckQueryForAutoCompletion(query); comp != "" {
                if err := wf.Alfred.RunTrigger(comp, query); err != nil {
                    wf.FatalError(err)
                }
                return
            }

            parsedQuery := util.ParseQuery(query)
            var repos []*bb.RepositoryList
            if err := alfredutils.LoadCache(wf, repoCacheName, &repos); err != nil {
                wf.FatalError(err)
            }

            maxCacheAge := time.Duration(cfg.CacheAge * int(time.Minute))
            if err := alfredutils.RefreshCache(wf, repoCacheName, maxCacheAge); err != nil {
                wf.FatalError(err)
            }

            for _, list := range repos {
                for _, repo := range list.Values {
                    if len(parsedQuery.Projects) > 0 && !slices.Contains(parsedQuery.Projects, repo.Project.Key) {
                        continue
                    }

                    i := wf.NewItem(repo.Name).
                        Subtitle(repo.Project.Name).
                        Match(fmt.Sprintf("%s %s %s %s", repo.Name, repo.Slug, repo.Project.Name, repo.Project.Key)).
                        UID(fmt.Sprintf("%s/%s", repo.Project.Key, repo.Slug)).
                        Var("projectKey", repo.Project.Key).
                        Var("repoSlug", repo.Slug).
                        Var("link", repo.Links["self"][0].Href).
                        Var("lastQuery", query).
                        Valid(true)

                    i.NewModifier(aw.ModCmd).
                        Subtitle("Show Commits").
                        Arg("commits").
                        Valid(true)

                    i.NewModifier(aw.ModCtrl).
                        Subtitle("Show Tags").
                        Arg("tags").
                        Valid(true)

                    i.NewModifier(aw.ModOpt).
                        Subtitle("Show Pull Requests").
                        Arg("pullrequests").
                        Valid(true)

                    i.NewModifier(aw.ModShift).
                        Subtitle("Show Branches").
                        Arg("branches").
                        Valid(true)

                    sshIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "ssh" })
                    httpIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "http" })
                    i.NewModifier(aw.ModOpt, aw.ModShift).
                        Subtitle("Copy HTTP clone URL").
                        Arg("copy").
                        Var("copy_value", repo.Links["clone"][httpIdx].Href).
                        Valid(true)

                    i.NewModifier(aw.ModCmd, aw.ModShift).
                        Subtitle("Copy SSH clone URL").
                        Arg("copy").
                        Var("copy_value", repo.Links["clone"][sshIdx].Href).
                        Valid(true)
                }
            }

            wf.Filter(parsedQuery.Text)
            alfredutils.HandleFeedback(wf)
        },
    }
)

func init() {
    rootCmd.AddCommand(reposCmd)
}
