package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func main() {
	flagHelp := flag.Bool("help", false, "print this help")
	flagNoColor := flag.Bool("no-color", false, "disable color")
	flagErrorCode := flag.Bool("e", false, "exit with error code if changes not pushed")
	flagFetchUpstream := flag.Bool("f", false, "fetch upstream")
	flagAll := flag.Bool("a", false, "print all branches")

	flag.Parse()

	if *flagHelp {
		flag.Usage()
		return
	}

	if *flagNoColor {
		color.NoColor = true
	}

	all := *flagAll
	errorCode := *flagErrorCode
	fetchUpstream := *flagFetchUpstream

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current directory: %v", err)
		os.Exit(1)
	}

	repos, err := getRepos(dir, fetchUpstream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting git worktrees: %v", err)
		os.Exit(1)
	}

	nameWidth, upstreamWidth := maxColumnWidths(repos)

	for _, r := range repos {
		if !all && !r.ChangesNotPushed() {
			continue
		}
		relPath, err := filepath.Rel(dir, r.path)
		if err != nil {
			relPath = r.path
		}
		fmt.Printf("%s%c\n", relPath, filepath.Separator)
		for _, b := range r.branches {
			if !all && !b.ChangesNotPushed() {
				continue
			}

			dirty, ahead, behind := iconChars(b)
			upstream := b.upstream
			if upstream == "" {
				upstream = color.HiRedString("<none>")
			}

			icons := color.HiRedString("%c%c%c", dirty, ahead, behind)

			fmt.Printf("  %-*s %s %-*s\n", nameWidth, b.name, icons, upstreamWidth, upstream)
		}
	}

	if errorCode {
		for _, r := range repos {
			if r.ChangesNotPushed() {
				os.Exit(1)
			}
		}
	}
}

type repo struct {
	path     string
	branches []branch
}

func (r repo) ChangesNotPushed() bool {
	for _, b := range r.branches {
		if b.ChangesNotPushed() {
			return true
		}
	}
	return false
}

type branch struct {
	head     bool
	name     string
	dirty    bool
	ahead    bool
	behind   bool
	upstream string
}

func (b branch) ChangesNotPushed() bool {
	return b.dirty || b.ahead || b.behind || b.upstream == ""
}

func getRepos(dir string, fetchUpstream bool) ([]repo, error) {
	repos := []repo{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		dotGitFileInfo, err := os.Stat(filepath.Join(path, ".git"))
		if os.IsNotExist(err) || !dotGitFileInfo.IsDir() {
			return nil
		}

		repo, err := getRepo(path, fetchUpstream)

		repos = append(repos, repo)

		return filepath.SkipDir
	})
	return repos, err
}

func getRepo(path string, fetchUpstream bool) (repo, error) {
	if fetchUpstream {
		remotes, err := getBrancheRemotes(path)
		if err != nil {
			return repo{}, err
		}

		for _, r := range remotes {
			err := fetch(path, r)
			if err != nil {
				return repo{}, err
			}
		}
	}

	branches, err := getBranches(path)
	if err != nil {
		return repo{}, err
	}

	r := repo{
		path:     path,
		branches: branches,
	}

	return r, nil
}

func isDirty(path string) (bool, error) {
	out := strings.Builder{}
	cmd := exec.Command("git", "--no-pager", "status", "--porcelain")
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return false, err
	}
	return out.Len() > 0, nil
}

func getBranches(path string) ([]branch, error) {
	dirty, err := isDirty(path)
	if err != nil {
		return nil, err
	}

	out := strings.Builder{}
	cmd := exec.Command(
		"git", "--no-pager", "branch", "--sort=committerdate",
		"--format=%(HEAD):%(refname:short):%(upstream:trackshort):%(upstream:short)",
	)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	branches := make([]branch, len(lines))
	for i, l := range lines {
		f := strings.Split(l, ":")
		branches[i] = branch{
			head:     f[0] == "*",
			dirty:    f[0] == "*" && dirty,
			name:     f[1],
			ahead:    f[2] == ">" || f[2] == "<>",
			behind:   f[2] == "<" || f[2] == "<>",
			upstream: f[3],
		}
	}

	return branches, nil
}

func getBrancheRemotes(path string) ([]string, error) {
	out := strings.Builder{}
	cmd := exec.Command(
		"git", "--no-pager", "branch", "--sort=committerdate",
		"--format=%(upstream:remotename)",
	)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	remotes := strings.Split(strings.TrimSpace(out.String()), "\n")
	remotesSeen := map[string]struct{}{}
	for _, r := range remotes {
		remotesSeen[r] = struct{}{}
	}
	dedupedRemotes := make([]string, 0, len(remotesSeen))
	for r := range remotesSeen {
		dedupedRemotes = append(dedupedRemotes, r)
	}

	return dedupedRemotes, nil
}

func fetch(path, remote string) error {
	cmd := exec.Command("git", "--no-pager", "fetch", remote)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func maxColumnWidths(repos []repo) (nameWidth, upstreamWidth int) {
	for _, r := range repos {
		for _, b := range r.branches {
			if len(b.name) > nameWidth {
				nameWidth = len(b.name)
			}
			if len(b.upstream) > upstreamWidth {
				upstreamWidth = len(b.upstream)
			}
		}
	}
	return
}

func iconChars(b branch) (dirty, ahead, behind rune) {
	dirty = ' '
	if b.dirty {
		dirty = 'M'
	}
	ahead = ' '
	if b.ahead {
		ahead = '↑'
	}
	behind = ' '
	if b.behind {
		behind = '↓'
	}
	return dirty, ahead, behind
}
