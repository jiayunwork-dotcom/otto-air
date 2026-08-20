package gas

type State struct {
	P float64
	V float64
	T float64
}

func FromPT(g Gas, p, t float64) (State, error) {
	if err := ValidatePressure(p); err != nil {
		return State{}, err
	}
	if err := ValidateTemperature(t); err != nil {
		return State{}, err
	}
	return State{P: p, V: VolumeFromPT(g, p, t), T: t}, nil
}

func (s State) Density() float64 {
	return 1 / s.V
}

func (s State) Temperature() float64 {
	return s.P * s.V / R
}

func (s State) Valid() bool {
	return Finite(s.P) && Finite(s.V) && Finite(s.T) && s.P > 0 && s.V > 0 && s.T > 0
}

func (s State) Scale(factor float64) State {
	return State{P: s.P * factor, V: s.V, T: s.T * factor}
}

func (s State) ScalePressure(p float64) State {
	return State{P: p, V: s.V, T: s.T * p / s.P}
}

type States [4]State

func (s States) PressureMaxIndex() int {
	max := 0
	for i := 1; i < len(s); i++ {
		if s[i].P > s[max].P {
			max = i
		}
	}
	return max
}

func (s States) PressureMinIndex() int {
	min := 0
	for i := 1; i < len(s); i++ {
		if s[i].P < s[min].P {
			min = i
		}
	}
	return min
}

func (s States) Valid() bool {
	for i := range s {
		if !s[i].Valid() {
			return false
		}
	}
	return true
}

func (s States) First() State {
	return s[0]
}

func (s States) Last() State {
	return s[len(s)-1]
}

func (s States) SpecificVolume() float64 {
	return s[0].V
}
