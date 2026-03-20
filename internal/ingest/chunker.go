package ingest

import (
	"strings"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

const (
	defaultChunkSize    = 900
	defaultChunkOverlap = 120
)

func NormalizeText(input string) string {
	lines := strings.Split(input, "\n")
	blocks := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(blocks) == 0 || blocks[len(blocks)-1] == "" {
				continue
			}
			blocks = append(blocks, "")
			continue
		}
		blocks = append(blocks, strings.Join(strings.Fields(trimmed), " "))
	}

	return strings.TrimSpace(strings.Join(blocks, "\n"))
}

func SplitText(input string) []domain.CreateChunkInput {
	normalized := NormalizeText(input)
	if normalized == "" {
		return nil
	}

	parts := strings.Split(normalized, "\n\n")
	rawChunks := make([]string, 0)
	var current strings.Builder

	flush := func() {
		if current.Len() == 0 {
			return
		}
		rawChunks = append(rawChunks, strings.TrimSpace(current.String()))
		current.Reset()
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if current.Len() == 0 {
			current.WriteString(part)
			continue
		}

		if current.Len()+2+len(part) <= defaultChunkSize {
			current.WriteString("\n\n")
			current.WriteString(part)
			continue
		}

		flush()
		current.WriteString(part)
	}
	flush()

	items := make([]domain.CreateChunkInput, 0, len(rawChunks))
	for idx, chunk := range rawChunks {
		if idx > 0 && defaultChunkOverlap > 0 {
			prev := rawChunks[idx-1]
			overlap := tail(prev, defaultChunkOverlap)
			if overlap != "" && !strings.HasPrefix(chunk, overlap) {
				chunk = overlap + "\n" + chunk
			}
		}

		items = append(items, domain.CreateChunkInput{
			ChunkIndex: idx,
			Content:    chunk,
			TokenCount: len(strings.Fields(chunk)),
		})
	}

	return items
}

func tail(input string, size int) string {
	if len(input) <= size {
		return input
	}
	return strings.TrimSpace(input[len(input)-size:])
}
