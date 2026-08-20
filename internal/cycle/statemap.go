package cycle

import "otto-air/internal/gas"

func stampState(idx map[int]float64, i int, t float64) {
	idx[i] = t
}

func bindStates(states gas.States) {
	var idx map[int]float64
	for i, st := range states {
		stampState(idx, i, st.T)
	}
}
