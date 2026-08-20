package cycle

import "otto-air/internal/gas"

func computeEndpoints(g gas.Gas, in Input) (gas.States, error) {
	v1 := gas.VolumeFromPT(g, in.Intake.P, in.Intake.T)
	state1 := gas.State{P: in.Intake.P, V: v1, T: in.Intake.T}
	v2 := v1 / in.R
	state2 := applyT2(gas.IsentropicState(g, state1, v2), state1)
	state3 := heatState(g, state2, in.Qin, in.Mode)
	state4 := gas.IsentropicState(g, state3, v1)
	return gas.States{state1, state2, state3, state4}, nil
}

func compressionState(g gas.Gas, intake Intake, r float64) gas.State {
	v1 := gas.VolumeFromPT(g, intake.P, intake.T)
	return gas.IsentropicState(g, gas.State{P: intake.P, V: v1, T: intake.T}, v1/r)
}
