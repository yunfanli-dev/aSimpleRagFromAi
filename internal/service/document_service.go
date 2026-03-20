package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/ingest"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

var errUnsupportedSourceType = errors.New("unsupported source_type: only txt and markdown are supported in current ingest")
var errEmptyContent = errors.New("content is empty after normalization")

type DocumentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) Create(ctx context.Context, input domain.CreateDocumentInput) (domain.DocumentIngestResult, error) {
	if input.SourceType != "txt" && input.SourceType != "markdown" {
		return domain.DocumentIngestResult{}, errUnsupportedSourceType
	}

	normalized := ingest.NormalizeText(input.Content)
	if normalized == "" {
		return domain.DocumentIngestResult{}, errEmptyContent
	}

	input.Content = normalized
	if strings.TrimSpace(input.ContentHash) == "" {
		sum := sha256.Sum256([]byte(normalized))
		input.ContentHash = hex.EncodeToString(sum[:])
	}

	chunks := ingest.SplitText(normalized)
	if len(chunks) == 0 {
		return domain.DocumentIngestResult{}, errEmptyContent
	}

	doc, err := s.repo.CreateDocumentWithChunks(ctx, input, chunks)
	if err != nil {
		return domain.DocumentIngestResult{}, err
	}

	return domain.DocumentIngestResult{
		Document:   doc,
		ChunkCount: len(chunks),
	}, nil
}

func (s *DocumentService) List(ctx context.Context, kbID string) ([]domain.Document, error) {
	return s.repo.ListDocuments(ctx, kbID)
}
