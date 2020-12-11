package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var version = "<dev>"
var commit = ""
var date = ""

func main() {
	exitCode := run(os.Args, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	flag := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flag.SetOutput(stderr)
	flagVersion := flag.Bool("version", false, "print the version")
	flagHelp := flag.Bool("help", false, "print this help")
	flagNoColor := flag.Bool("no-color", false, "disable color")
	flagErrorCode := flag.Bool("e", false, "exit with error code if changes not pushed")
	flagFetchUpstream := flag.Bool("f", false, "fetch upstream")
	flagAll := flag.Bool("a", false, "print all branches")
	err := flag.Parse(args[1:])
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 2
	}

	if *flagVersion {
		fmt.Fprintf(stderr, "gas %s %s %s\n", version, commit, date)
		return 0
	}

	if *flagHelp {
		flag.Usage()
		return 0
	}

	if *flagNoColor {
		color.NoColor = true
	}

	all := *flagAll
	errorCode := *flagErrorCode
	fetchUpstream := *flagFetchUpstream

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "error getting current directory: %v", err)
		return 1
	}

	repos, err := getRepos(dir, fetchUpstream)
	if err != nil {
		fmt.Fprintf(stderr, "error getting git worktrees: %v", err)
		return 1
	}

	nameWidth, upstreamWidth, authorDateWidth := maxColumnWidths(repos)

	for _, r := range repos {
		if !all && !r.ChangesNotPushed() {
			continue
		}
		relPath, err := filepath.Rel(dir, r.path)
		if err != nil {
			relPath = r.path
		}
		fmt.Fprintf(stdout, "%s%c\n", relPath, filepath.Separator)
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

			authorDate := color.HiGreenString("%*s", authorDateWidth, b.authorDate)

			fmt.Fprintf(stdout, "  %-*s %s %-*s %s %s\n", nameWidth, b.name, icons, upstreamWidth, upstream, authorDate, b.commitMessageSubject)
		}
	}

	if errorCode {
		for _, r := range repos {
			if r.ChangesNotPushed() {
				return 1
			}
		}
	}

	return 0
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
	head                 bool
	name                 string
	dirty                bool
	ahead                bool
	behind               bool
	upstream             string
	authorDate           string
	commitMessageSubject string
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
		"--format=%(HEAD):%(refname:short):%(upstream:trackshort):%(upstream:short):%(authordate:relative):%(contents:subject)",
	)
	cmd.Dir = path
	cmd.Stdin = nil
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	outStr := strings.TrimSpace(out.String())
	if outStr == "" {
		headBranch, err := getHeadBranch(path)
		if err != nil {
			return nil, err
		}
		return []branch{{head: true, dirty: dirty, name: headBranch}}, nil
	}
	lines := strings.Split(outStr, "\n")
	branches := make([]branch, len(lines))
	for i, l := range lines {
		f := strings.Split(l, ":")
		branches[i] = branch{
			head:                 f[0] == "*",
			dirty:                f[0] == "*" && dirty,
			name:                 f[1],
			ahead:                f[2] == ">" || f[2] == "<>",
			behind:               f[2] == "<" || f[2] == "<>",
			upstream:             f[3],
			authorDate:           f[4],
			commitMessageSubject: f[5],
		}
	}

	return branches, nil
}

func getHeadBranch(path string) (string, error) {
	out := strings.Builder{}
	cmd := exec.Command(
		"git", "--no-pager", "symbolic-ref", "HEAD",
	)
	cmd.Dir = path
	cmd.Stdin = nil
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	outStr := strings.TrimSpace(out.String())
	if outStr == "" {
		return "", nil
	}
	branchName := strings.Replace(outStr, "refs/heads/", "", 1)
	return branchName, nil
}

func getBrancheRemotes(path string) ([]string, error) {
	out := strings.Builder{}
	cmd := exec.Command(
		"git", "--no-pager", "branch", "--sort=committerdate",
		"--format=%(upstream:remotename)",
	)
	cmd.Dir = path
	cmd.Stdin = nil
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
	cmd.Stdin = nil
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func maxColumnWidths(repos []repo) (nameWidth, upstreamWidth, authorDateWidth int) {
	for _, r := range repos {
		for _, b := range r.branches {
			if len(b.name) > nameWidth {
				nameWidth = len(b.name)
			}
			if len(b.upstream) > upstreamWidth {
				upstreamWidth = len(b.upstream)
			}
			if len(b.authorDate) > authorDateWidth {
				authorDateWidth = len(b.authorDate)
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
