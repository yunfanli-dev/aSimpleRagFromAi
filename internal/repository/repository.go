package repository

import (
	"context"
	"errors"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

var ErrNotFound = errors.New("repository: not found")

type KnowledgeBaseRepository interface {
	CreateKnowledgeBase(ctx context.Context, input domain.CreateKnowledgeBaseInput) (domain.KnowledgeBase, error)
	GetKnowledgeBase(ctx context.Context, id string) (domain.KnowledgeBase, error)
	ListKnowledgeBases(ctx context.Context) ([]domain.KnowledgeBase, error)
}

type DocumentRepository interface {
	CreateDocumentWithChunks(ctx context.Context, input domain.CreateDocumentInput, chunks []domain.CreateChunkInput) (domain.Document, error)
	ListDocuments(ctx context.Context, kbID string) ([]domain.Document, error)
}
