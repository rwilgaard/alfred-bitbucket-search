package cmd

import (
	"fmt"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/rwilgaard/alfred-bitbucket-search/src/internal/bitbucket"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "authenticate",
	RunE: func(_ *cobra.Command, _ []string) error {
		if cfg.URL == "" || cfg.Username == "" {
			return fmt.Errorf("please configure bitbucket_url and username in workflow variables")
		}
		_, pwd, err := zenity.Password(zenity.Title(fmt.Sprintf("Enter API Token for %s", cfg.Username)))
		if err != nil {
			return err
		}
		pwd = strings.TrimSpace(pwd)

		client, err := bitbucket.NewClient(cfg.URL, cfg.Username, pwd)
		if err != nil {
			return err
		}
		if err := client.TestAuthentication(); err != nil {
			zerr := zenity.Error(fmt.Sprintf("Authentication failed: %s", err), zenity.ErrorIcon)
			if zerr != nil {
				return err
			}
			return err
		}

		if err := wf.Keychain.Set(keychainAccount, pwd); err != nil {
			return err
		}
		fmt.Println("Successfully authenticated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
