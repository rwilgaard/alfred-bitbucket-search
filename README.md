# Bitbucket Search

Quickly find repositories in your Bitbucket Server/Data Center installation from Alfred.

## Requirements

- [Alfred](https://www.alfredapp.com/) with the Powerpack
- A Bitbucket Server / Data Center instance and an HTTP access token

## Installation

1. [Download the latest release](https://github.com/rwilgaard/alfred-bitbucket-search/releases)
2. Open the downloaded `.alfredworkflow` file to import it into Alfred.
3. On macOS Catalina or later you _**MUST**_ add Alfred to the list of security exceptions for running unsigned software. See [this guide](https://github.com/deanishe/awgo/wiki/Catalina).
4. Set `Bitbucket URL` and `Username` in the workflow configuration.
5. Trigger `gs`, press `⏎` on "You're not logged in", and enter your API token.

### Creating an API token

In Bitbucket: **Manage account → HTTP access tokens → Create token**. Grant at least read permission on projects and repositories.

## Usage

### `gs [query]` — search repositories

Fuzzy-find repositories. Matches on repo name, slug, project name and project key.

| Action | Result |
|--------|--------|
| `⏎` | Open the repository in your browser |
| `⌘` + `⏎` | Open the details / actions menu for the repository |

### Filters

- Filter by project with `@projectkey` syntax (e.g. `gs @infra deploy`). Type `@` on its own to open a project picker.

### Details menu (`⌘` + `⏎`)

From a repository you can:

- Open in Browser
- View Commits (`⌘` + `⏎` on a commit shows the full message)
- View Branches
- View Tags
- View Pull Requests
- Copy HTTP Clone URL
- Copy SSH Clone URL
- Back to Repositories (restores your previous search)

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `bitbucket_url` | — | Base URL of your Bitbucket installation (e.g. `https://bitbucket.example.com`) |
| `username` | — | Your Bitbucket username |
| `cache_age` | `180` | Repository/project cache TTL in minutes |
| `keyword` | `gs` | Keyword that triggers the search |

The API token is stored in the macOS Keychain, never in workflow variables.

## Development

```sh
make build          # build arch-specific binaries
make package-alfred # build + zip into .alfredworkflow
make zip-alfred     # zip into .alfredworkflow without rebuilding
make fmt            # format with gofumpt
make release V=x.y.z # bump version, package, and tag locally
```

Requires Go 1.26+ and `golangci-lint`.
