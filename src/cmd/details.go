package cmd

import (
	"fmt"

	"github.com/rwilgaard/alfred-bitbucket-search/src/internal/util"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "list actions for a repository",
	Run: func(_ *cobra.Command, _ []string) {
		if cfg.ProjectKey == "" || cfg.RepoSlug == "" {
			wf.FatalError(fmt.Errorf("project_key and repo_slug variables must be set"))
		}

		wf.NewItem("Open in Browser").
			Subtitle("Open "+cfg.RepoName+" in your browser").
			Icon(util.GetIcon("browser")).
			Var("link", cfg.RepoHTMLURL).
			Arg("open_browser").
			Valid(true)

		wf.NewItem("View Commits").
			Subtitle("Show commits for " + cfg.RepoName).
			Icon(util.GetIcon("commit")).
			Arg("commits").
			Valid(true)

		wf.NewItem("View Branches").
			Subtitle("Show branches for " + cfg.RepoName).
			Icon(util.GetIcon("branch")).
			Arg("branches").
			Valid(true)

		wf.NewItem("View Tags").
			Subtitle("Show tags for " + cfg.RepoName).
			Icon(util.GetIcon("tag")).
			Arg("tags").
			Valid(true)

		wf.NewItem("View Pull Requests").
			Subtitle("Show pull requests for " + cfg.RepoName).
			Icon(util.GetIcon("pull-request")).
			Arg("pullrequests").
			Valid(true)

		if cfg.RepoCloneHTTP != "" {
			wf.NewItem("Copy HTTP Clone URL").
				Subtitle(cfg.RepoCloneHTTP).
				Icon(util.GetIcon("copy")).
				Var("action", "copy").
				Var("notification_title", "HTTP Clone URL Copied").
				Var("notification_text", cfg.RepoName).
				Var("copy_value", cfg.RepoCloneHTTP).
				Arg("copy").
				Valid(true)
		}

		if cfg.RepoCloneSSH != "" {
			wf.NewItem("Copy SSH Clone URL").
				Subtitle(cfg.RepoCloneSSH).
				Icon(util.GetIcon("copy")).
				Var("action", "copy").
				Var("notification_title", "SSH Clone URL Copied").
				Var("notification_text", cfg.RepoName).
				Var("copy_value", cfg.RepoCloneSSH).
				Arg("copy").
				Valid(true)
		}

		wf.NewItem("Back to Repositories").
			Subtitle("Return to your repository list").
			Icon(util.GetIcon("go-back")).
			Var("list_query", cfg.ListQuery).
			Arg("go-back").
			Valid(true)

		alfredutils.HandleFeedback(wf)
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)
}
