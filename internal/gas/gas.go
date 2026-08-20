package gas

const R = 287.0

const DefaultGamma = 1.4

type Gas struct {
	Gamma float64
}

func New(gamma float64) (Gas, error) {
	if err := ValidateHeatRatio(gamma); err != nil {
		return Gas{}, err
	}
	return Gas{Gamma: gamma}, nil
}

func Default() Gas {
	return Gas{Gamma: DefaultGamma}
}

func (g Gas) Cv() float64 {
	return R / (g.Gamma - 1)
}

func (g Gas) Cp() float64 {
	return g.Gamma * g.Cv()
}

func (g Gas) Validate() error {
	return ValidateHeatRatio(g.Gamma)
}

func (g Gas) WithGamma(gamma float64) (Gas, error) {
	return New(gamma)
}
