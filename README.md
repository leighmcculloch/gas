# gs

Check the state of all your git worktrees.

## Install

### Binary (Linux; macOS; Windows)

Download the binary from the [releases][] page.

[releases]: https://github.com/leighmcculloch/gs/releases

### From Source

```
go get 4d63.com/gs
```

## Usage

Print a list of repositories in the current directory that have branches not
pushed upstream, or dirty working directories:

```
gs
```

Example:

```
$ gs
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

## Usage (all)

Print all branches:

```
gs -all
```
