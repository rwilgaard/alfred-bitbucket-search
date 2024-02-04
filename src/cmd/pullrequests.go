package cmd

import (
	"github.com/spf13/cobra"
)

var (
    pullrequestsCmd     = &cobra.Command{
        Use:   "commits",
        Short: "List pullrequests for repository",
        Run: func(cmd *cobra.Command, args []string) {
        },
    }
)

func init() {
    rootCmd.AddCommand(pullrequestsCmd)
}
