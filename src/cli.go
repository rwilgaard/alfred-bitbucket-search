package main

import "flag"

var (
    opts = &options{}
    cli = flag.NewFlagSet("alfred-bitbucket-search", flag.ContinueOnError)
)

type options struct {
    // Arguments
    Query string

    // Commands
    Commits      bool
    Branches     bool
    Tags         bool
    PullRequests bool
    Update       bool
}

func init() {
    cli.BoolVar(&opts.Commits, "commits", false, "show commits for repository")
    cli.BoolVar(&opts.Branches, "branches", false, "show branches for repository")
    cli.BoolVar(&opts.Tags, "tags", false, "show tags for repository")
    cli.BoolVar(&opts.PullRequests, "pullrequests", false, "show pull requests for repository")
    cli.BoolVar(&opts.Update, "update", false, "check for updates")
}
