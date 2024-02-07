package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"web/internal/server/context/request"
	"web/internal/server/context/response"
	"web/internal/storage"
	tagsList "web/storage/tags-list"
)

type UserHandler struct {
	Db      storage.Storage
	Log     *slog.Logger
	AllTags *tagsList.TagsList
}

func NewUserHandler(Db *storage.Storage, Log *slog.Logger, AllTags *tagsList.TagsList) *UserHandler {
	return &UserHandler{
		Db:      *Db,
		Log:     Log,
		AllTags: AllTags,
	}
}

func (h *UserHandler) JSON(w http.ResponseWriter, resp response.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.GetStatus())

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.Log.Error(fmt.Sprintf("error encoding response: %v", err))
	}
}

func (h *UserHandler) DecodeJSON(body io.ReadCloser, req request.Request) error {
	const op = "handlers.DecodeJSON"

	// Todo time problem parsing time if time error decode error and all error in response

	err := json.NewDecoder(body).Decode(req)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("request body is empty")
	}
	if err != nil {
		h.Log.Error(fmt.Sprintf("%v: error reading request body: %v", op, err))
	}
	return nil
}

func (h *UserHandler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	h.JSON(w, response.Error(http.StatusMethodNotAllowed, fmt.Errorf("method '%s' not allowed", r.Method)))
}

func (h *UserHandler) NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	h.JSON(w, response.Error(http.StatusNotFound, fmt.Errorf("path '%s' not found", req.URL.Path)))
}