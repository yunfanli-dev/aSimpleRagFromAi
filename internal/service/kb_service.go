package service

import (
	"context"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

type KnowledgeBaseService struct {
	repo repository.KnowledgeBaseRepository
}

func NewKnowledgeBaseService(repo repository.KnowledgeBaseRepository) *KnowledgeBaseService {
	return &KnowledgeBaseService{repo: repo}
}

func (s *KnowledgeBaseService) Create(ctx context.Context, input domain.CreateKnowledgeBaseInput) (domain.KnowledgeBase, error) {
	return s.repo.CreateKnowledgeBase(ctx, input)
}

func (s *KnowledgeBaseService) Get(ctx context.Context, id string) (domain.KnowledgeBase, error) {
	return s.repo.GetKnowledgeBase(ctx, id)
}

func (s *KnowledgeBaseService) List(ctx context.Context) ([]domain.KnowledgeBase, error) {
	return s.repo.ListKnowledgeBases(ctx)
}
