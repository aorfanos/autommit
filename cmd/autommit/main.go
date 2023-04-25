package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aorfanos/autommit/utils"
)

const version = "0.0.7"

var (
	openAiApiKey = flag.String("openai-api-key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	path = flag.String("path", ".", "Path to the git repository")
	pgpKeyPath = flag.String("pgp-key-path", "", "Path to the PGP key")
	signCommitsMessage = flag.String("sign-commits-with-message", "Created by autommit 🦄", "Will add the provided message to the long commit message")
	convCommitsType = flag.String("conventional-commits-type", "feat", "Will add the provided type to the commit message")
	gitUser = flag.String("git-user", "", "Will set the git user")
	gitEmail = flag.String("git-mail", "", "Will set the git email")
	gitConfigPath = flag.String("git-config-path", "~/.gitconfig", "Will set the git config path")
	showVersion = flag.Bool("version", false, "Will show the version of autommit")
	// nonInteractive = flag.Bool("non-interactive", false, "Will automatically add, commit and push the commit to the remote repository")
)

func init() {
	flag.StringVar(convCommitsType, "t", "feat", "Alias of --conventional-commits-type")
}

func main() {
	flag.Parse()
	if (*showVersion) {
		utils.ShowVersion(version)
		os.Exit(0)
	}
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

	autommit, err := utils.NewAutommit(
		version,
		*openAiApiKey,
		*convCommitsType,
		*path,
		*pgpKeyPath,
		*gitConfigPath,
	)
	utils.ErrCheck(err)

	// get git user info from .gitconfig
	err = autommit.PopulateGitUserInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get pgp keyring
	if (*pgpKeyPath != "") {
		err = autommit.GetOpenPGPKeyring()
		utils.ErrCheck(err)
	} else {
		fmt.Println("No PGP key provided. Will not sign commits.")
	}

	// add files to the commit
	autommit.GitAddDialogue()

	COMPLETIONLOOP:
	answer, err := autommit.CreateCompletionRequest(autommit.GeneratePrompt(utils.GitDiff(true, nil), *signCommitsMessage))
	if err != nil {
		// if there's an issue reaching the OpenAI API, we unstage the files
		fmt.Println(err)
		autommit.UnstageFiles()
		os.Exit(1)
		return
	}

	err = autommit.ParseStringAsJson(answer)
	utils.ErrCheck(err)

	if (autommit.GitCommitDialogue()) {
		fmt.Println("Commit successful. Proceeding to push routine.")
		autommit.GitPush()
	} else {
		fmt.Println("Will recreate a message")
		goto COMPLETIONLOOP
	}
}
