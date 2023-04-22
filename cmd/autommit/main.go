package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aorfanos/autommit/utils"
)

var (
	openAiApiKey = flag.String("openai-api-key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	path = flag.String("path", ".", "Path to the git repository")
)

func main() {
	flag.Parse()
	// redundant loop since we use the envvar as default for openAiApiKey
	// @TODO: reassess
	if (*openAiApiKey == "") {
		if (os.Getenv("OPENAI_API_KEY") == "") {
			fmt.Println("Please provide an OpenAI API token")
			os.Exit(1)
			return
		} else {
			*openAiApiKey = os.Getenv("OPENAI_API_KEY")
		}
	}

	// check if git is present in the system
	utils.CheckGitPresence()

	// add the file to the git repository
	utils.GitAdd(*path)

	var autommit = utils.NewAutommit(*openAiApiKey)

	err := autommit.CreateCompletionRequest(autommit.GeneratePrompt(utils.GitDiff()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// fmt.Println(prompt)

	// utils.ParseStringAsJson(prompt)

	// autommit.GitCommit(false)
}