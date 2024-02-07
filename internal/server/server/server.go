package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"os"
	"web/internal/config"
	taskHandlers "web/internal/server/server/handlers/task_handlers"
	userHandlers "web/internal/server/server/handlers/user_handlers"
	"web/internal/storage"
	"web/storage/tags-list"
)

type Server struct {
	Router       chi.Router
	TaskHandlers *taskHandlers.TaskHandler
	UserHandlers *userHandlers.UserHandler
	Db           storage.Storage
	Log          *slog.Logger
	AllTags      *tagsList.TagsList
}

// NewServer create new http server
func NewServer(db *storage.Storage, log *slog.Logger) *Server {
	allTags := tagsList.NewTagsMemoryList(db, log)
	return &Server{
		Router:       chi.NewRouter(),
		Db:           *db,
		Log:          log,
		AllTags:      allTags,
		TaskHandlers: taskHandlers.NewTaskHandler(db, log, allTags),
		UserHandlers: userHandlers.NewUserHandler(db, log, allTags),
	}
}

// Start start http server
func (s *Server) Start(cfg *config.Config, log *slog.Logger) {
	const op = "httpserver.Server.Start"

	log.Info("Starting server", slog.String("host", cfg.Host), slog.String("port", cfg.Port))
	err := http.ListenAndServeTLS(cfg.Host+":"+cfg.Port, "cert.pem", "key.pem", s.Router)
	if err != nil {
		log.Error(fmt.Sprintf("%v: %v", op, err.Error()))
		os.Exit(1)
	}

}
