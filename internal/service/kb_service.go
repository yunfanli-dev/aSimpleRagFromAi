package service

import (
	"context"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

type KnowledgeBaseService struct {
	repo repository.KnowledgeBaseRepository
}

// NewKnowledgeBaseService builds a service around knowledge base persistence.
func NewKnowledgeBaseService(repo repository.KnowledgeBaseRepository) *KnowledgeBaseService {
	return &KnowledgeBaseService{repo: repo}
}

// Create persists a new knowledge base record.
func (s *KnowledgeBaseService) Create(ctx context.Context, input domain.CreateKnowledgeBaseInput) (domain.KnowledgeBase, error) {
	return s.repo.CreateKnowledgeBase(ctx, input)
}

// Get loads a single knowledge base by ID.
func (s *KnowledgeBaseService) Get(ctx context.Context, id string) (domain.KnowledgeBase, error) {
	return s.repo.GetKnowledgeBase(ctx, id)
}

// List returns all knowledge bases ordered by the repository.
func (s *KnowledgeBaseService) List(ctx context.Context) ([]domain.KnowledgeBase, error) {
	return s.repo.ListKnowledgeBases(ctx)
}
