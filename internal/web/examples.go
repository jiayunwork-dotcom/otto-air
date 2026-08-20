package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

var errExampleNotFound = errors.New("example not found")

type examplesResponse struct {
	R8 json.RawMessage `json:"r8"`
}

func (s *Server) handleExamples(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	r8, ok := s.assets.Examples["r8"]
	if !ok {
		writeError(w, http.StatusNotFound, errExampleNotFound)
		return
	}
	writeJSON(w, http.StatusOK, examplesResponse{R8: r8})
}
