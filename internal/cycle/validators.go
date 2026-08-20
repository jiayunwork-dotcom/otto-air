package cycle

import (
	"errors"
	"math"
)

var (
	ErrEtaOutOfRange      = errors.New("thermal efficiency out of (0,1) range")
	ErrNonPositiveWork    = errors.New("net work must be positive")
	ErrNonPositiveMEP     = errors.New("mep must be positive")
	ErrAreaWorkMismatch   = errors.New("p-v closed area does not match net work")
	ErrNonPositiveQOut    = errors.New("rejected heat must be positive")
)

type Validator struct {
	EtaTol  float64
	WorkTol float64
}

func DefaultValidator() Validator {
	return Validator{EtaTol: 1e-9, WorkTol: 5e-3}
}

func (v Validator) Validate(res Result) error {
	if res.Eta <= 0 || res.Eta >= 1 {
		return ErrEtaOutOfRange
	}
	if res.WNet <= 0 {
		return ErrNonPositiveWork
	}
	if res.MEP <= 0 {
		return ErrNonPositiveMEP
	}
	if res.QOut <= 0 {
		return ErrNonPositiveQOut
	}
	scale := math.Max(math.Abs(res.Area), math.Abs(res.WNet))
	if scale > 0 && math.Abs(res.Area-res.WNet) > v.WorkTol*scale {
		return ErrAreaWorkMismatch
	}
	return nil
}

func CheckConsistency(res Result) error {
	return DefaultValidator().Validate(res)
}
