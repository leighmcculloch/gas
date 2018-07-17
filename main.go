package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current directory: %v", err)
		os.Exit(1)
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path
		}
		fmt.Printf("%s:\n", color.CyanString(relPath))
		cmdBranches(path)

		return filepath.SkipDir
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func cmdDirtyStatus(path string) string {
	out := strings.Builder{}
	cmd := exec.Command("git", "--no-pager", "status", "--porcelain")
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	if out.Len() > 0 {
		return "M"
	}
	return " "
}

func cmdBranches(path string) {
	cmd := exec.Command(
		"git",
		"--no-pager",
		"branch",
		"-vv",
		"--color",
		"--sort=committerdate",
		"--format="+
			"  "+
			"%(align:20,left)%(refname:short)%(end)"+
			"%(align:3,left)%(color:bold red)"+
			"%(if)%(HEAD)%(then)"+cmdDirtyStatus(path)+"%(else) %(end)"+
			"%(if:equals=>)%(upstream:trackshort)%(then)↑ %(end)"+
			"%(if:equals=<)%(upstream:trackshort)%(then) ↓%(end)"+
			"%(if:equals=<>)%(upstream:trackshort)%(then)↑↓%(end)"+
			"%(if:equals==)%(upstream:trackshort)%(then)  %(end)"+
			"%(if:equals=)%(upstream:trackshort)%(then)  %(end)"+
			"%(color:reset) %(end)"+
			"%(align:20,left)%(if)%(upstream:short)%(then)%(upstream:short)%(else)%(color:bold red)no upstream%(color:reset)%(end)%(end)",
	)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
}
