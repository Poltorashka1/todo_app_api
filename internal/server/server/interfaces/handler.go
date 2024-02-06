package handlerInterfaces

import (
	"net/http"
)

// Todo comment all this method

type HandlerMethods interface {
	// MethodNotAllowedHandler is the handler for handling method not allowed requests.
	MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request)
	// NotFoundHandler is the handler for handling not found requests.
	NotFoundHandler(w http.ResponseWriter, r *http.Request)
	// GetTasksHandler get all tasks
	GetTasksHandler(w http.ResponseWriter, r *http.Request)
	// GetTaskHandler get task by id
	GetTaskHandler(w http.ResponseWriter, r *http.Request)
	// GetTagsHandler get all tags
	GetTagsHandler(w http.ResponseWriter, r *http.Request)
	// GetTagHandler get tag by name
	GetTagHandler(w http.ResponseWriter, r *http.Request)
	GetTasksByModeAndTagHandler(w http.ResponseWriter, r *http.Request)
	GetTasksByTagHandler(w http.ResponseWriter, r *http.Request)
	// GetTasksByDueDateHandler get tasks by due date
	GetTasksByDueDateHandler(w http.ResponseWriter, r *http.Request)
	// CreateTaskHandler create new task with specified params
	CreateTaskHandler(w http.ResponseWriter, req *http.Request)
	// DeleteTasksHandler delete all tasks
	DeleteTasksHandler(w http.ResponseWriter, req *http.Request)
	// DeleteTaskHandler delete task by id
	DeleteTaskHandler(w http.ResponseWriter, req *http.Request)
	// CreateTagHandler create new tag
	CreateTagHandler(w http.ResponseWriter, r *http.Request)
	// DeleteTagsHandler delete all tags
	DeleteTagsHandler(w http.ResponseWriter, r *http.Request)
	// DeleteTagHandler delete tag by name
	DeleteTagHandler(w http.ResponseWriter, r *http.Request)
}
