package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

const (
	defaultRetrieveLimit = 5
	maxAnswerLength      = 1200
	maxCitationTextRunes = 240
)

type QueryService struct {
	repo repository.QueryRepository
}

func NewQueryService(repo repository.QueryRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) Ask(ctx context.Context, req domain.QueryRequest) (domain.QueryResponse, error) {
	start := time.Now()
	question := normalizeQuestion(req.Question)

	chunks, err := s.repo.SearchChunks(ctx, req.KnowledgeBaseID, question, defaultRetrieveLimit)
	if err != nil {
		return domain.QueryResponse{}, err
	}

	resp := domain.QueryResponse{
		Answer:    buildAnswer(chunks),
		Citations: buildCitations(chunks),
	}

	if req.Debug {
		resp.DebugInfo = map[string]any{
			"query":            question,
			"knowledge_base":   req.KnowledgeBaseID,
			"retrieved_count":  len(chunks),
			"retrieved_chunks": chunks,
			"latency_ms":       time.Since(start).Milliseconds(),
		}
	}

	if err := s.repo.LogQuery(ctx, domain.QueryLogInput{
		KnowledgeBaseID:   req.KnowledgeBaseID,
		Question:          question,
		Answer:            resp.Answer,
		LatencyMS:         int(time.Since(start).Milliseconds()),
		RetrievedChunkIDs: chunkIDs(chunks),
	}); err != nil {
		return domain.QueryResponse{}, err
	}

	return resp, nil
}

func buildAnswer(chunks []domain.RetrievedChunk) string {
	if len(chunks) == 0 {
		return "I couldn't find relevant content in the current knowledge base."
	}

	parts := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		excerpt := clipText(chunk.Content, maxCitationTextRunes)
		parts = append(parts, chunk.DocumentTitle+" [chunk "+formatChunkIndex(chunk.ChunkIndex)+"]: "+excerpt)
	}

	answer := strings.Join(parts, "\n\n")
	answerRunes := []rune(answer)
	if len(answerRunes) > maxAnswerLength {
		return strings.TrimSpace(string(answerRunes[:maxAnswerLength])) + "..."
	}
	return answer
}

func buildCitations(chunks []domain.RetrievedChunk) []domain.Citation {
	items := make([]domain.Citation, 0, len(chunks))
	for _, chunk := range chunks {
		items = append(items, domain.Citation{
			ChunkID:       chunk.ChunkID,
			DocumentID:    chunk.DocumentID,
			DocumentTitle: chunk.DocumentTitle,
			ChunkIndex:    chunk.ChunkIndex,
			Text:          clipText(chunk.Content, maxCitationTextRunes),
			Source:        chunk.DocumentTitle,
			Score:         chunk.Score,
		})
	}
	return items
}

func chunkIDs(chunks []domain.RetrievedChunk) []string {
	items := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		items = append(items, chunk.ChunkID)
	}
	return items
}

func normalizeQuestion(question string) string {
	return strings.Join(strings.Fields(question), " ")
}

func clipText(text string, limit int) string {
	trimmed := strings.TrimSpace(text)
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
