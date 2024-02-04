package cmd

import (
	"github.com/spf13/cobra"
)

var (
    brancesCmd     = &cobra.Command{
        Use:   "commits",
        Short: "List brances for repository",
        Run: func(cmd *cobra.Command, args []string) {
        },
    }
)

func init() {
    rootCmd.AddCommand(brancesCmd)
}
