package utils

import (
	"fmt"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/manifoldco/promptui"
)

var (
	gitCommitSelectorQTitle = "Proceed with the commit?"
	gitCommitSelectorQChoices = []string{"✅ Yes", "❌ No", "💸 Regenerate", "🔍 Edit"}
	gitPushSelectorQTitle = "Push to remote?"
	gitPushSelectorQChoices = []string{"✅ Yes", "❌ No"}
	gitAddSelectorQTitle = "Select files to add to the commit"

    cmd *exec.Cmd
)

type GitConfig struct {
	RepoPath string
	Repo *git.Repository
	Worktree *git.Worktree
}


type Commit struct {
	Message string `json:"commit_message"`
	MessageLong string `json:"commit_message_long"`
}

func CheckGitPresence() bool {
	cmd := exec.Command("git", "--version")
	_, err := cmd.Output()
	return err == nil
}

func (a *Autommit) GitAddDialogue() {
	fileList, err := a.PopulateFileList()
	ErrCheck(err)

	index := -1
	var result string

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label: gitAddSelectorQTitle,
			Items: fileList,
		}

		index, result, err = prompt.Run()

		if (index == -1) {
			fileList = append(fileList, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	_, err = a.GitConfig.Worktree.Add(result)
	ErrCheck(err)

	addAnother, err := ProceedSelector("Add another file?", []string{"✅ Yes", "❌ No"})
	ErrCheck(err)
	if (addAnother == "✅ Yes") {
		a.GitAddDialogue()
	} else {
		return
	}
}

// if not staged, it will accept a list of strings to be passed to git
// as arguments
func GitDiff(staged bool, args []string) (string) {
	if (staged) {
		cmd = exec.Command("git", "diff", "--staged", "HEAD")
	} else {
		cmd = exec.Command("git", args...)
	}
	diff, err := cmd.Output()
	ErrCheck(err)
	return string(diff)
}

func (a *Autommit) GitCommit() (regenerate bool) {
	result, err := ProceedSelector(gitCommitSelectorQTitle, gitCommitSelectorQChoices)
	ErrCheck(err)
	IF_EVAL_START:
	if (result == gitCommitSelectorQChoices[1]) { // no
		fmt.Println("Commit aborted")
		os.Exit(0)
	} else if (result == gitCommitSelectorQChoices[2]) { // regenerate
		return false
	} else if (result == gitCommitSelectorQChoices[3]) { // edit
		a.CommitInfo.Message, err = ProceedEditor("Change commit message", a.CommitInfo.Message)
		ErrCheck(err)
		a.CommitInfo.MessageLong, err = ProceedEditor("Change commit message long", a.CommitInfo.MessageLong)
		ErrCheck(err)
		// set result to yes after editing and goto start of if
		// so that logic proceeds as if the user selected yes
		// since we assume intent to commit after editing
		result = gitCommitSelectorQChoices[0]
		goto IF_EVAL_START
	} else if (result == gitCommitSelectorQChoices[0]) { // yes
		if (a.PgpSign) {
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

// PopulateFileList populates the file list for the add prompt selector
func (a *Autommit) PopulateFileList() ([]string, error) {
	var fileList []string

	// get files
	files, err := a.GitConfig.Worktree.Status()
	ErrCheck(err)

	// populate file list
	for fileName, fileStatus := range files {
		// skip files that are staged
		if (fileStatus.Staging == git.Modified || 
			fileStatus.Staging == git.Added) {
			continue
		} else {
			fileList = append(fileList, fileName)
		}
	}
	return fileList, err
}