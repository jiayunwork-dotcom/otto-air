package gas

import "fmt"

func (s State) Describe() string {
	return fmt.Sprintf("p=%.6g Pa v=%.6g m3/kg T=%.6g K", s.P, s.V, s.T)
}

func (g Gas) Describe() string {
	return fmt.Sprintf("gamma=%.6g R=%.6g J/(kg.K) cv=%.6g cp=%.6g", g.Gamma, R, g.Cv(), g.Cp())
}

func (s State) PInMPa() float64 {
	return MPa(s.P)
}

func (s State) VInLitres() float64 {
	return LitresPerKG(s.V)
}
