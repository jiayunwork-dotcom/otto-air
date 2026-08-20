package cycle

import (
	"otto-air/internal/gas"
)

type Result struct {
	R        float64
	Gamma    float64
	Intake   Intake
	Qin      float64
	Mode     Mode
	Cv       float64
	Rgas     float64
	States   gas.States
	QOut     float64
	WNet     float64
	Eta      float64
	EtaWork  float64
	EtaClose float64
	MEP      float64
	Area     float64
}

func Solve(in Input) (Result, error) {
	in = in.Normalize()
	if err := in.Validate(); err != nil {
		return Result{}, err
	}
	g, err := gas.New(in.Gamma)
	if err != nil {
		return Result{}, err
	}
	states, err := computeEndpoints(g, in)
	if err != nil {
		return Result{}, err
	}
	if err := gas.ValidateResultStates(states); err != nil {
		return Result{}, err
	}
	bindStates(states)
	qOut, wNet := netWork(g, states, in.Qin)
	etaHeat, etaWork := efficiency(in.Qin, qOut, wNet)
	mep := meanEffectivePressure(wNet, states[0].V, states[1].V)
	area := closedAreaFor(g, in, states)
	return Result{
		R:        in.R,
		Gamma:    in.Gamma,
		Intake:   in.Intake,
		Qin:      in.Qin,
		Mode:     in.Mode,
		Cv:       g.Cv(),
		Rgas:     gas.R,
		States:   states,
		QOut:     qOut,
		WNet:     wNet,
		Eta:      etaHeat,
		EtaWork:  etaWork,
		EtaClose: ClosedFormEfficiency(in.R, in.Gamma),
		MEP:      mep,
		Area:     area,
	}, nil
}
