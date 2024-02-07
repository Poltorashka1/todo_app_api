package repository

import (
	"database/sql"
	"log/slog"
)

type TaskTagMethods struct {
	DataBase *sql.DB
	Log      *slog.Logger
}

func TaskTagMethodsNew(db *sql.DB, log *slog.Logger) *TaskTagMethods {
	return &TaskTagMethods{
		// todo check tihs
		DataBase: db,
		Log:      log,
	}
}

func (t *TaskTagMethods) GetDb() *sql.DB {
	return t.DataBase

}
