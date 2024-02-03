package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"strings"
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

// validate id don't use because chi do it automatically +-
func validateId(idString string) (*int, error) {
	if idString[0] == ' ' {
		return nil, fmt.Errorf("id param must not be empty")
	}
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return nil, fmt.Errorf("id param must be integer") // Todo check using method by due date
	}
	if idInt < 0 {
		return nil, fmt.Errorf("id param must be positive") // Todo check using method by due date
	}
	return &idInt, nil
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
	err = h.Db.CreateTask(requestData.Text, requestData.Tags, requestData.Due)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	h.JSON(w, response.OK())
}

// Todo merge to one handler 2 delete handlers

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

// GetTasksByTagHandler get tasks by tag(in Url)
// @Summary Get tasks by tag
// @Description Mode: 'full' - returns tasks who have includes specified tag or tags list in the task tags, 'short' - returns tasks who have only specified tag or tags list included in the task tags. Tag: tag or tags using , as separator
// @Tags tasks
// @Accept json
// @Produce json
// @Param mode path string true "Mode"
// @Param tag query string true "Tags"
// @Success 200 {object} response.OkResponse{data=storage.Task}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/{mode}/ [get]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.GetTasksByTagHandler
func (h *Handlers) GetTasksByTagHandler(w http.ResponseWriter, r *http.Request) {
	mode := chi.URLParam(r, "mode")
	// first ex for tag in url
	//tags := chi.URLParam(r, "tag")
	//fmt.Println(tags)
	//tagList := strings.Split(tags, "%2C")

	// second ex for tag in query
	query := r.URL.Query()
	tagList := strings.Split(query.Get("tag"), ",")

	err := validateTags(tagList, h.AllTags)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	var tasks *storage.Tasks

	switch mode {
	case "full":
		tasks, err = h.Db.GetAllTasksByTag(tagList)
	case "short":
		tasks, err = h.Db.GetTaskByTag(tagList)
	default:
		h.JSON(w, response.Error(http.StatusBadRequest, fmt.Errorf("expect mode == full or short, got %v", mode)))
		return
	}

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

	tasks, err := h.Db.GetTasksByDueDate(&due)
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
