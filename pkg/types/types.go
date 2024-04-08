package types

// Content represents the content of a message.
type Content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Message represents a message in the conversation.
type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

// Output represents the output data from the API response.
type Output struct {
	Response         string `json:"response"`
	TokensConsumed   uint64 `json:"tokens_consumed"`
	ResponseTimeMS   uint64 `json:"response_time_ms"`
	PromptTokenCount uint64 `json:"prompt_token_count"`
}

// APIResponse represents the response from the AWS Bedrock chat API.
type APIResponse struct {
	Output struct {
		Message        Message `json:"message"`
		TokensConsumed uint64  `json:"tokens_consumed"`
	} `json:"output"`
}
