package service

import (
	"context"
	"strings"
	"time"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/embedding"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/llm"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/rerank"
)

const (
	defaultRetrieveLimit   = 5
	candidateRetrieveLimit = 12
	maxCitationTextRunes   = 240
)

type QueryService struct {
	repo     repository.QueryRepository
	embedder embedding.Provider
	llm      llm.Provider
}

// NewQueryService wires retrieval, embedding, and generation dependencies.
func NewQueryService(repo repository.QueryRepository, embedder embedding.Provider, llmProvider llm.Provider) *QueryService {
	return &QueryService{repo: repo, embedder: embedder, llm: llmProvider}
}

// Ask executes retrieval, rerank, answer generation, and query logging for one request.
func (s *QueryService) Ask(ctx context.Context, req domain.QueryRequest) (domain.QueryResponse, error) {
	start := time.Now()
	question := normalizeQuestion(req.Question)
	questionEmbedding, err := s.embedQuestion(ctx, question)
	if err != nil {
		return domain.QueryResponse{}, err
	}

	chunks, err := s.repo.SearchChunks(ctx, req.KnowledgeBaseID, question, questionEmbedding, candidateRetrieveLimit)
	if err != nil {
		return domain.QueryResponse{}, err
	}
	chunks = selectTopChunks(question, chunks, defaultRetrieveLimit)
	answer, err := s.llm.Generate(ctx, question, chunks)
	if err != nil {
		return domain.QueryResponse{}, err
	}

	resp := domain.QueryResponse{
		Answer:    answer,
		Citations: buildCitations(chunks),
		Model:     s.llm.Model(),
	}

	if req.Debug {
		resp.DebugInfo = map[string]any{
			"query":            question,
			"knowledge_base":   req.KnowledgeBaseID,
			"retrieved_count":  len(chunks),
			"retrieved_chunks": chunks,
			"embedding_model":  s.embedder.Model(),
			"llm_model":        s.llm.Model(),
			"latency_ms":       time.Since(start).Milliseconds(),
		}
	}

	if err := s.repo.LogQuery(ctx, domain.QueryLogInput{
		KnowledgeBaseID:   req.KnowledgeBaseID,
		Question:          question,
		Answer:            answer,
		LatencyMS:         int(time.Since(start).Milliseconds()),
		RetrievedChunkIDs: chunkIDs(chunks),
	}); err != nil {
		return domain.QueryResponse{}, err
	}

	return resp, nil
}

// buildCitations converts retrieved chunks into API citation payloads.
func buildCitations(chunks []domain.RetrievedChunk) []domain.Citation {
	items := make([]domain.Citation, 0, len(chunks))
	for _, chunk := range chunks {
		items = append(items, domain.Citation{
			ChunkID:         chunk.ChunkID,
			DocumentID:      chunk.DocumentID,
			DocumentTitle:   chunk.DocumentTitle,
			ChunkIndex:      chunk.ChunkIndex,
			Text:            clipText(chunk.Content, maxCitationTextRunes),
			Source:          chunk.DocumentTitle,
			Score:           chunk.Score,
			RetrievalSource: chunk.RetrievalSource,
		})
	}
	return items
}

// chunkIDs extracts chunk identifiers for query logging.
func chunkIDs(chunks []domain.RetrievedChunk) []string {
	items := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		items = append(items, chunk.ChunkID)
	}
	return items
}

// normalizeQuestion collapses repeated whitespace in user input.
func normalizeQuestion(question string) string {
	return strings.Join(strings.Fields(question), " ")
}

// clipText trims long text for citation and prompt payloads.
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

// embedQuestion skips empty questions and delegates to the embedding provider.
func (s *QueryService) embedQuestion(ctx context.Context, question string) ([]float32, error) {
	if question == "" {
		return nil, nil
	}

	return s.embedder.Embed(ctx, question)
}

// selectTopChunks applies rerank, de-duplication, and document diversity limits.
func selectTopChunks(question string, chunks []domain.RetrievedChunk, limit int) []domain.RetrievedChunk {
	// TODO: replace heuristic rerank and selection rules with a real rerank provider.
	if len(chunks) == 0 || limit <= 0 {
		return []domain.RetrievedChunk{}
	}

	ranked := rerank.Rank(question, chunks)
	selected := make([]domain.RetrievedChunk, 0, limit)
	documentCounts := make(map[string]int)
	seen := make(map[string]struct{}, len(ranked))

	for _, chunk := range ranked {
		if _, ok := seen[chunk.ChunkID]; ok {
			continue
		}
		if documentCounts[chunk.DocumentID] >= 2 {
			continue
		}
		if hasNearbyChunk(selected, chunk) {
			continue
		}

		seen[chunk.ChunkID] = struct{}{}
		documentCounts[chunk.DocumentID]++
		selected = append(selected, chunk)
		if len(selected) == limit {
			break
		}
	}

	if len(selected) == 0 && len(ranked) > 0 {
		return ranked[:min(limit, len(ranked))]
	}

	return selected
}

// hasNearbyChunk suppresses adjacent chunks from the same document in final citations.
func hasNearbyChunk(selected []domain.RetrievedChunk, candidate domain.RetrievedChunk) bool {
	for _, item := range selected {
		if item.DocumentID != candidate.DocumentID {
			continue
		}
		diff := item.ChunkIndex - candidate.ChunkIndex
		if diff < 0 {
			diff = -diff
		}
		if diff <= 1 {
			return true
		}
	}
	return false
}

// min returns the smaller integer value.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
