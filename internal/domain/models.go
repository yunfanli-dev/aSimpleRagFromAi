package domain

type KnowledgeBase struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type CreateKnowledgeBaseInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type Document struct {
	ID              string `json:"id"`
	KnowledgeBaseID string `json:"knowledge_base_id"`
	Title           string `json:"title"`
	SourceType      string `json:"source_type"`
	Status          string `json:"status"`
	StoragePath     string `json:"storage_path,omitempty"`
	Content         string `json:"content,omitempty"`
}

type CreateDocumentInput struct {
	KnowledgeBaseID string `json:"knowledge_base_id"`
	Title           string `json:"title" binding:"required"`
	SourceType      string `json:"source_type" binding:"required"`
	StoragePath     string `json:"storage_path"`
	ContentHash     string `json:"content_hash"`
	Content         string `json:"content" binding:"required"`
}

type Chunk struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	ChunkIndex int    `json:"chunk_index"`
	Content    string `json:"content"`
	TokenCount int    `json:"token_count"`
}

type CreateChunkInput struct {
	ChunkIndex int
	Content    string
	TokenCount int
}

type DocumentIngestResult struct {
	Document   Document `json:"document"`
	ChunkCount int      `json:"chunk_count"`
}

type QueryRequest struct {
	KnowledgeBaseID string `json:"knowledge_base_id" binding:"required"`
	Question        string `json:"question" binding:"required"`
	Debug           bool   `json:"debug"`
}

type RetrievedChunk struct {
	ChunkID       string  `json:"chunk_id"`
	DocumentID    string  `json:"document_id"`
	DocumentTitle string  `json:"document_title"`
	ChunkIndex    int     `json:"chunk_index"`
	Content       string  `json:"content"`
	Score         float64 `json:"score"`
}

type Citation struct {
	ChunkID       string  `json:"chunk_id"`
	DocumentID    string  `json:"document_id"`
	DocumentTitle string  `json:"document_title"`
	ChunkIndex    int     `json:"chunk_index"`
	Text          string  `json:"text"`
	Source        string  `json:"source"`
	Score         float64 `json:"score,omitempty"`
}

type QueryResponse struct {
	Answer    string     `json:"answer"`
	Citations []Citation `json:"citations"`
	DebugInfo any        `json:"debug_info,omitempty"`
}

type QueryLogInput struct {
	KnowledgeBaseID   string
	Question          string
	Answer            string
	LatencyMS         int
	RetrievedChunkIDs []string
}
