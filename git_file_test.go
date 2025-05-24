package main

import (
	"testing"

	"4d63.com/testcli"
	"4d63.com/want"
)

func TestGitFileNotDirectory(t *testing.T) {
	setupGit(t)

	// Create a parent directory
	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	
	// Create a directory with a .git file (not a directory)
	testcli.Mkdir(t, "tricky-repo")
	testcli.Chdir(t, "tricky-repo")
	testcli.WriteFile(t, ".git", []byte{})
	testcli.Chdir(t, "..")
	
	// Run gas in the parent directory
	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	
	// Check that we get a clean output without errors
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "") // There should be no errors
	want.Eq(t, stdout, "") // No repositories should be found
}