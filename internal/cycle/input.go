package cycle

import (
	"fmt"

	"otto-air/internal/gas"
)

type Intake struct {
	T float64
	P float64
}

type Input struct {
	R      float64
	Gamma  float64
	Intake Intake
	Qin    float64
	Mode   Mode
}

func (in Input) Normalize() Input {
	if in.Mode == "" {
		in.Mode = Otto
	}
	return in
}

func (in Input) Validate() error {
	if err := gas.ValidateCompressionRatio(in.R); err != nil {
		return err
	}
	if err := gas.ValidateHeatRatio(in.Gamma); err != nil {
		return err
	}
	if err := gas.ValidateHeatInput(in.Qin); err != nil {
		return err
	}
	if err := gas.ValidateIntakeState(in.Intake.T, in.Intake.P); err != nil {
		return err
	}
	switch in.Mode {
	case Otto, Diesel:
	default:
		return fmt.Errorf("unknown cycle mode %q, expected %q or %q", in.Mode, Otto, Diesel)
	}
	return nil
}

func (in Input) WithR(r float64) Input {
	in.R = r
	return in
}

func (in Input) WithGamma(gamma float64) Input {
	in.Gamma = gamma
	return in
}

func (in Input) WithQin(qin float64) Input {
	in.Qin = qin
	return in
}

func (in Input) WithMode(mode Mode) Input {
	in.Mode = mode
	return in
}
