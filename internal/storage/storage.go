package storage

import (
	"log/slog"
	"web/internal/config"
)

// Storage - database interface
type Storage interface {
	// Connect connect to database.
	// cfg *config.Config - configuration of database connection.
	// log *slog.Logger - logger.
	Connect(cfg *config.Config, log *slog.Logger) Storage

	// CreateTask creates new task with selected parameters.
	CreateTask(text string, tags []string, dueDate string) error

	// GetTask gets task by ID.
	GetTask(id int) (*Task, error)

	// GetAllTasksByTag returns all tasks by tag or tags list.
	// tag []string - list of tags or single tag
	// returns *storage.AllTasks - list of storage.Task
	GetAllTasksByTag(tagList []string) (*Tasks, error)

	// GetTaskByTag returns tasks only with specified tag or tags list.
	// tag []string - list of tags or single tag.
	// returns *storage.AllTasks - list of storage.Task.
	GetTaskByTag(tagList []string) (*Tasks, error)

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
	GetTasksByDueDate(due *string) (*Tasks, error)

	// CreateTag creates new tag
	CreateTag(name string) error

	// DeleteTag deletes tag
	DeleteTag(name ...string) error
}

type SqlError interface {
	Error() string
	GetCode() int
}
