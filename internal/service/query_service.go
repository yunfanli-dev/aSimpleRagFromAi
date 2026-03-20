package service

import "github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"

type QueryService struct{}

func NewQueryService() *QueryService {
	return &QueryService{}
}

func (s *QueryService) Ask(req domain.QueryRequest) domain.QueryResponse {
	resp := domain.QueryResponse{
		Answer: "RAG pipeline scaffold is ready. Retrieval and generation will be implemented in the next phase.",
		Citations: []domain.Citation{
			{
				ChunkID: "chunk_demo_1",
				Text:    "SimpleRAG-Go bootstrap response.",
				Source:  "system/demo",
			},
		},
	}
	if req.Debug {
		resp.DebugInfo = map[string]any{
			"query":       req.Question,
			"knowledgeId": req.KnowledgeBaseID,
			"stage":       "bootstrap",
		}
	}
	return resp
}
