package cmd

import (
	"github.com/spf13/cobra"
)

var (
    tagsCmd     = &cobra.Command{
        Use:   "commits",
        Short: "List commits for repository",
        Run: func(cmd *cobra.Command, args []string) {
        },
    }
)

func init() {
    rootCmd.AddCommand(tagsCmd)
}
