package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func ErrCheck(err error) {
	if (err != nil) {
		fmt.Printf("error: %s", err)
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

func PopulateFileAddSelector(gitDiffChangesString string) ([]string, error) {
	if (gitDiffChangesString == "") {
		return []string{}, fmt.Errorf("No files to add")
	}
	// Split the string by newlines, trim the trailing newline and trim the spaces
	splitted := strings.Split(strings.TrimSuffix(gitDiffChangesString, "\n"), "\n")
	for i, v := range splitted {
		splitted[i] = strings.TrimSpace(v)
		splitted[i] = strings.TrimSuffix(splitted[i], "\n")
	}
	return splitted, nil
}
