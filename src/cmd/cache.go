package cmd

import (
	"fmt"
	"log"

	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"
)

var (
    cacheCmd = &cobra.Command{
        Use:   "cache",
        Short: "Refresh cache",
        RunE: func(cmd *cobra.Command, args []string) error {
            log.Printf("[cache] fetching repositories...")
            repos, err := bs.GetAllRepositories()
            if err != nil {
                zerr := zenity.Error(
                    fmt.Sprintf("Repository caching failed: %s", err),
                    zenity.ErrorIcon,
                )
                if zerr != nil {
                    return zerr
                }
                return err
            }

            if err := wf.Cache.StoreJSON(repoCacheName, repos); err != nil {
                return err
            }
            log.Printf("[cache] repositories fetched")

            log.Printf("[cache] fetching projects...")
            projects, err := bs.GetAllProjects()
            if err != nil {
                zerr := zenity.Error(
                    fmt.Sprintf("Project caching failed: %s", err),
                    zenity.ErrorIcon,
                )
                if zerr != nil {
                    return zerr
                }
                return err
            }

            if err := wf.Cache.StoreJSON(projectCacheName, projects); err != nil {
                return err
            }
            log.Printf("[cache] projects fetched")

            return nil
        },
    }
)

func init() {
    rootCmd.AddCommand(cacheCmd)
}
