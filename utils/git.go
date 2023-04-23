package utils

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	gitCommitSelectorQTitle = "Proceed with the commit?"
	gitCommitSelectorQChoices = []string{"✅ Yes", "❌ No", "💸 Regenerate"}
	gitPushSelectorQTitle = "Push to remote?"
	gitPushSelectorQChoices = []string{"✅ Yes", "❌ No"}

    cmd *exec.Cmd
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
	cmd := exec.Command("git", "diff", "--staged", "HEAD")
	diff, err := cmd.Output()
	ErrCheck(err)
	return string(diff)
}

func (a *Autommit) GitCommit(signCommits bool) (regenerate bool) {
	result, err := ProceedSelector(gitCommitSelectorQTitle, gitCommitSelectorQChoices)
	ErrCheck(err)
	if (result == gitCommitSelectorQChoices[1]) { // no
		fmt.Println("Commit aborted")
		os.Exit(0)
	} else if (result == gitCommitSelectorQChoices[2]) { // regenerate
		return false
	} else if (result == gitCommitSelectorQChoices[0]) { // yes
		if (signCommits) {
			cmd = exec.Command("git", "commit", "-S", "-m", a.CommitInfo.Message, "-m", a.CommitInfo.MessageLong)
		} else {
			cmd = exec.Command("git", "commit", "-m", a.CommitInfo.Message, "-m", a.CommitInfo.MessageLong)
		}
		_, err := cmd.Output()
		ErrCheck(err)
		return true
	}
	return
}

func GitPush() error {
	result, err := ProceedSelector(gitPushSelectorQTitle, gitPushSelectorQChoices)
	ErrCheck(err)

	if (result == gitPushSelectorQChoices[1]) { // no
		fmt.Println("Push aborted")
		return nil
	} else if (result == gitPushSelectorQChoices[0]) { // yes
		cmd := exec.Command("git", "push")
		_, err := cmd.Output()
		return err
	}
	return nil
}
