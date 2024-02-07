package storage

import (
	"log/slog"
	"time"
	"web/internal/config"
	"web/internal/storage/models"
)

// Storage - database interface
type Storage interface {
	// Connect connect to database.
	// cfg *config.Config - configuration of database connection.
	// log *slog.Logger - logger.
	Connect(cfg *config.Config, log *slog.Logger) Storage

	//CreateTask creates new task with selected parameters.
	CreateTask(text string, tags []string, dueDate *time.Time) error

	// GetTask gets task by ID.
	GetTask(id int) (*models.Task, error)

	// GetTasksByTagShort
	GetTasksByTagShort(tagList []string) (*models.Tasks, error)

	// GetTasksByTagFull
	GetTasksByTagFull(tagList []string) (*models.Tasks, error)

	// Delete deletes task by ID or all tasks.
	// TODO check admin rules
	DeleteTask(id ...string) error

	// GetTag returns tag by name.
	GetTag(name string) (*models.Tag, error)

	// GetAllTags returns all tags.
	GetAllTags() (*models.Tags, error)

	//GetAllTasks returns all tasks.
	GetAllTasks() (*models.Tasks, error)

	// GetTasksByDueDate returns tasks by due date.
	GetTasksByDueDate(due *time.Time) (*models.Tasks, error)

	// CreateTag creates new tag
	CreateTag(name string) error

	// DeleteTag deletes tag
	DeleteTag(name ...string) error

	// GetTasksByDueAndTag returns tasks by due date and tag
	GetTasksByDueAndTagFull(tagList []string, due *time.Time) (*models.Tasks, error)
	GetTasksByDueAndTagShort(tagList []string, due *time.Time) (*models.Tasks, error)

	GetTasksByTag(tagList []string) (*models.Tasks, error)

	GetTasksByTagAndDue(tagList []string, due *time.Time) (*models.Tasks, error)
}

type SqlError interface {
	Error() string
	GetCode() int
}

type ErrorSql struct {
	Code    int
	Message string
}

func ErrorSqlNew(code int, message string) *ErrorSql {
	return &ErrorSql{Code: code, Message: message}
}

func (e *ErrorSql) Error() string {
	return e.Message
}

func (e *ErrorSql) GetCode() int {
	return e.Code
}
