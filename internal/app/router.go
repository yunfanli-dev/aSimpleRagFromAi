package app

import "github.com/gin-gonic/gin"

func NewRouter(handlers Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", handlers.Health.Health)
	router.GET("/readyz", handlers.Health.Ready)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/kbs", handlers.KnowledgeBase.List)
		v1.POST("/kbs/:id/documents", handlers.Document.Create)
		v1.GET("/kbs/:id/documents", handlers.Document.List)
		v1.POST("/documents/:id/reindex", handlers.Document.Reindex)
		v1.POST("/query", handlers.Query.Query)
		v1.POST("/query/debug", handlers.Query.Debug)
	}

	return router
}
