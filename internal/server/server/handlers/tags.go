package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
	"web/internal/server/context/request"
	"web/internal/server/context/response"
	"web/internal/storage"
)

// GetTagsHandler returns all tags
// @Summary Get all tags
// @Description Get all tags
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} response.OkResponse{data=storage.Tags}
// @Failure 404 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tag [get]
func (h *Handlers) GetTagsHandler(w http.ResponseWriter, req *http.Request) {
	Tags, err := h.Db.GetAllTags()
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	h.JSON(w, response.OK(Tags))
}

// GetTagHandler returns tag by name
// @Summary Get tag by name
// @Description Get tag by name
// @Tags tags
// @Accept json
// @Produce json
// @Param name path string true "Tag name"
// @Success 200 {object} response.OkResponse{data=storage.Tag}
// @Failure 404 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tag/{name} [get]
func (h *Handlers) GetTagHandler(w http.ResponseWriter, req *http.Request) {
	tagName := chi.URLParam(req, "name")

	tag, err := h.Db.GetTag(tagName)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	h.JSON(w, response.OK(tag))
}

// CreateTagHandler creates new tag
// @Summary Create new tag
// @Description Create new tag with uniq name
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body request.TagRequest true "Tag name"
// @Success 200 {object} response.OkResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tag [post]
func (h *Handlers) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var requestData request.TagRequest

	err := h.DecodeJSON(r.Body, &requestData)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	err = validateTagName(requestData.Name)
	if err != nil {
		h.JSON(w, response.Error(http.StatusBadRequest, err))
		return
	}

	err = h.Db.CreateTag(requestData.Name)
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

// DeleteTagHandler deletes tag by name
// @Summary Delete tag by name
// @Description Delete tag
// @Tags tags
// @Accept json
// @Produce json
// @Param name path string true "Tag name"
// @Success 200 {object} response.OkResponseEmpty
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tag/{name} [delete]
func (h *Handlers) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tagName := chi.URLParam(r, "name")

	err := h.Db.DeleteTag(tagName)
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

// DeleteTagsHandler deletes all tags
// @Summary Delete tags
// @Description Delete tags
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} response.OkResponseEmpty
// @Failure 404 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tag/ [delete]
func (h *Handlers) DeleteTagsHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Db.DeleteTag()
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
