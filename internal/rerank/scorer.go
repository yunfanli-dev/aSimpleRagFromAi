package rerank

import (
	"sort"
	"strings"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

type Candidate struct {
	Chunk domain.RetrievedChunk
	Score float64
}

// Rank applies the current lightweight heuristic rerank to retrieved chunks.
func Rank(question string, chunks []domain.RetrievedChunk) []domain.RetrievedChunk {
	// TODO: replace this heuristic scorer with a real rerank provider.
	if len(chunks) == 0 {
		return nil
	}

	terms := tokenSet(question)
	candidates := make([]Candidate, 0, len(chunks))
	for _, chunk := range chunks {
		score := chunk.Score
		score += lexicalOverlapScore(terms, chunk.DocumentTitle) * 0.35
		score += lexicalOverlapScore(terms, chunk.Content) * 0.55
		score += sourceBonus(chunk.RetrievalSource)
		score += chunk.QualityScore * 0.25

		candidates = append(candidates, Candidate{
			Chunk: chunk,
			Score: score,
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			if candidates[i].Chunk.DocumentTitle == candidates[j].Chunk.DocumentTitle {
				return candidates[i].Chunk.ChunkIndex < candidates[j].Chunk.ChunkIndex
			}
			return candidates[i].Chunk.DocumentTitle < candidates[j].Chunk.DocumentTitle
		}
		return candidates[i].Score > candidates[j].Score
	})

	ranked := make([]domain.RetrievedChunk, 0, len(candidates))
	for _, item := range candidates {
		item.Chunk.Score = item.Score
		ranked = append(ranked, item.Chunk)
	}

	return ranked
}

// tokenSet tokenizes text into a de-duplicated lowercase term set.
func tokenSet(text string) map[string]struct{} {
	parts := strings.Fields(strings.ToLower(text))
	items := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = strings.Trim(part, ".,!?;:\"'()[]{}")
		if len(part) < 2 {
			continue
		}
		items[part] = struct{}{}
	}
	return items
}

// lexicalOverlapScore measures how much text overlaps with query terms.
func lexicalOverlapScore(queryTerms map[string]struct{}, text string) float64 {
	if len(queryTerms) == 0 {
		return 0
	}

	matches := 0
	seen := make(map[string]struct{}, len(queryTerms))
	for _, token := range strings.Fields(strings.ToLower(text)) {
		token = strings.Trim(token, ".,!?;:\"'()[]{}")
		if _, ok := queryTerms[token]; !ok {
			continue
		}
		if _, duplicated := seen[token]; duplicated {
			continue
		}
		seen[token] = struct{}{}
		matches++
	}

	return float64(matches) / float64(len(queryTerms))
}

// sourceBonus gives small preference to hybrid hits over single-source hits.
func sourceBonus(source string) float64 {
	switch source {
	case "hybrid":
		return 0.2
	case "keyword":
		return 0.1
	case "vector":
		return 0.05
	default:
		return 0
	}
}
