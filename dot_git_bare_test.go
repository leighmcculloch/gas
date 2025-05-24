package main

import (
	"testing"

	"4d63.com/testcli"
	"4d63.com/want"
)

func TestBareGitDirectoryNamedDotGit(t *testing.T) {
	setupGit(t)

	// Create a parent directory
	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	
	// Create a bare git repository named .git
	testcli.Exec(t, "git init --bare .git")
	
	// Run gas in the parent directory
	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	
	// Check that we get a clean output without errors
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "") // There should be no errors
	want.Eq(t, stdout, "") // No repositories should be found
}