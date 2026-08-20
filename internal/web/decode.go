package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"otto-air/internal/cycle"
)

type intakeJSON struct {
	T *float64 `json:"t"`
	P *float64 `json:"p"`
}

type cycleRequest struct {
	R      *float64    `json:"r"`
	Gamma  *float64    `json:"gamma"`
	Intake *intakeJSON `json:"intake"`
	Qin    *float64    `json:"qin"`
	Mode   string      `json:"mode"`
}

func decodeCycleRequest(r *http.Request) (cycle.Input, error) {
	var req cycleRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		return cycle.Input{}, fmt.Errorf("invalid JSON body: %v", err)
	}
	if req.R == nil {
		return cycle.Input{}, errors.New("field r is required")
	}
	if req.Gamma == nil {
		return cycle.Input{}, errors.New("field gamma is required")
	}
	if req.Qin == nil {
		return cycle.Input{}, errors.New("field qin is required")
	}
	if req.Intake == nil {
		return cycle.Input{}, errors.New("field intake is required")
	}
	if req.Intake.T == nil {
		return cycle.Input{}, errors.New("field intake.t is required")
	}
	if req.Intake.P == nil {
		return cycle.Input{}, errors.New("field intake.p is required")
	}
	mode, err := cycle.ParseMode(req.Mode)
	if err != nil {
		return cycle.Input{}, err
	}
	return cycle.Input{
		R:      *req.R,
		Gamma:  *req.Gamma,
		Intake: cycle.Intake{T: *req.Intake.T, P: *req.Intake.P},
		Qin:    *req.Qin,
		Mode:   mode,
	}, nil
}
