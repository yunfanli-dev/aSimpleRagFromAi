package service

import (
	"context"
	"strings"
	"time"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

type QueryService struct {
	repo repository.QueryRepository
}

func NewQueryService(repo repository.QueryRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) Ask(ctx context.Context, req domain.QueryRequest) (domain.QueryResponse, error) {
	start := time.Now()

	chunks, err := s.repo.SearchChunks(ctx, req.KnowledgeBaseID, req.Question, 5)
	if err != nil {
		return domain.QueryResponse{}, err
	}

	resp := domain.QueryResponse{
		Answer:    buildAnswer(chunks),
		Citations: buildCitations(chunks),
	}

	if req.Debug {
		resp.DebugInfo = map[string]any{
			"query":           req.Question,
			"knowledge_base":  req.KnowledgeBaseID,
			"retrieved_count": len(chunks),
			"retrieved_chunks": chunks,
			"latency_ms":      time.Since(start).Milliseconds(),
		}
	}

	if err := s.repo.LogQuery(ctx, domain.QueryLogInput{
		KnowledgeBaseID:   req.KnowledgeBaseID,
		Question:          req.Question,
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
		parts = append(parts, chunk.Content)
	}

	answer := strings.Join(parts, "\n\n")
	if len(answer) > 1200 {
		return strings.TrimSpace(answer[:1200]) + "..."
	}
	return answer
}

func buildCitations(chunks []domain.RetrievedChunk) []domain.Citation {
	items := make([]domain.Citation, 0, len(chunks))
	for _, chunk := range chunks {
		items = append(items, domain.Citation{
			ChunkID: chunk.ChunkID,
			Text:    chunk.Content,
			Source:  chunk.DocumentTitle,
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
