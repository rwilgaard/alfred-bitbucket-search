# Bitbucket Search

A workflow for quickly finding repositories in your Bitbucket installation.

## Installation
* [Download the latest release](https://github.com/rwilgaard/alfred-bitbucket-search/releases)
* Open the downloaded file in Finder.
* If running on macOS Catalina or later, you _**MUST**_ add Alfred to the list of security exceptions for running unsigned software. See [this guide](https://github.com/deanishe/awgo/wiki/Catalina) for instructions on how to do this.

## Keywords

*You can change the default 'gs' keyword in the User configuration.*

* With `gs` you can fuzzy find repositories in Bitbucket. The default `⏎` action will open the highlighted repository in your browser.

## Actions

The following actions can be used on a highlighted repository:
* `⏎` opens the repository in your browser.
* `⌘` + `⏎` will show commits for the repository.
* `⌥` + `⏎` will show pull requests for the repository.
* `⌃` + `⏎` will show tags for the repository.
* `⇧` + `⏎` will show branches for the repository
* `⌘` `⇧` + `⏎` will copy SSH clone URL.
* `⌥` `⇧` + `⏎` will copy HTTP clone URL.
