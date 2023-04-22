package utils

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Autommit struct {
	OpenAiApiKey string
	Context context.Context
	OpenAiClient openai.Client
	CommitInfo Commit
}

func NewAutommit(openAiApiKey string) (*Autommit) {
	ctx := context.Background()
	client := openai.NewClient(openAiApiKey)
	return &Autommit{
		OpenAiApiKey: openAiApiKey,
		Context: ctx,
		OpenAiClient: *client,
	}
}

func (a *Autommit) GeneratePrompt(gitDiff string, footer string) (string) {
	var prologue = "Analyze the following output of git diff and create a git commit message in json format, following conventional commits specification. Json fields are: commit_message, commit_message_long."
	var restrictions = "Include only information found on the git diff output. If there are multiple files, provide a detailed changelog in commit_message_long. commit_message should always be short and concise. Explain code changes in simple terms in commit_message_long."
	var messageFooter = fmt.Sprintf("Sign your commits in the bottom of commit_message_long. Use the following format: %s\n", footer)
	return fmt.Sprintf("%s. Diff is: ```\n%s\n```\n.%s\n%s\n", prologue, gitDiff, restrictions, messageFooter)
}
