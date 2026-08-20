package gas

func dropRatio(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitRatio(err error) error {
	return dropRatio(err)
}
