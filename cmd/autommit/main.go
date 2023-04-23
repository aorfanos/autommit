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
	pgpSignedCommit = flag.Bool("pgp-sign", true, "Will sign the commit with the default PGP key")
	signCommitsMessage = flag.String("sign-commits-with-message", "Created by autommit 🦄", "Will add the provided message to the long commit message")
	// nonInteractive = flag.Bool("non-interactive", false, "Will automatically add, commit and push the commit to the remote repository")
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
	// utils.GitAdd(*path)
	utils.GitAdd()

	var autommit = utils.NewAutommit(*openAiApiKey)

	COMPLETIONLOOP:
	answer, err := autommit.CreateCompletionRequest(autommit.GeneratePrompt(utils.GitDiff(true, nil), *signCommitsMessage))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	autommit.ParseStringAsJson(answer)

	if (autommit.GitCommit()) {
		fmt.Println("Commit successful. Proceeding to push routine.")
		utils.GitPush()
	} else {
		fmt.Println("Will recreate a message")
		goto COMPLETIONLOOP
	}
}