package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"os"
	"web/internal/config"
	"web/internal/server/server/interfaces"
	"web/internal/storage"
	"web/storage/tags-list"
)

type Server struct {
	Router   chi.Router
	Handlers handlerInterfaces.HandlerMethods
	Db       storage.Storage
	Log      *slog.Logger
	AllTags  *tagsList.TagsList
}

// NewServer create new http server
func NewServer(db *storage.Storage, log *slog.Logger) *Server {
	allTags := tagsList.NewTagsMemoryList(db, log)

	return &Server{
		Router:  chi.NewRouter(),
		Db:      *db,
		Log:     log,
		AllTags: allTags,
	}
}

func (s *Server) InitHandlers(handlers handlerInterfaces.HandlerMethods) {
	s.Handlers = handlers
}

// Start start http server
func (s *Server) Start(cfg *config.Config, log *slog.Logger) {
	const op = "httpserver.Server.Start"

	log.Info(fmt.Sprintf("Starting server on %s:%s", cfg.Host, cfg.Port))
	err := http.ListenAndServe(cfg.Host+":"+cfg.Port, s.Router)
	if err != nil {
		log.Error(fmt.Sprintf("%v: %v", op, err.Error()))
		os.Exit(1)
	}

}

// AddRoute adds a new path to the Server mux.
//func (s *Server) AddRoute(pattern string, fn http.HandlerFunc) {
//
//}
