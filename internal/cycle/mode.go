package cycle

import "fmt"

type Mode string

const (
	Otto   Mode = "otto"
	Diesel Mode = "diesel"
)

func ParseMode(s string) (Mode, error) {
	switch s {
	case "":
		return Otto, nil
	case string(Otto):
		return Otto, nil
	case string(Diesel):
		return Diesel, nil
	}
	return "", fmt.Errorf("unknown cycle mode %q, expected %q or %q", s, Otto, Diesel)
}

func (m Mode) String() string {
	return string(m)
}

func (m Mode) Valid() bool {
	switch m {
	case Otto, Diesel:
		return true
	}
	return false
}

func (m Mode) HeatingKind() string {
	if m == Diesel {
		return "isobaric"
	}
	return "isochoric"
}
