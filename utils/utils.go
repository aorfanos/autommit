package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/muja/goconfig"
)

func ErrCheck(err error) {
	if (err != nil) {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
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

func (a *Autommit) PopulateGitUserInfo() (error) {
	user, err := user.Current()
	ErrCheck(err)

	gitconfig := filepath.Join(user.HomeDir, ".gitconfig")
	bytes, err := ioutil.ReadFile(gitconfig)
	ErrCheck(err)

	config, _, err := goconfig.Parse(bytes)
	ErrCheck(err)

	a.GitConfig.Author = config["user.name"]
	a.GitConfig.AuthorMail = config["user.email"]

	if (a.GitConfig.Author == "" || a.GitConfig.AuthorMail == "") {
		err = fmt.Errorf("No git user info found, will not proceed")
	}

	return err
}
