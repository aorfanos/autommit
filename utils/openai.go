package utils

import (
	"errors"
	"fmt"
	"io"
	"strings"

	// "github.com/sashabaranov/go-openai"
	openai "github.com/sashabaranov/go-openai"
)

func (a *Autommit) NewOpenAiClient() *openai.Client {
	return openai.NewClient(a.OpenAiApiKey)
}

func (a *Autommit) CreateCompletionRequest(prompt string) (string, error) {

	var buildResp []string
	var formattedResp string

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1000,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := a.OpenAiClient.CreateChatCompletionStream(a.Context, req)
	if err != nil {
		return "exited with error while creating stream", err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		// when response is nil, it means the stream is closed
		// we then join the text segment chatgpt provides and return it as string
		if errors.Is(err, io.EOF) {
			formattedResp = strings.Join(buildResp, "")
			fmt.Println(formattedResp)
			return formattedResp, nil
		}

		if err != nil {
			return "exited with error", err
		}

		buildResp = append(buildResp, fmt.Sprintf(response.Choices[0].Delta.Content))
	}
}
