package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/config"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/handler"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/observability"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
)

type Handlers struct {
	Health        *handler.HealthHandler
	KnowledgeBase *handler.KnowledgeBaseHandler
	Document      *handler.DocumentHandler
	Query         *handler.QueryHandler
}

func Run() error {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.PostgresDSN)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		return err
	}

	repo := repository.NewPostgresRepository(pool)

	handlers := Handlers{
		Health:        handler.NewHealthHandler(service.NewHealthService()),
		KnowledgeBase: handler.NewKnowledgeBaseHandler(service.NewKnowledgeBaseService(repo)),
		Document:      handler.NewDocumentHandler(service.NewDocumentService(repo)),
		Query:         handler.NewQueryHandler(service.NewQueryService(repo)),
	}

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      NewRouter(handlers),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		observability.Printf("starting %s on %s", cfg.AppName, cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-stopCh:
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		return server.Shutdown(ctx)
	}
}
