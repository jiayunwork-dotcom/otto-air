package process

const (
	DefaultIsentropicSamples = 48
	DefaultIsobaricSamples   = 16
)

func SampleCountFor(kind Type, n int) int {
	switch kind {
	case Isentropic:
		if n <= 0 {
			return DefaultIsentropicSamples
		}
		return n
	case Isobaric:
		if n <= 0 {
			return DefaultIsobaricSamples
		}
		return n
	case Isochoric:
		return 2
	}
	return 2
}
