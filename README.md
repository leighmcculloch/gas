<div align="center"><img alt="gas" src="README-gas.png" /></div>
<p align="center">
<a href="https://github.com/leighmcculloch/tldr/actions"><img alt="Build" src="https://github.com/leighmcculloch/tldr/workflows/build/badge.svg" /></a>
<a href="https://goreportcard.com/report/github.com/leighmcculloch/tldr"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/leighmcculloch/tldr" /></a>
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
  master               ↑  origin/master
  base                    <none>
gs/
  master               ↑↓ origin/master
helloworld/
  push-with-request       <none>
  add-vcr-2           M   origin/add-vcr-2
  add-vcr              ↑↓ origin/add-vcr
```
