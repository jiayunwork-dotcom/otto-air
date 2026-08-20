package cycle

import "otto-air/internal/process"

type StateReport struct {
	Point int     `json:"point"`
	P     float64 `json:"p"`
	V     float64 `json:"v"`
	T     float64 `json:"t"`
}

type CycleReport struct {
	R        float64       `json:"r"`
	Gamma    float64       `json:"gamma"`
	Mode     string        `json:"mode"`
	Cv       float64       `json:"cv"`
	Rgas     float64       `json:"R"`
	States   [4]StateReport `json:"states"`
	QIn      float64       `json:"q_in"`
	QOut     float64       `json:"q_out"`
	WNet     float64       `json:"w_net"`
	Eta      float64       `json:"eta"`
	EtaWork  float64       `json:"eta_work"`
	EtaClose float64       `json:"eta_closed"`
	MEP      float64       `json:"mep"`
	Area     float64       `json:"area"`
}

func (r Result) StateReports() [4]StateReport {
	return [4]StateReport{
		{Point: 1, P: r.States[0].P, V: r.States[0].V, T: r.States[0].T},
		{Point: 2, P: r.States[1].P, V: r.States[1].V, T: r.States[1].T},
		{Point: 3, P: r.States[2].P, V: r.States[2].V, T: r.States[2].T},
		{Point: 4, P: r.States[3].P, V: r.States[3].V, T: r.States[3].T},
	}
}

func (r Result) CycleReport() CycleReport {
	return CycleReport{
		R:        r.R,
		Gamma:    r.Gamma,
		Mode:     string(r.Mode),
		Cv:       r.Cv,
		Rgas:     r.Rgas,
		States:   r.StateReports(),
		QIn:      r.Qin,
		QOut:     r.QOut,
		WNet:     r.WNet,
		Eta:      r.Eta,
		EtaWork:  r.EtaWork,
		EtaClose: r.EtaClose,
		MEP:      r.MEP,
		Area:     r.Area,
	}
}

type PVPoint struct {
	P float64 `json:"p"`
	V float64 `json:"v"`
}

type PVReport struct {
	Mode  string    `json:"mode"`
	Curve []PVPoint `json:"curve"`
	Area  float64   `json:"area"`
	WNet  float64   `json:"w_net"`
	Eta   float64   `json:"eta"`
}

func BuildPVReport(in Input, pts []process.Point, area, wNet, eta float64) PVReport {
	curve := make([]PVPoint, 0, len(pts))
	for _, p := range pts {
		curve = append(curve, PVPoint{P: p.P, V: p.V})
	}
	return PVReport{
		Mode:  string(in.Mode),
		Curve: curve,
		Area:  area,
		WNet:  wNet,
		Eta:   eta,
	}
}
