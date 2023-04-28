package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/muja/goconfig"
)

func ErrCheck(err error) {
	if (err != nil) {
		log.Fatal(fmt.Sprintf("error: %s\n", err.Error()))
	}
}

func ErrReturn(err error) error {
	if (err != nil) {
		return err
	}
	return nil
}

func (a *Autommit) ParseStringAsJson(strSrc string) (error) {
	var c Commit
	err := json.Unmarshal([]byte(strSrc), &c)
	if err != nil {
		return err
	} else if (c.Message == "" || c.MessageLong == "") {
		return fmt.Errorf("Invalid JSON")
	}
	a.CommitInfo.Message = c.Message
	a.CommitInfo.MessageLong = c.MessageLong
	return nil
}

func ProceedSelector(title string, choices []string) (string, error) {
	selector := promptui.Select{
		Label: title,
		Items: choices,
	}
	_, result, err := selector.Run()
	return result, err
}

func ProceedEditor(title, target string) (string, error) {
	prompt := promptui.Prompt{
		Label: title,
		Default: target,
	}
	result, err := prompt.Run()
	ErrCheck(err)
	return result, err
}

// func PopulateGitUserInfo() will parse a .gitconfig file
// and populate the Autommit struct with the user's name and email
func (a *Autommit) PopulateGitUserInfo() (error) {
	user, err := user.Current()
	ErrReturn(err)

	if (a.GitConfig.FilePath == "~/.gitconfig") {
		a.GitConfig.FilePath = filepath.Join(user.HomeDir, ".gitconfig")
	}

	bytes, err := ioutil.ReadFile(a.GitConfig.FilePath)
	ErrReturn(err)

	config, _, err := goconfig.Parse(bytes)
	ErrReturn(err)

	a.GitConfig.Author = config["user.name"]
	a.GitConfig.AuthorMail = config["user.email"]

	if (a.GitConfig.Author == "" || a.GitConfig.AuthorMail == "") {
		err = fmt.Errorf("No git user info found")
	}

	return err
}

func ShowVersion(version string) {
	fmt.Printf("Autommit version %s 🦄\n", version)
}

func FindDotGit(repoPath string) (path string, err error) {
	var dotGit string = fmt.Sprintf("%s/.git", repoPath)
    // Check current directory for file
    _, err = os.Stat(dotGit);
    if err == nil {
        // File found in current directory
        return repoPath, nil
    }

    // Traverse parent directories
    for i := 0; i < getDirectoryLevelsToRoot(); i++ {
        err = os.Chdir("..")
        if err != nil {
            // Unable to change directory
            return "", err
        }

        // Check new directory for file
        _, err = os.Stat(dotGit)
        if err == nil {
            // File found in parent directory
			pwd, err := os.Getwd()
			if err != nil {
				return "", err
			}
            return fmt.Sprintf("%s", pwd), nil
        }
    }

    // File not found in any directory
    return "", fmt.Errorf("Dir .git not found in your current dir, nor any parent dir - exiting")
}

func getDirectoryLevelsToRoot() int {
    var levels int
    dir, err := os.Getwd()
    if err != nil {
        // Error getting current directory
        return -1
    }

    for dir != "/" {
        dir = filepath.Dir(dir)
        levels++
    }

    return levels
}