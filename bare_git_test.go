package main

import (
	"testing"

	"4d63.com/testcli"
	"4d63.com/want"
)

func TestBareGitDirectory(t *testing.T) {
	setupGit(t)

	// Create a parent directory
	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	
	// Create a regular git repository
	testcli.Mkdir(t, "regular-repo")
	testcli.Chdir(t, "regular-repo")
	testcli.Exec(t, "git init")
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Chdir(t, "..")

	// Create a bare git repository
	testcli.Exec(t, "git init --bare bare-repo.git")
	
	// Manually create a problem case - a directory that looks like a git directory but isn't
	// Create a regular directory with .git inside but not a proper git structure
	testcli.Mkdir(t, "fake-git")
	testcli.Mkdir(t, "fake-git/.git")
	
	// Run gas in the parent directory
	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	
	// It should only find the regular repo, not the bare repo or fake-git
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "") // There should be no errors for bare repos
	want.Eq(t, stdout, `regular-repo/
  master     <none> 0 seconds ago Add files
`)
}

func TestBareGitWithSubdirectories(t *testing.T) {
	setupGit(t)

	// Create a parent directory
	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	
	// Create a bare git repository with some subdirectories
	testcli.Exec(t, "git init --bare bare-repo.git")
	testcli.Mkdir(t, "bare-repo.git/subdirectory")
	
	// Run gas in the parent directory
	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	
	// Check that we get a clean output without errors
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "") // There should be no errors
	want.Eq(t, stdout, "") // No repositories should be found
}