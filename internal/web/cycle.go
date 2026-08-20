package web

import (
	"net/http"

	"otto-air/internal/cycle"
)

func (s *Server) handleCycle(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	in, err := decodeCycleRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	res, err := cycle.Solve(in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, res.CycleReport())
}
