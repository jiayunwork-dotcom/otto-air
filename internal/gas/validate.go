package gas

import (
	"errors"
	"fmt"
	"math"
)

var ErrInvalidState = errors.New("state contains non-finite or non-positive values")

func Finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func CheckFinite(name string, v float64) error {
	if !Finite(v) {
		return fmt.Errorf("field %s must be finite, got %v", name, v)
	}
	return nil
}

func ValidateCompressionRatio(r float64) error {
	if err := CheckFinite("r", r); err != nil {
		return err
	}
	if r <= 1 {
		return fmt.Errorf("compression ratio r must be greater than 1, got %v", r)
	}
	return nil
}

func ValidateHeatRatio(gamma float64) error {
	if err := CheckFinite("gamma", gamma); err != nil {
		return err
	}
	if gamma <= 1 {
		return fmt.Errorf("specific heat ratio gamma must be greater than 1, got %v", gamma)
	}
	return nil
}

func ValidateHeatInput(q float64) error {
	if err := CheckFinite("qin", q); err != nil {
		return err
	}
	if q <= 0 {
		return fmt.Errorf("heat input q_in must be positive, got %v", q)
	}
	return nil
}

func ValidateTemperature(t float64) error {
	if err := CheckFinite("intake.t", t); err != nil {
		return err
	}
	if t <= 0 {
		return fmt.Errorf("intake temperature must be positive, got %v", t)
	}
	return nil
}

func ValidatePressure(p float64) error {
	if err := CheckFinite("intake.p", p); err != nil {
		return err
	}
	if p <= 0 {
		return fmt.Errorf("intake pressure must be positive, got %v", p)
	}
	return nil
}

func ValidateIntakeState(t, p float64) error {
	if err := ValidateTemperature(t); err != nil {
		return err
	}
	return ValidatePressure(p)
}

func ValidateResultStates(states States) error {
	if !states.Valid() {
		return ErrInvalidState
	}
	return nil
}
