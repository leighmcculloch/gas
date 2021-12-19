<div align="center"><img alt="gas" src="README-gas.png" /></div>
<p align="center">
<a href="https://github.com/leighmcculloch/gas/actions"><img alt="Build" src="https://github.com/leighmcculloch/gas/workflows/build/badge.svg" /></a>
<a href="https://goreportcard.com/report/github.com/leighmcculloch/gas"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/leighmcculloch/gas" /></a>
<a href="https://github.com/leighmcculloch/gas/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/leighmcculloch/gas.svg" /></a>
</p>

Check the state of all your git worktrees in subdirectories.

Checks for any untracked or modified files, and any unpushed branches in all git repositories below the current directory.

## Install

### Binary (Linux; macOS; Windows)

Download the binary from the [releases][] page.

[releases]: https://github.com/leighmcculloch/gas/releases

### Homebrew (Linux; macOS)

```
brew install 4d63/gas/gas
```

### From Source

```
go get 4d63.com/gas
```

## Usage

Print a list of repositories in the current directory that have branches not
pushed upstream, or dirty working directories:

```
gas [-e] [-a]
```

Add `-e` to exit with an error code if there are changes not pushed.

Add `-a` to print all branches.

Example:

```
$ gas
fork-stretchr-testify/
  master               ↑  origin/master    15 minutes ago Add new assertion
  base                    <none>              7 hours ago Fix test
gs/
  master               ↑↓ origin/master      3 months ago Release version 2.1.0
helloworld/
  push-with-request       <none>             23 hours ago Add push with request
  add-vcr-2           M   origin/add-vcr-2  12 months ago Add vcr
  add-vcr              ↑↓ origin/add-vcr    12 months ago Add vcr
```
