package main

import (
	"os"
	"testing"

	"4d63.com/testcli"
	"4d63.com/want"
)

func setupGit(t *testing.T) {
	h := testcli.Helper{TB: t}
	dir := h.MkdirTemp()
	os.Setenv("HOME", dir)
	h.Exec("git config --global user.email 'tests@example.com'")
	h.Exec("git config --global user.name 'Tests'")
}

func TestNoRemote(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master     <none>
`)
}

func TestNoRemoteUntrackedFiles(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.WriteFile("file1", []byte{})
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   <none>
`)
}

func TestNoRemoteStagedFiles(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   <none>
`)
}

func TestNoRemoteCommitted(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master     <none>
`)
}

func TestRemoteCommittedPushed(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	remote := h.MkdirTemp()
	h.Chdir(remote)
	h.Exec("git init --bare")
	h.Exec("git status")

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.Exec("git remote add origin " + remote)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, ``)
}

func TestRemoteUntracked(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	remote := h.MkdirTemp()
	h.Chdir(remote)
	h.Exec("git init --bare")
	h.Exec("git status")

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.Exec("git remote add origin " + remote)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")
	h.WriteFile("file2", []byte{})

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   origin/master
`)
}

func TestRemoteStaged(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	remote := h.MkdirTemp()
	h.Chdir(remote)
	h.Exec("git init --bare")
	h.Exec("git status")

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.Exec("git remote add origin " + remote)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")
	h.WriteFile("file2", []byte{})
	h.Exec("git add . -v")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master M   origin/master
`)
}

func TestRemoteNotPushed(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	remote := h.MkdirTemp()
	h.Chdir(remote)
	h.Exec("git init --bare")
	h.Exec("git status")

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.Exec("git remote add origin " + remote)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")
	h.WriteFile("file2", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add more files'")
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  master  ↑  origin/master
`)
}

func TestRemoteNotPushedOtherBranch(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	remote := h.MkdirTemp()
	h.Chdir(remote)
	h.Exec("git init --bare")
	h.Exec("git status")

	dir := h.MkdirTemp()
	h.Chdir(dir)
	h.Exec("git init")
	h.Exec("git remote add origin " + remote)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")
	h.Exec("git checkout -b branch1")
	h.WriteFile("file2", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add more files'")
	h.Exec("git status")
	h.Exec("git checkout -")
	h.Exec("git status")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `./
  branch1     <none>       
`)
}

func TestMultipleRepos(t *testing.T) {
	setupGit(t)

	h := testcli.Helper{TB: t}

	remote1 := h.MkdirTemp()
	h.Chdir(remote1)
	h.Exec("git init --bare")
	h.Exec("git status")

	remote2 := h.MkdirTemp()
	h.Chdir(remote2)
	h.Exec("git init --bare")
	h.Exec("git status")

	dir := h.MkdirTemp()
	h.Chdir(dir)

	h.Mkdir("repo1")
	h.Chdir("repo1")
	h.Exec("git init")
	h.Exec("git remote add origin " + remote1)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")
	h.Exec("git checkout -b branch1")
	h.WriteFile("file2", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add more files'")
	h.Exec("git status")
	h.Exec("git checkout -")
	h.Exec("git status")
	h.Chdir("..")

	h.Mkdir("repo2")
	h.Chdir("repo2")
	h.Exec("git init")
	h.Exec("git remote add origin " + remote2)
	h.WriteFile("file1", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add files'")
	h.Exec("git push -u origin HEAD")
	h.Exec("git status")
	h.Exec("git checkout -b branch1")
	h.WriteFile("file2", []byte{})
	h.Exec("git add . -v")
	h.Exec("git commit -m 'Add more files'")
	h.Exec("git status")
	h.Exec("git checkout -")
	h.Exec("git status")
	h.Chdir("..")

	args := []string{"gas", "-no-color"}
	exitCode, stdout, stderr := h.Main(args, nil, run)
	want.Eq(t, exitCode, 0)
	want.Eq(t, stderr, "")
	want.Eq(t, stdout, `repo1/
  branch1     <none>       
repo2/
  branch1     <none>       
`)
}
