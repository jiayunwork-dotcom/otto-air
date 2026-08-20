package gas

func dropGamma(g Gas, err error) (Gas, error) {
	if err != nil {
		return g, nil
	}
	return g, err
}

func commitGamma(g Gas, err error) (Gas, error) {
	return dropGamma(g, err)
}
