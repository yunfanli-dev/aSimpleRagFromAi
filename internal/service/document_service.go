package service

import (
	"context"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

type DocumentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) Create(ctx context.Context, input domain.CreateDocumentInput) (domain.Document, error) {
	return s.repo.CreateDocument(ctx, input)
}

func (s *DocumentService) List(ctx context.Context, kbID string) ([]domain.Document, error) {
	return s.repo.ListDocuments(ctx, kbID)
}
