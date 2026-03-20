package service

import (
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

type DocumentService struct {
	repo *repository.MemoryRepository
}

func NewDocumentService(repo *repository.MemoryRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) List(kbID string) []domain.Document {
	return s.repo.ListDocuments(kbID)
}
