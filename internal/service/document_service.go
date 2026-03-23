package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/embedding"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/ingest"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
)

var errUnsupportedSourceType = errors.New("unsupported source_type: only txt and markdown are supported in current ingest")
var errEmptyContent = errors.New("content is empty after normalization")

type DocumentService struct {
	repo     repository.DocumentRepository
	embedder embedding.Provider
}

// NewDocumentService wires ingest and embedding dependencies into the document service.
func NewDocumentService(repo repository.DocumentRepository, embedder embedding.Provider) *DocumentService {
	return &DocumentService{repo: repo, embedder: embedder}
}

// Create validates source content, chunks it, and persists the document with chunks.
func (s *DocumentService) Create(ctx context.Context, input domain.CreateDocumentInput) (domain.DocumentIngestResult, error) {
	// TODO: extend ingest to support PDF and other planned source types.
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

// List returns documents under the given knowledge base.
func (s *DocumentService) List(ctx context.Context, kbID string) ([]domain.Document, error) {
	return s.repo.ListDocuments(ctx, kbID)
}

// Get loads a single document by ID.
func (s *DocumentService) Get(ctx context.Context, id string) (domain.Document, error) {
	return s.repo.GetDocument(ctx, id)
}

// ListChunks returns stored chunks for one document.
func (s *DocumentService) ListChunks(ctx context.Context, documentID string) ([]domain.Chunk, error) {
	return s.repo.ListChunks(ctx, documentID)
}

// Reindex regenerates embeddings for all chunks in one document.
func (s *DocumentService) Reindex(ctx context.Context, documentID string) (domain.ReindexDocumentResult, error) {
	document, err := s.repo.GetDocument(ctx, documentID)
	if err != nil {
		return domain.ReindexDocumentResult{}, err
	}

	chunks, err := s.repo.ListChunkEmbeddingsInput(ctx, documentID)
	if err != nil {
		return domain.ReindexDocumentResult{}, err
	}

	result := domain.ReindexDocumentResult{
		DocumentID:      document.ID,
		ChunkCount:      len(chunks),
		EmbeddingModel:  s.embedder.Model(),
		EmbeddingStatus: "completed",
	}

	for _, chunk := range chunks {
		vector, err := s.embedder.Embed(ctx, chunk.Content)
		if err != nil {
			return domain.ReindexDocumentResult{}, err
		}

		if err := s.repo.UpsertChunkEmbedding(ctx, chunk.ChunkID, vector, s.embedder.Model()); err != nil {
			return domain.ReindexDocumentResult{}, err
		}
		result.EmbeddedCount++
	}

	return result, nil
}
