package utils

import (
	"os/exec"
)

type Commit struct {
	Message string `json:"commit_message"`
	MessageLong string `json:"commit_message_long"`
}

func CheckGitPresence() bool {
	cmd := exec.Command("git", "--version")
	_, err := cmd.Output()
	return err == nil
}

func GitAdd(path string) {
	cmd := exec.Command("git", "add", path)
	_, err := cmd.Output()
	ErrCheck(err)
}

func GitDiff() (string) {
	cmd := exec.Command("git", "diff", "-m", "HEAD")
	diff, err := cmd.Output()
	ErrCheck(err)
	return string(diff)
}

func (a *Autommit) GitCommit(signCommits bool) (string) {
	var cmd *exec.Cmd
	if (signCommits) {
		cmd = exec.Command("git", "commit", "-S", "-m", a.CommitInfo.Message, "-m", a.CommitInfo.MessageLong)
	} else {
		cmd = exec.Command("git", "commit", "-m", a.CommitInfo.Message, "-m", a.CommitInfo.MessageLong)
	}
	commit, err := cmd.Output()
	ErrCheck(err)
	return string(commit)
}
 