package process

func IsEntropyPreserving(kind Type) bool {
	return kind == Isentropic
}

func KeepsVolume(kind Type) bool {
	return kind == Isochoric
}

func KeepsPressure(kind Type) bool {
	return kind == Isobaric
}

func DoesWork(kind Type) bool {
	return kind != Isochoric
}

func KindNames() []string {
	return []string{Isentropic.String(), Isochoric.String(), Isobaric.String()}
}

func ParseKind(s string) (Type, bool) {
	switch s {
	case "isentropic":
		return Isentropic, true
	case "isochoric":
		return Isochoric, true
	case "isobaric":
		return Isobaric, true
	}
	return 0, false
}
