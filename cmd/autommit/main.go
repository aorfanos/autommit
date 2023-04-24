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
	convCommitsType = flag.String("conventional-commits-type", "feat", "Will add the provided type to the commit message")
	// nonInteractive = flag.Bool("non-interactive", false, "Will automatically add, commit and push the commit to the remote repository")
)

func init() {
	flag.StringVar(convCommitsType, "t", "feat", "Alias of --conventional-commits-type")
}

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

	autommit, err := utils.NewAutommit(*openAiApiKey, *convCommitsType, *path)
	utils.ErrCheck(err)

	// add files to the commit
	autommit.GitAddDialogue()

	COMPLETIONLOOP:
	answer, err := autommit.CreateCompletionRequest(autommit.GeneratePrompt(utils.GitDiff(true, nil), *signCommitsMessage))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	err = autommit.ParseStringAsJson(answer)
	utils.ErrCheck(err)

	if (autommit.GitCommitDialogue()) {
		fmt.Println("Commit successful. Proceeding to push routine.")
		utils.GitPush()
	} else {
		fmt.Println("Will recreate a message")
		goto COMPLETIONLOOP
	}
}
