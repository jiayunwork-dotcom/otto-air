package process

type Type int

const (
	Isentropic Type = iota
	Isochoric
	Isobaric
)

func (t Type) String() string {
	switch t {
	case Isentropic:
		return "isentropic"
	case Isochoric:
		return "isochoric"
	case Isobaric:
		return "isobaric"
	}
	return "unknown"
}

func (t Type) Valid() bool {
	switch t {
	case Isentropic, Isochoric, Isobaric:
		return true
	}
	return false
}
