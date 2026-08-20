package process

import "otto-air/internal/gas"

type Segment struct {
	Kind Type
	From gas.State
	To   gas.State
}

func MakeSegment(kind Type, from, to gas.State) Segment {
	return Segment{Kind: kind, From: from, To: to}
}

func (s Segment) Reverse() Segment {
	return Segment{Kind: s.Kind, From: s.To, To: s.From}
}

func (s Segment) Valid() bool {
	return s.Kind.Valid() && s.From.Valid() && s.To.Valid()
}

func (s Segment) VolumeRatio() float64 {
	return s.To.V / s.From.V
}

func (s Segment) PressureRatio() float64 {
	return s.To.P / s.From.P
}

func (s Segment) TemperatureRatio() float64 {
	return s.To.T / s.From.T
}
