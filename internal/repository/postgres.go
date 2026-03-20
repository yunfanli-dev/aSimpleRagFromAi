package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateKnowledgeBase(ctx context.Context, input domain.CreateKnowledgeBaseInput) (domain.KnowledgeBase, error) {
	const query = `
		INSERT INTO knowledge_bases (name, description, status)
		VALUES ($1, $2, 'active')
		RETURNING id::text, name, description, status
	`

	var kb domain.KnowledgeBase
	err := r.pool.QueryRow(ctx, query, input.Name, input.Description).Scan(
		&kb.ID,
		&kb.Name,
		&kb.Description,
		&kb.Status,
	)
	return kb, err
}

func (r *PostgresRepository) GetKnowledgeBase(ctx context.Context, id string) (domain.KnowledgeBase, error) {
	const query = `
		SELECT id::text, name, description, status
		FROM knowledge_bases
		WHERE id = $1
	`

	var kb domain.KnowledgeBase
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&kb.ID,
		&kb.Name,
		&kb.Description,
		&kb.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.KnowledgeBase{}, ErrNotFound
	}
	return kb, err
}

func (r *PostgresRepository) ListKnowledgeBases(ctx context.Context) ([]domain.KnowledgeBase, error) {
	const query = `
		SELECT id::text, name, description, status
		FROM knowledge_bases
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.KnowledgeBase, 0)
	for rows.Next() {
		var kb domain.KnowledgeBase
		if err := rows.Scan(&kb.ID, &kb.Name, &kb.Description, &kb.Status); err != nil {
			return nil, err
		}
		items = append(items, kb)
	}

	return items, rows.Err()
}

func (r *PostgresRepository) CreateDocument(ctx context.Context, input domain.CreateDocumentInput) (domain.Document, error) {
	const query = `
		INSERT INTO documents (
			knowledge_base_id,
			source_type,
			title,
			storage_path,
			content_hash,
			status
		)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		RETURNING id::text, knowledge_base_id::text, title, source_type, status
	`

	var doc domain.Document
	err := r.pool.QueryRow(
		ctx,
		query,
		input.KnowledgeBaseID,
		input.SourceType,
		input.Title,
		input.StoragePath,
		input.ContentHash,
	).Scan(
		&doc.ID,
		&doc.KnowledgeBaseID,
		&doc.Title,
		&doc.SourceType,
		&doc.Status,
	)
	return doc, err
}

func (r *PostgresRepository) ListDocuments(ctx context.Context, kbID string) ([]domain.Document, error) {
	const query = `
		SELECT id::text, knowledge_base_id::text, title, source_type, status
		FROM documents
		WHERE knowledge_base_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, kbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Document, 0)
	for rows.Next() {
		var doc domain.Document
		if err := rows.Scan(&doc.ID, &doc.KnowledgeBaseID, &doc.Title, &doc.SourceType, &doc.Status); err != nil {
			return nil, err
		}
		items = append(items, doc)
	}

	return items, rows.Err()
}
