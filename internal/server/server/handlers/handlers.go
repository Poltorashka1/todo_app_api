package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"web/internal/server/context/request"
	"web/internal/server/context/response"
	"web/internal/server/server"
)

type Handlers struct {
	*server.Server
}

func NewHandlers(s *server.Server) *Handlers {
	return &Handlers{Server: s}
}

func (h *Handlers) JSON(w http.ResponseWriter, resp response.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.GetStatus())

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.Log.Error(fmt.Sprintf("error encoding response: %v", err))
	}
}

func (h *Handlers) DecodeJSON(body io.ReadCloser, req request.Request) error {
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

func (h *Handlers) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	h.JSON(w, response.Error(http.StatusMethodNotAllowed, fmt.Errorf("method '%s' not allowed", r.Method)))
}

func (h *Handlers) NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	h.JSON(w, response.Error(http.StatusNotFound, fmt.Errorf("path '%s' not found", req.URL.Path)))
}
