package main

import (
	"strings"
	"testing"

	"4d63.com/testcli"
	"4d63.com/want"
)

func TestMixedReposWithBareGitNamedDotGit(t *testing.T) {
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
	
	// Create a bare git repository named .git
	testcli.Exec(t, "git init --bare .git")
	
	// Run gas in the parent directory
	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	
	// It should only find the regular repo, not the .git bare repo
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "") // There should be no errors
	want.True(t, len(stdout) > 0) // Should have some output
	want.True(t, stdout != "" && stdout != "\n") // Should have non-empty output
	want.True(t, strings.Contains(stdout, "regular-repo/")) // Should find the regular repo
}