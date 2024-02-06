package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web/internal/server/context/request"
	"web/internal/server/context/response"
	"web/internal/storage"
)

// Todo get by due date
// Todo mutex lock for tasks

// GetTaskHandler returns task by id
// @Summary Get task by id
// @Description Get task by id
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task id"
// @Success 200 {object} response.OkResponse{data=storage.Task}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/{id} [get]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.GetTaskHandler
func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	// Todo check using method by due date
	task, err := h.Db.GetTask(idInt)
	if err != nil {
		switch errSql := err.(type) {
		case storage.SqlError:
			h.JSON(w, response.Error(errSql.GetCode(), errSql))
		default:
			h.JSON(w, response.Error(http.StatusBadRequest, err))
		}
		return
	}

	h.JSON(w, response.OK(task))
}

// GetTasksHandler returns task
// @Summary Get tasks
// @Description Get tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} response.OkResponse{data=storage.Tasks} "Successful response"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/ [get]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.GetTasksHandler
func (h *Handlers) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	allTasks, err := h.Db.GetAllTasks()

	if err != nil {
		switch errSql := err.(type) {
		case storage.SqlError:
			h.JSON(w, response.Error(errSql.GetCode(), errSql))
		default:
			h.JSON(w, response.Error(http.StatusBadRequest, err))
		}
		return
	}

	h.JSON(w, response.OK(allTasks))
}

// CreateTaskHandler creates new task
// @Summary Create new task
// @Description "Create new task object with the following fields: text (string, required) - text of the task, tags ([]string, required) - tags associated with the task, due (string, required) - due date of the task in '2006-01-02T15:04:05Z' format"
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body request.TaskRequest true "Task"
// @Success 200 {object} response.OkResponseEmpty
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/ [post]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.CreateTaskHandler
func (h *Handlers) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var requestData request.TaskRequest

	err := h.DecodeJSON(r.Body, &requestData)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	// Todo refactor all validate request
	err = requestData.ValidateRequest(h.AllTags)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}
	// reqData
	// Todo why reqData param doing cycle import
	// Todo remove this and mak it more beautiful
	dueDate, _ := time.Parse(time.RFC3339, requestData.Due)
	err = h.Db.CreateTask(requestData.Text, requestData.Tags, &dueDate)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	h.JSON(w, response.OK())
}

func (h *Handlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	///
}

// DeleteTasksHandler deletes all tasks
// @Summary Delete tasks
// @Description Delete tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} response.OkResponseEmpty
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/ [delete]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.DeleteTasksHandler
func (h *Handlers) DeleteTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Todo check admin rules

	err := h.Db.DeleteTask()
	if err != nil {
		switch errSql := err.(type) {
		case storage.SqlError:
			h.JSON(w, response.Error(errSql.GetCode(), errSql))
		default:
			h.JSON(w, response.Error(http.StatusBadRequest, err))
		}
		return
	}

	h.JSON(w, response.OK())
}

// DeleteTaskHandler deletes task by id
// @Summary Delete task
// @Description Delete task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} response.OkResponseEmpty
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/{id} [delete]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.DeleteTaskHandler
func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Db.DeleteTask(id)
	if err != nil {
		switch errSql := err.(type) {
		case storage.SqlError:
			h.JSON(w, response.Error(errSql.GetCode(), errSql))
		default:
			h.JSON(w, response.Error(http.StatusBadRequest, err))
		}
		return
	}

	h.JSON(w, response.OK())
}

// GetTasksByDueDateHandler get tasks by due date
// @Summary Get tasks by due date
// @Description Get tasks by due date format: 2006-01-02T15:04:05Z
// @Tags tasks
// @Accept json
// @Produce json
// @Param due path string true "Due date"
// @Success 200 {object} response.OkResponse{data=storage.Tasks}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/{due} [get]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.GetTasksByDueDateHandler
func (h *Handlers) GetTasksByDueDateHandler(w http.ResponseWriter, r *http.Request) {
	due := chi.URLParam(r, "due")
	// INFO: use it because swagger (is shit) request using '%3A' instead of ':'
	due = strings.Replace(due, "%3A", ":", -1)

	err := validateDue(due)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}
	// Todo remove this and mak it more beautiful
	dueDate, _ := time.Parse(time.RFC3339, due)
	tasks, err := h.Db.GetTasksByDueDate(&dueDate)
	if err != nil {
		switch errSql := err.(type) {
		case storage.SqlError:
			h.JSON(w, response.Error(errSql.GetCode(), errSql))
		default:
			h.JSON(w, response.Error(http.StatusBadRequest, err))
		}
		return
	}

	h.JSON(w, response.OK(tasks))
}
