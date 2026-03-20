CREATE INDEX IF NOT EXISTS idx_documents_kb_id ON documents (knowledge_base_id);
CREATE INDEX IF NOT EXISTS idx_documents_status ON documents (status);

CREATE INDEX IF NOT EXISTS idx_chunks_document_id ON chunks (document_id);
CREATE INDEX IF NOT EXISTS idx_chunks_tsv ON chunks USING GIN (tsv);

CREATE INDEX IF NOT EXISTS idx_chunk_vectors_embedding
    ON chunk_vectors
    USING hnsw (embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_ingest_jobs_status_scheduled
    ON ingest_jobs (status, scheduled_at);

CREATE INDEX IF NOT EXISTS idx_query_logs_kb_id_created_at
    ON query_logs (knowledge_base_id, created_at DESC);
