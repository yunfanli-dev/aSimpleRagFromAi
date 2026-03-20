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

func (r *PostgresRepository) CreateDocumentWithChunks(ctx context.Context, input domain.CreateDocumentInput, chunks []domain.CreateChunkInput) (domain.Document, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Document{}, err
	}
	defer tx.Rollback(ctx)

	const insertDocument = `
		INSERT INTO documents (
			knowledge_base_id,
			source_type,
			title,
			storage_path,
			content_hash,
			content_text,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
		RETURNING id::text, knowledge_base_id::text, title, source_type, status
	`

	var doc domain.Document
	err = tx.QueryRow(
		ctx,
		insertDocument,
		input.KnowledgeBaseID,
		input.SourceType,
		input.Title,
		input.StoragePath,
		input.ContentHash,
		input.Content,
	).Scan(
		&doc.ID,
		&doc.KnowledgeBaseID,
		&doc.Title,
		&doc.SourceType,
		&doc.Status,
	)
	if err != nil {
		return domain.Document{}, err
	}

	const insertChunk = `
		INSERT INTO chunks (document_id, chunk_index, content, token_count, metadata_json)
		VALUES ($1, $2, $3, $4, '{}'::jsonb)
	`
	for _, chunk := range chunks {
		if _, err := tx.Exec(ctx, insertChunk, doc.ID, chunk.ChunkIndex, chunk.Content, chunk.TokenCount); err != nil {
			return domain.Document{}, err
		}
	}

	if _, err := tx.Exec(ctx, `UPDATE documents SET status = 'indexed', updated_at = NOW() WHERE id = $1`, doc.ID); err != nil {
		return domain.Document{}, err
	}
	doc.Status = "indexed"

	if err := tx.Commit(ctx); err != nil {
		return domain.Document{}, err
	}

	return doc, nil
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

func (r *PostgresRepository) GetDocument(ctx context.Context, id string) (domain.Document, error) {
	const query = `
		SELECT
			id::text,
			knowledge_base_id::text,
			title,
			source_type,
			status,
			storage_path,
			content_text
		FROM documents
		WHERE id = $1
	`

	var doc domain.Document
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&doc.ID,
		&doc.KnowledgeBaseID,
		&doc.Title,
		&doc.SourceType,
		&doc.Status,
		&doc.StoragePath,
		&doc.Content,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Document{}, ErrNotFound
	}
	return doc, err
}

func (r *PostgresRepository) ListChunks(ctx context.Context, documentID string) ([]domain.Chunk, error) {
	const query = `
		SELECT id::text, document_id::text, chunk_index, content, token_count
		FROM chunks
		WHERE document_id = $1
		ORDER BY chunk_index ASC
	`

	rows, err := r.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Chunk, 0)
	for rows.Next() {
		var chunk domain.Chunk
		if err := rows.Scan(
			&chunk.ID,
			&chunk.DocumentID,
			&chunk.ChunkIndex,
			&chunk.Content,
			&chunk.TokenCount,
		); err != nil {
			return nil, err
		}
		items = append(items, chunk)
	}

	return items, rows.Err()
}

func (r *PostgresRepository) SearchChunks(ctx context.Context, kbID, question string, limit int) ([]domain.RetrievedChunk, error) {
	if limit <= 0 || question == "" {
		return []domain.RetrievedChunk{}, nil
	}

	const query = `
		WITH ranked_chunks AS (
			SELECT
				c.id::text,
				c.document_id::text,
				d.title,
				c.chunk_index,
				c.content,
				ts_rank_cd(
					setweight(to_tsvector('simple', coalesce(d.title, '')), 'A') ||
					setweight(c.tsv, 'B'),
					websearch_to_tsquery('simple', $2)
				) AS score
			FROM chunks c
			INNER JOIN documents d ON d.id = c.document_id
			WHERE d.knowledge_base_id = $1
			  AND (
				setweight(to_tsvector('simple', coalesce(d.title, '')), 'A') ||
				setweight(c.tsv, 'B')
			  ) @@ websearch_to_tsquery('simple', $2)
		)
		SELECT
			id,
			document_id,
			title,
			chunk_index,
			content,
			score
		FROM ranked_chunks
		ORDER BY score DESC, title ASC, chunk_index ASC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, kbID, question, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.RetrievedChunk, 0, limit)
	for rows.Next() {
		var item domain.RetrievedChunk
		if err := rows.Scan(
			&item.ChunkID,
			&item.DocumentID,
			&item.DocumentTitle,
			&item.ChunkIndex,
			&item.Content,
			&item.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PostgresRepository) LogQuery(ctx context.Context, input domain.QueryLogInput) error {
	const query = `
		INSERT INTO query_logs (
			knowledge_base_id,
			question,
			answer,
			latency_ms,
			retrieved_chunk_ids
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			COALESCE((SELECT array_agg(value::uuid) FROM unnest($5::text[]) AS value), '{}'::uuid[])
		)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		input.KnowledgeBaseID,
		input.Question,
		input.Answer,
		input.LatencyMS,
		input.RetrievedChunkIDs,
	)
	return err
}
