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

var ErrMissingAPIKey = errors.New("llm: missing OPENAI_API_KEY for openai provider")

type OpenAIProvider struct {
	model           string
	baseURL         string
	apiKey          string
	reasoningEffort string
	client          *http.Client
}

func NewOpenAIProvider(model, baseURL, apiKey, reasoningEffort string, timeout time.Duration) (*OpenAIProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrMissingAPIKey
	}
	if strings.TrimSpace(model) == "" {
		model = "gpt-5-mini"
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &OpenAIProvider{
		model:           model,
		baseURL:         strings.TrimRight(baseURL, "/"),
		apiKey:          apiKey,
		reasoningEffort: normalizeReasoningEffort(reasoningEffort),
		client:          &http.Client{Timeout: timeout},
	}, nil
}

func (p *OpenAIProvider) Model() string {
	return p.model
}

func (p *OpenAIProvider) Generate(ctx context.Context, question string, chunks []domain.RetrievedChunk) (string, error) {
	if len(chunks) == 0 {
		return "I couldn't find relevant content in the current knowledge base.", nil
	}

	payload := openAIResponsesRequest{
		Model:        p.model,
		Instructions: buildInstructions(),
		Input:        buildInput(question, chunks),
	}
	if p.reasoningEffort != "" {
		payload.Reasoning = &openAIReasoning{Effort: p.reasoningEffort}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/responses", bytes.NewReader(body))
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
		return "", fmt.Errorf("llm: openai responses api returned %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	var parsed openAIResponsesResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}

	text := strings.TrimSpace(parsed.OutputText)
	if text == "" {
		text = strings.TrimSpace(parsed.extractText())
	}
	if text == "" {
		return "", errors.New("llm: openai responses api returned empty output")
	}

	return text, nil
}

type openAIResponsesRequest struct {
	Model        string           `json:"model"`
	Instructions string           `json:"instructions,omitempty"`
	Input        string           `json:"input"`
	Reasoning    *openAIReasoning `json:"reasoning,omitempty"`
}

type openAIReasoning struct {
	Effort string `json:"effort,omitempty"`
}

type openAIResponsesResponse struct {
	OutputText string             `json:"output_text"`
	Output     []openAIOutputItem `json:"output"`
}

type openAIOutputItem struct {
	Type    string              `json:"type"`
	Content []openAIContentPart `json:"content"`
}

type openAIContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (r openAIResponsesResponse) extractText() string {
	parts := make([]string, 0)
	for _, item := range r.Output {
		for _, content := range item.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				parts = append(parts, strings.TrimSpace(content.Text))
			}
		}
	}
	return strings.Join(parts, "\n\n")
}

func buildInstructions() string {
	return strings.Join([]string{
		"You answer strictly from the retrieved context.",
		"If the context is insufficient, say so clearly.",
		"Keep the answer concise and factual.",
		"When useful, mention the source document title in prose.",
	}, " ")
}

func buildInput(question string, chunks []domain.RetrievedChunk) string {
	parts := make([]string, 0, len(chunks)+2)
	parts = append(parts, "Question:\n"+question)
	parts = append(parts, "Retrieved context:")
	for _, chunk := range chunks {
		parts = append(parts, fmt.Sprintf("- %s [chunk %d | source=%s]: %s",
			chunk.DocumentTitle,
			chunk.ChunkIndex,
			chunk.RetrievalSource,
			clipText(chunk.Content, 400),
		))
	}

	return strings.Join(parts, "\n\n")
}

func normalizeReasoningEffort(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "low", "medium", "high":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}
