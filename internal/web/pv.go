package web

import (
	"net/http"

	"otto-air/internal/cycle"
)

func (s *Server) handlePV(w http.ResponseWriter, r *http.Request) {
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
	pts, area, err := cycle.Curve(in, 0)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	report := cycle.BuildPVReport(in, pts, area, res.WNet, res.Eta)
	writeJSON(w, http.StatusOK, report)
}
