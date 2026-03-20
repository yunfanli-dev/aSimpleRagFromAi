package llm

import (
	"context"
	"strconv"
	"strings"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

type Provider interface {
	Model() string
	Generate(ctx context.Context, question string, chunks []domain.RetrievedChunk) (string, error)
}

type ExtractiveProvider struct {
	model string
}

func NewExtractiveProvider(model string) *ExtractiveProvider {
	if strings.TrimSpace(model) == "" {
		model = "local-extractive-v1"
	}

	return &ExtractiveProvider{model: model}
}

func (p *ExtractiveProvider) Model() string {
	return p.model
}

func (p *ExtractiveProvider) Generate(_ context.Context, question string, chunks []domain.RetrievedChunk) (string, error) {
	if len(chunks) == 0 {
		return "I couldn't find relevant content in the current knowledge base.", nil
	}

	lines := make([]string, 0, len(chunks)+1)
	lines = append(lines, "Answer based on retrieved context for: "+question)
	for _, chunk := range chunks {
		lines = append(lines, chunk.DocumentTitle+" [chunk "+formatChunkIndex(chunk.ChunkIndex)+"]: "+clipText(chunk.Content, 220))
	}

	answer := strings.Join(lines, "\n\n")
	answerRunes := []rune(answer)
	if len(answerRunes) > 1200 {
		return strings.TrimSpace(string(answerRunes[:1200])) + "...", nil
	}

	return answer, nil
}

func clipText(text string, limit int) string {
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
	if trimmed == "" {
		return trimmed
	}

	runes := []rune(trimmed)
	if len(runes) <= limit {
		return trimmed
	}

	return strings.TrimSpace(string(runes[:limit])) + "..."
}

func formatChunkIndex(index int) string {
	return strconv.Itoa(index)
}
