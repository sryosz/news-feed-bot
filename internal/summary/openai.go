package summary

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"sync"
)

type OpenAiSummarizer struct {
	client  *openai.Client
	prompt  string
	enabled bool
	mu      sync.Mutex
}

func New(apiKey string, prompt string) *OpenAiSummarizer {
	s := &OpenAiSummarizer{
		client: openai.NewClient(apiKey),
		prompt: prompt,
	}

	log.Printf("openai summarizer enabled: %v", apiKey != "")

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAiSummarizer) Summarize(ctx context.Context, text string) (string, error) {
	const op = "summary.Summarize"

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return "", nil
	}

	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("%s%s", text, s.prompt),
			},
		},
		MaxTokens:   256,
		Temperature: 0.7,
		TopP:        1,
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	rawSummary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary, nil
	}

	sentences := strings.Split(rawSummary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}
