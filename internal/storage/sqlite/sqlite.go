package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"web/internal/config"
	"web/internal/storage"
)

type StoreSqlite struct {
	DataBase *sql.DB
	Log      *slog.Logger
}

// Connect connect to database
func (s *StoreSqlite) Connect(cfg *config.Config, log *slog.Logger) storage.Storage {
	const op = "sqlite.Connect"

	db, err := sql.Open("sqlite3", cfg.DatabaseConfig.Config["storagePath"])
	if err != nil {
		log.Error(fmt.Sprintf("%v: %v", op, err.Error()))
		os.Exit(1)
	}
	return &StoreSqlite{DataBase: db, Log: log}
}

type ErrorSqlite struct {
	Code    int
	Message string
}

func ErrorSqliteNew(code int, message string) *ErrorSqlite {
	return &ErrorSqlite{Code: code, Message: message}
}

func (e *ErrorSqlite) Error() string {
	return e.Message
}

func (e *ErrorSqlite) GetCode() int {
	return e.Code
}
