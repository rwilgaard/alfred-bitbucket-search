package cmd

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	bb "github.com/rwilgaard/bitbucket-go-api"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

type parsedQuery struct {
	Text     string
	Projects []string
}

func parseQuery(query string) *parsedQuery {
	q := new(parsedQuery)
	projectRegex := regexp.MustCompile(`^@\w+`)
	var text []string
	for w := range strings.FieldsSeq(query) {
		if projectRegex.MatchString(w) {
			q.Projects = append(q.Projects, w[1:])
		} else {
			text = append(text, w)
		}
	}
	q.Text = strings.Join(text, " ")
	return q
}

func autocomplete(query string) string {
	if slices.Contains(strings.Fields(query), "@") {
		return "projects"
	}
	return ""
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "search bitbucket repositories",
	Args:  cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
			return
		}

		var query string
		if len(args) > 0 {
			query = args[0]
		}

		if a := autocomplete(query); a != "" {
			if err := wf.Alfred.RunTrigger(a, query); err != nil {
				wf.FatalError(err)
			}
			return
		}

		var repos []*bb.RepositoryList
		if err := alfredutils.LoadCache(wf, repoCacheName, &repos); err != nil {
			wf.FatalError(err)
		}

		maxAge := time.Duration(cfg.CacheAge) * time.Minute
		if err := alfredutils.RefreshCache(wf, repoCacheName, maxAge, []string{"cache"}); err != nil {
			wf.FatalError(err)
		}

		parsed := parseQuery(query)

		for _, list := range repos {
			for _, repo := range list.Values {
				if len(parsed.Projects) > 0 && !slices.Contains(parsed.Projects, repo.Project.Key) {
					continue
				}

				var selfURL string
				if s := repo.Links["self"]; len(s) > 0 {
					selfURL = s[0].Href
				}

				item := wf.NewItem(repo.Name).
					Subtitle(repo.Project.Name).
					Match(fmt.Sprintf("%s %s %s %s", repo.Name, repo.Slug, repo.Project.Name, repo.Project.Key)).
					UID(fmt.Sprintf("%s/%s", repo.Project.Key, repo.Slug)).
					Var("item_url", selfURL).
					Arg("open").
					Valid(true)

				var sshURL, httpURL string
				sshIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "ssh" })
				if sshIdx != -1 {
					sshURL = repo.Links["clone"][sshIdx].Href
				}
				httpIdx := slices.IndexFunc(repo.Links["clone"], func(l bb.Link) bool { return l.Name == "http" })
				if httpIdx != -1 {
					httpURL = repo.Links["clone"][httpIdx].Href
				}

				item.NewModifier(aw.ModCmd).
					Subtitle("View repository details and actions").
					Var("repo_name", repo.Name).
					Var("repo_slug", repo.Slug).
					Var("project_key", repo.Project.Key).
					Var("repo_html_url", selfURL).
					Var("repo_clone_http", httpURL).
					Var("repo_clone_ssh", sshURL).
					Var("list_query", query).
					Arg("details").
					Valid(true)
			}
		}

		if len(parsed.Text) > 0 {
			wf.Filter(parsed.Text)
		}

		alfredutils.HandleFeedback(wf)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
