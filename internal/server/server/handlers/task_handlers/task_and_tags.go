package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
	"time"
	"web/internal/server/context/response"
	"web/internal/storage"
	"web/internal/storage/models"
)

// GetTasksByTagHandler returns tasks that have one of the specified tags from the query
// @Summary Get tasks by tag and due date
// @Description Tag: returns tasks that have one of the specified tags from the query. Tag: a tag or multiple tags separated by a comma(',') without spaces. Due: due date format: 2006-01-02T15:04:05Z
// @Tags tasks_tags
// @Accept json
// @Produce json
// @Param tag query string true "Tags"
// @Param due query string false "Due"
// @Success 200 {object} response.OkResponse{data=storage.Tasks}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/tag/ [get]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.GetTasksByTagOrByTagAndDueHandler
func (h *TaskHandler) GetTasksByTagHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tagList := strings.Split(query.Get("tag"), ",")
	due := query.Get("due")

	var tasks *models.Tasks
	var err error

	switch due {
	case "":
		err = validateTags(tagList, h.AllTags)
		if err != nil {
			h.JSON(w, response.Error(http.StatusBadRequest, err))
			return
		}
		tasks, err = h.Db.GetTasksByTag(tagList)
	default:
		err = validateTagsAndDue(tagList, due, h.AllTags)
		if err != nil {
			h.JSON(w, response.Error(http.StatusBadRequest, err))
			return
		}

		dueDate, _ := time.Parse(time.RFC3339, due)
		tasks, err = h.Db.GetTasksByTagAndDue(tagList, &dueDate)
	}
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	h.JSON(w, response.OK(tasks))
}

// GetTasksByModeAndTagHandler returns tasks that have one of the specified tags from the query
// @Summary Get tasks by mode and tag
// @Description Mode: "full" returns tasks with the specified tag, or all of the specified tags in the query. "short" returns tasks with only the specified tag, or only all specified tags in the query. Tag: a tag or multiple tags separated by a comma(',') without spaces. Due: due date format: 2006-01-02T15:04:05Z
// @Tags tasks_tags
// @Accept json
// @Produce json
// @Param mode path string true "Mode"
// @Param tag query string true "Tags"
// @Param due query string false "Due"
// @Success 200 {object} response.OkResponse{data=storage.Tasks}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /task/tag/{mode}/ [get]
// Context from Function internal/server/server/handlers/task.go:handlers.*Handlers.GetTasksByModeAndTagHandler
func (h *TaskHandler) GetTasksByModeAndTagHandler(w http.ResponseWriter, r *http.Request) {
	mode := chi.URLParam(r, "mode")

	var tasks *models.Tasks
	var err error

	switch mode {
	case "full":
		tasks, err = h.getTasksFull(r)
	case "short":
		tasks, err = h.getTasksShort(r)
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

func (h *TaskHandler) getTasksFull(r *http.Request) (*models.Tasks, error) {
	query := r.URL.Query()
	tagList := strings.Split(query.Get("tag"), ",")
	due := query.Get("due")

	switch {
	case due == "":
		err := validateTags(tagList, h.AllTags)
		if err != nil {
			return nil, err
		}
		return h.Db.GetTasksByTagFull(tagList)
	default:
		err := validateTagsAndDue(tagList, due, h.AllTags)
		deuDate, _ := time.Parse(time.RFC3339, due)
		if err != nil {
			return nil, err
		}
		return h.Db.GetTasksByDueAndTagFull(tagList, &deuDate)
	}

}

func (h *TaskHandler) getTasksShort(r *http.Request) (*models.Tasks, error) {
	query := r.URL.Query()
	tagList := strings.Split(query.Get("tag"), ",")
	due := query.Get("due")

	switch {
	case due == "":
		err := validateTags(tagList, h.AllTags)
		if err != nil {
			return nil, err
		}
		return h.Db.GetTasksByTagShort(tagList)
	default:
		err := validateTagsAndDue(tagList, due, h.AllTags)
		deuDate, _ := time.Parse(time.RFC3339, due)
		if err != nil {
			return nil, err
		}
		return h.Db.GetTasksByDueAndTagShort(tagList, &deuDate)
	}
}
