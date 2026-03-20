package repository

import "github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"

type MemoryRepository struct{}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) ListKnowledgeBases() []domain.KnowledgeBase {
	return []domain.KnowledgeBase{
		{
			ID:          "kb_demo",
			Name:        "Demo KB",
			Description: "bootstrap knowledge base",
			Status:      "ready",
		},
	}
}

func (r *MemoryRepository) ListDocuments(kbID string) []domain.Document {
	return []domain.Document{
		{
			ID:              "doc_demo",
			KnowledgeBaseID: kbID,
			Title:           "Getting Started",
			SourceType:      "markdown",
			Status:          "indexed",
		},
	}
}
