package utils

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	Author string
	AuthorMail string
	RepoPath string
	Repo *git.Repository
	Worktree *git.Worktree
	HeadRef *plumbing.Reference
	PGPKeyRing openpgp.KeyRing
}


type Commit struct {
	Message string `json:"commit_message"`
	MessageLong string `json:"commit_message_long"`
	FilesAffected []string `json:"files_affected"`
}

func CheckGitPresence() bool {
	cmd := exec.Command("git", "--version")
	_, err := cmd.Output()
	return err == nil
}

// func GitAddDialogue will ask the user to select files to add to the commit
// if no files are available to select and no files are already staged, it will exit the program
// if no files are available to select but some files are already staged, it will proceed to the commit dialogue
func (a *Autommit) GitAddDialogue() {
	fileList, err := a.PopulateFileList()
	ErrCheck(err)

	stagedFilesExist, _ := a.CheckForStagedFiles()

	if (len(fileList) == 0) {
		if (!stagedFilesExist) {
			fmt.Println("No files to stage, or already staged - exiting")
			os.Exit(0)
		} else {
			fmt.Println("No new files to stage, proceeding to commit dialogue")
		}
		return
	}

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

func (a *Autommit) GitCommitDialogue() (regenerate bool) {
	result, err := ProceedSelector(gitCommitSelectorQTitle, gitCommitSelectorQChoices)
	ErrCheck(err)
	IF_EVAL_START:
	if (result == gitCommitSelectorQChoices[1]) { // no
		// unstage all files
		a.UnstageFiles()
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
		err = a.GitCommit()
		ErrCheck(err)
		return true
	}
	return
}

func (a *Autommit) GitCommit() (error) {
	// Read the PGP key from a file, or use any other method to obtain the key.
	keyRingFile, err := os.Open(a.PgpKeyPath)
	if err != nil {
		return err
	}
	defer keyRingFile.Close()

	// Read the armored keyring
	block, err := armor.Decode(keyRingFile)
	if err != nil {
		return err
	}

	// Read the keyring
	keyRing, err := openpgp.ReadKeyRing(block.Body)
	if err != nil {
		return err
	}

	// Get the current HEAD reference
	headRef, err := a.GitConfig.Repo.Head()
	if err != nil {
		return err
	}

	// Get the commit at the HEAD
	headCommit, err := a.GitConfig.Repo.CommitObject(headRef.Hash())
	if err != nil {
		return err
	}

	// Create a new commit message
	commitMessage := fmt.Sprintf("%s\n\n%s", a.CommitInfo.Message, a.CommitInfo.MessageLong)
	author := &object.Signature{
		Name:  fmt.Sprintf("%s", a.GitConfig.Author),
		Email: fmt.Sprintf("%s", a.GitConfig.AuthorMail),
		When:  time.Now(),
	}

	// Create the commit with the PGP signature
	commitHash, err := a.GitConfig.Worktree.Commit(
		commitMessage,
		&git.CommitOptions{
			Author:    author,
			Committer: author,
			Parents:   []plumbing.Hash{headCommit.Hash},
			SignKey:   keyRing[0],
		},
	)
	ErrCheck(err)

	// commit the file(s)
	_, err = a.GitConfig.Repo.CommitObject(commitHash)
	ErrCheck(err)
	return err
}

func (a *Autommit) GitPush() error {
	result, err := ProceedSelector(gitPushSelectorQTitle, gitPushSelectorQChoices)
	ErrCheck(err)

	if (result == gitPushSelectorQChoices[1]) { // no
		fmt.Println("Push aborted")
		a.UnstageFiles()
		return nil
	} else if (result == gitPushSelectorQChoices[0]) { // yes
		err = a.GitConfig.Repo.Push(&git.PushOptions{})
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
			a.CommitInfo.FilesAffected = append(a.CommitInfo.FilesAffected, fileName)
		}
	}
	return fileList, err
}

// CheckForStagedFiles checks if there are any staged files
// if there are, it will return true, and list them
func (a *Autommit) CheckForStagedFiles() (exist bool, fileNames []string) {
	files, err := a.GitConfig.Worktree.Status()
	ErrCheck(err)

	for fileName, fileStatus := range files {
		if (fileStatus.Staging == git.Modified ||
			fileStatus.Staging == git.Added) {
			fileNames = append(fileNames, fileName)
			return true, fileNames
		}
	}
	return false, nil
}

// UnstageFiles unstages all files
func (a *Autommit) UnstageFiles() (error) {
	err := a.GitConfig.Worktree.Reset(&git.ResetOptions{
		Mode: git.MixedReset,
	})
	ErrCheck(err)
	return err
}
