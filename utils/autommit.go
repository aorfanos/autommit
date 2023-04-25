package utils

import (
	"context"
	"fmt"

	git "github.com/go-git/go-git/v5"
	openai "github.com/sashabaranov/go-openai"
)

type Autommit struct {
	OpenAiApiKey string
	Context context.Context
	OpenAiClient openai.Client
	PgpSign bool
	CommitInfo Commit
	Type string
	GitConfig GitConfig
}

func NewAutommit(openAiApiKey, commitType, path, gitUser, gitEmail string) (*Autommit, error) {
	ctx := context.Background()
	client := openai.NewClient(openAiApiKey)
	repo, err := git.PlainOpen(path)
	ErrCheck(err)
	workTree, err := repo.Worktree()
	ErrCheck(err)
	headRef, err := repo.Head()
	ErrCheck(err)
	return &Autommit{
		OpenAiApiKey: openAiApiKey,
		Context: ctx,
		OpenAiClient: *client,
		Type: commitType,
		GitConfig: GitConfig{
			Author: gitUser,
			AuthorMail: gitEmail,
			RepoPath: path,
			Repo: repo,
			Worktree: workTree,
			HeadRef: headRef,
		},
	}, nil
}


// function GeneratePrompt generates a prompt for the OpenAI API
// @param gitDiff string
// @param footer string
// @return string
func (a *Autommit) GeneratePrompt(gitDiff string, footer string) (string) {
	var prologue = "Analyze the following output of git diff and create a git commit message in json format, describing the changes following the Conventional Commits specification (study https://www.conventionalcommits.org/en/v1.0.0/#specification). Json fields are only: commit_message, commit_message_long."
	var restrictions = fmt.Sprintf("Include only truthful information relevant to the git diff output. If there are multiple files, provide a detailed changelog in commit_message_long. commit_message should always be short and concise. Both commit_message and commit_message_long should be of the same type. Do not return anything else other than JSON response. Answers should only be in Conventional Commits format. The commit is of the conventional commits type %s.", a.Type)
	var messageFooter = fmt.Sprintf("The last words of commit_message_long should be: %s.\n", footer)
	return fmt.Sprintf("%s. Diff is: ```\n%s\n```\n.%s %s\n", prologue, gitDiff, restrictions, messageFooter)
}
