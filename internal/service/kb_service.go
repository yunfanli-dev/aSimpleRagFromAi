package service

import (
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

type KnowledgeBaseService struct {
	repo *repository.MemoryRepository
}

func NewKnowledgeBaseService(repo *repository.MemoryRepository) *KnowledgeBaseService {
	return &KnowledgeBaseService{repo: repo}
}

func (s *KnowledgeBaseService) List() []domain.KnowledgeBase {
	return s.repo.ListKnowledgeBases()
}
