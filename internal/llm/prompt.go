package llm

import (
	"fmt"
	"strings"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

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
