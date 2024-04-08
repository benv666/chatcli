package api

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"benvbin/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	defaultModel = "anthropic.claude-v2"
	version      = "bedrock-2023-05-31"
	region       = "us-east-1"
)

type ClaudeMessageContentSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      []byte `json:"data"`
}

type ClaudeMessageContent struct {
	Type   string                      `json:"type"`
	Text   string                      `json:"text,omitempty"`
	Source *ClaudeMessageContentSource `json:"source,omitempty"`
}

type ClaudeMessage struct {
	Role    string                 `json:"role"`
	Content []ClaudeMessageContent `json:"content"`
}

type ClaudeRequest struct { // see https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-anthropic-claude-messages.html
	AntropicVersion string          `json:"anthropic_version"` // bedrock-2023-05-31
	MaxTokens       int             `json:"max_tokens"`
	System          string          `json:"system,omitempty"`
	Messages        []ClaudeMessage `json:"messages"`
	// Omitting optional request parameters for now
}

type ClaudeUsage struct {
	InputTokens  uint64 `json:"input_tokens"`
	OutputTokens uint64 `json:"output_tokens"`
}

type ClaudeResponse struct {
	Id           string                 `json:"id"`
	Model        string                 `json:"model"`
	Type         string                 `json:"type"`
	Role         string                 `json:"role"`
	Content      []ClaudeMessageContent `json:"content"`
	StopReason   string                 `json:"stop_reason"`
	StopSequence string                 `json:"stop_sequence"`
	Usage        ClaudeUsage            `json:"usage"`
}

// MakeRequest makes a request to the AWS Bedrock chat API for Claude.
// Inspired by https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/bedrock-runtime/hello/hello.go
func MakeRequest(prompt, model string) (*types.Output, error) {
	if model == "" {
		model = defaultModel
		// model = "fail" // DEBUG
	}

	// Logger client
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	service := bedrockruntime.NewFromConfig(sdkConfig)

	msgs := []ClaudeMessage{
		{
			Role: "user",
			Content: []ClaudeMessageContent{
				{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}

	request := ClaudeRequest{
		AntropicVersion: "bedrock-2023-05-31",
		MaxTokens:       2000,
		Messages:        msgs,
	}

	body, err := json.Marshal(request)
	if err != nil {
		log.Panicln("Couldn't marshal the request: ", err)
	}

	// Debug
	// log.Printf("Request:\n%s\n", body)
	/* DEBUG
	if len(body) > 0 {
		return nil, nil
	}
	*/

	startTime := time.Now()
	input := &bedrockruntime.InvokeModelInput{
		Body:        body,
		ModelId:     aws.String(model),
		ContentType: aws.String("application/json"),
	}

	result, err := service.InvokeModel(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	responseTimeMS := uint64(time.Since(startTime).Nanoseconds()) / 1e6

	// log.Printf("Result:\n%s\n", result.Body)
	var response ClaudeResponse
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	output := &types.Output{
		Response:         response.Content[0].Text, // TODO: check the rest, assume it's only 1 output for now
		ResponseTimeMS:   responseTimeMS,
		TokensConsumed:   response.Usage.OutputTokens,
		PromptTokenCount: response.Usage.InputTokens,
	}

	return output, nil
}
