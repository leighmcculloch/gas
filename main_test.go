package main

import (
	"os"
	"testing"

	"4d63.com/testcli"
	"4d63.com/want"
)

func setupGit(t *testing.T) {
	dir := testcli.MkdirTemp(t)
	os.Setenv("HOME", dir)
	testcli.Exec(t, "git config --global user.email 'tests@example.com'")
	testcli.Exec(t, "git config --global user.name 'Tests'")
}

func TestNoRemote(t *testing.T) {
	setupGit(t)

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master     <none>
`)
}

func TestNoRemoteUntrackedFiles(t *testing.T) {
	setupGit(t)

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   <none>
`)
}

func TestNoRemoteStagedFiles(t *testing.T) {
	setupGit(t)

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   <none>
`)
}

func TestNoRemoteCommitted(t *testing.T) {
	setupGit(t)

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master     <none>
`)
}

func TestRemoteCommittedPushed(t *testing.T) {
	setupGit(t)

	remote := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, ``)
}

func TestRemoteUntracked(t *testing.T) {
	setupGit(t)

	remote := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")
	testcli.WriteFile(t, "file2", []byte{})

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   origin/master
`)
}

func TestRemoteStaged(t *testing.T) {
	setupGit(t)

	remote := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")
	testcli.WriteFile(t, "file2", []byte{})
	testcli.Exec(t, "git add . -v")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   origin/master
`)
}

func TestRemoteNotPushed(t *testing.T) {
	setupGit(t)

	remote := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")
	testcli.WriteFile(t, "file2", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add more files'")
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master  ↑  origin/master
`)
}

func TestRemoteNotPushedOtherBranch(t *testing.T) {
	setupGit(t)

	remote := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")
	testcli.Exec(t, "git checkout -b branch1")
	testcli.WriteFile(t, "file2", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add more files'")
	testcli.Exec(t, "git status")
	testcli.Exec(t, "git checkout -")
	testcli.Exec(t, "git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  branch1     <none>       
`)
}

func TestMultipleRepos(t *testing.T) {
	setupGit(t)

	remote1 := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote1)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	remote2 := testcli.MkdirTemp(t)
	testcli.Chdir(t, remote2)
	testcli.Exec(t, "git init --bare")
	testcli.Exec(t, "git status")

	dir := testcli.MkdirTemp(t)
	testcli.Chdir(t, dir)

	testcli.Mkdir(t, "repo1")
	testcli.Chdir(t, "repo1")
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote1)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")
	testcli.Exec(t, "git checkout -b branch1")
	testcli.WriteFile(t, "file2", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add more files'")
	testcli.Exec(t, "git status")
	testcli.Exec(t, "git checkout -")
	testcli.Exec(t, "git status")
	testcli.Chdir(t, "..")

	testcli.Mkdir(t, "repo2")
	testcli.Chdir(t, "repo2")
	testcli.Exec(t, "git init")
	testcli.Exec(t, "git remote add origin "+remote2)
	testcli.WriteFile(t, "file1", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add files'")
	testcli.Exec(t, "git push -u origin HEAD")
	testcli.Exec(t, "git status")
	testcli.Exec(t, "git checkout -b branch1")
	testcli.WriteFile(t, "file2", []byte{})
	testcli.Exec(t, "git add . -v")
	testcli.Exec(t, "git commit -m 'Add more files'")
	testcli.Exec(t, "git status")
	testcli.Exec(t, "git checkout -")
	testcli.Exec(t, "git status")
	testcli.Chdir(t, "..")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := testcli.Main(t, args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `repo1/
  branch1     <none>       
repo2/
  branch1     <none>       
`)
}
