package domain

type KnowledgeBase struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Document struct {
	ID              string `json:"id"`
	KnowledgeBaseID string `json:"knowledge_base_id"`
	Title           string `json:"title"`
	SourceType      string `json:"source_type"`
	Status          string `json:"status"`
}

type QueryRequest struct {
	KnowledgeBaseID string `json:"knowledge_base_id" binding:"required"`
	Question        string `json:"question" binding:"required"`
	Debug           bool   `json:"debug"`
}

type Citation struct {
	ChunkID string `json:"chunk_id"`
	Text    string `json:"text"`
	Source  string `json:"source"`
}

type QueryResponse struct {
	Answer    string     `json:"answer"`
	Citations []Citation `json:"citations"`
	DebugInfo any        `json:"debug_info,omitempty"`
}
