package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
)

type Storage interface {
	SaveUsername(ctx context.Context, id int, username string) error
	DeleteUsername(ctx context.Context, id int) error
	ChangeUsername(ctx context.Context, id int, username string) error
}

type Server struct {
	httpServer *http.Server
	storage    Storage
}

func NewServer(cfg config.Server, storage Storage) *Server {
	s := &Server{storage: storage}

	router := bunrouter.New().Compat()
	router.POST("/notifications/:id", s.notificationsOn)
	router.DELETE("/notifications/:id", s.notificationsOff)
	router.PATCH("/notifications/:id", s.changeUsername)

	s.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s
}

func (s *Server) Run() {
	log.Println("server started") // TODO: логгер

	go func() {
		err := s.httpServer.ListenAndServe()
		log.Printf("http server stopped: %s", err.Error()) // TODO: логгер
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
