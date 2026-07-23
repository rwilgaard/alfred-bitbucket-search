package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:    "cache",
	Short:  "refresh cached data in the background",
	Hidden: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		log.Println("[cache] fetching repositories and projects...")
		if err := fetchAndCache(); err != nil {
			return err
		}
		log.Println("[cache] done")
		return nil
	},
}

func fetchAndCache() error {
	client, err := setupClient()
	if err != nil {
		msg := fmt.Sprintf("Cache failed: credentials error: %s", err)
		if zerr := wf.Alfred.RunTrigger("error", msg); zerr != nil {
			log.Printf("Alfred error trigger failed: %v", zerr)
		}
		return err
	}

	repos, err := client.GetAllRepositories()
	if err != nil {
		msg := fmt.Sprintf("Cache failed: could not fetch repos: %s", err)
		if zerr := wf.Alfred.RunTrigger("error", msg); zerr != nil {
			log.Printf("Alfred error trigger failed: %v", zerr)
		}
		return err
	}
	if err := wf.Cache.StoreJSON(repoCacheName, repos); err != nil {
		return err
	}

	projects, err := client.GetAllProjects()
	if err != nil {
		msg := fmt.Sprintf("Cache failed: could not fetch projects: %s", err)
		if zerr := wf.Alfred.RunTrigger("error", msg); zerr != nil {
			log.Printf("Alfred error trigger failed: %v", zerr)
		}
		return err
	}
	if err := wf.Cache.StoreJSON(projectCacheName, projects); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
