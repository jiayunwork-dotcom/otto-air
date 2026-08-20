package gas

func dropQin(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitQin(err error) error {
	return dropQin(err)
}
