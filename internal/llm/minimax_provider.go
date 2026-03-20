package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

var ErrMissingMiniMaxAPIKey = errors.New("llm: missing MINIMAX_API_KEY for minimax provider")

type MiniMaxProvider struct {
	model   string
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewMiniMaxProvider(model, baseURL, apiKey string, timeout time.Duration) (*MiniMaxProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrMissingMiniMaxAPIKey
	}
	if strings.TrimSpace(model) == "" {
		model = "MiniMax-M2.7"
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.minimax.io/v1"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &MiniMaxProvider{
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		client:  &http.Client{Timeout: timeout},
	}, nil
}

func (p *MiniMaxProvider) Model() string {
	return p.model
}

func (p *MiniMaxProvider) Generate(ctx context.Context, question string, chunks []domain.RetrievedChunk) (string, error) {
	if len(chunks) == 0 {
		return "I couldn't find relevant content in the current knowledge base.", nil
	}

	payload := miniMaxChatRequest{
		Model: p.model,
		Messages: []miniMaxMessage{
			{Role: "system", Content: buildInstructions()},
			{Role: "user", Content: buildInput(question, chunks)},
		},
		TokensToGenerate: 800,
		Temperature:      0.2,
		TopP:             0.95,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/text/chatcompletion_v2", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("llm: minimax api returned %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	var parsed miniMaxChatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}

	text := strings.TrimSpace(parsed.Reply)
	if text == "" && len(parsed.Choices) > 0 {
		text = strings.TrimSpace(parsed.Choices[0].Message.Content)
	}
	if text == "" {
		return "", errors.New("llm: minimax api returned empty output")
	}

	return text, nil
}

type miniMaxChatRequest struct {
	Model            string           `json:"model"`
	Messages         []miniMaxMessage `json:"messages"`
	TokensToGenerate int              `json:"tokens_to_generate,omitempty"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty"`
}

type miniMaxMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type miniMaxChatResponse struct {
	Reply   string          `json:"reply"`
	Choices []miniMaxChoice `json:"choices"`
}

type miniMaxChoice struct {
	Message miniMaxMessage `json:"message"`
}
