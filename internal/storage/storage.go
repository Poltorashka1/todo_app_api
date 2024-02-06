package storage

import (
	"log/slog"
	"time"
	"web/internal/config"
)

// Storage - database interface
type Storage interface {
	// Connect connect to database.
	// cfg *config.Config - configuration of database connection.
	// log *slog.Logger - logger.
	Connect(cfg *config.Config, log *slog.Logger) Storage

	// CreateTask creates new task with selected parameters.
	CreateTask(text string, tags []string, dueDate *time.Time) error

	// GetTask gets task by ID.
	GetTask(id int) (*Task, error)

	// GetAllTasksByTag returns all tasks by tag or tags list.
	// tag []string - list of tags or single tag
	// returns *storage.AllTasks - list of storage.Task
	GetTasksByTagShort(tagList []string) (*Tasks, error)

	// GetTaskByTag returns tasks only with specified tag or tags list.
	// tag []string - list of tags or single tag.
	// returns *storage.AllTasks - list of storage.Task.
	GetTasksByTagFull(tagList []string) (*Tasks, error)

	// Delete deletes task by ID or all tasks.
	// TODO check admin rules
	DeleteTask(id ...string) error

	// GetTag returns tag by name.
	GetTag(name string) (*Tag, error)

	// GetAllTags returns all tags.
	GetAllTags() (*Tags, error)

	// GetAllTasks returns all tasks.
	GetAllTasks() (*Tasks, error)

	// GetTasksByDueDate returns tasks by due date.
	GetTasksByDueDate(due *time.Time) (*Tasks, error)

	// CreateTag creates new tag
	CreateTag(name string) error

	// DeleteTag deletes tag
	DeleteTag(name ...string) error

	// GetTasksByDueAndTag returns tasks by due date and tag
	GetTasksByDueAndTagFull(tagList []string, due *time.Time) (*Tasks, error)
	GetTasksByDueAndTagShort(tagList []string, due *time.Time) (*Tasks, error)

	GetTasksByTag(tagList []string) (*Tasks, error)

	GetTasksByTagAndDue(tagList []string, due *time.Time) (*Tasks, error)
}

type SqlError interface {
	Error() string
	GetCode() int
}
