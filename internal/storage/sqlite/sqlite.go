package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"web/internal/config"
	"web/internal/storage"
	"web/internal/storage/sqlite/methods/task"
)

type StoreSqlite struct {
	DataBase *sql.DB
	Log      *slog.Logger
	// methods responsible for tasks and tags
	// methods responsible for user
	*repository.TaskTagMethods
}

// Connect connect to database
func (s *StoreSqlite) Connect(cfg *config.Config, log *slog.Logger) storage.Storage {
	const op = "sqlite.Connect"

	db, err := sql.Open("sqlite3", cfg.DatabaseConfig.Config["storagePath"])
	if err != nil {
		log.Error(fmt.Sprintf("%v: %v", op, err.Error()))
		os.Exit(1)
	}
	return &StoreSqlite{DataBase: db, Log: log, TaskTagMethods: repository.TaskTagMethodsNew(db, log)}
}
