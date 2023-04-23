package utils

import (
	"encoding/json"
	"fmt"

	"github.com/manifoldco/promptui"
)

func ErrCheck(err error) {
	if (err != nil) {
		fmt.Printf("error: %s", err)
		return
	}
}

func (a *Autommit) ParseStringAsJson(strSrc string)  {
	var c Commit
	err := json.Unmarshal([]byte(strSrc), &c)
	ErrCheck(err)
	a.CommitInfo.Message = c.Message
	a.CommitInfo.MessageLong = c.MessageLong
	// fmt.Printf("Message: %s\n", c.Message)
}

func ProceedSelector(title string, choices []string) (string, error) {
	selector := promptui.Select{
		Label: title,
		Items: choices,
	}
	_, result, err := selector.Run()
	return result, err
}