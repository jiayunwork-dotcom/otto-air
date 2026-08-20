package cycle

import (
	"otto-air/internal/gas"
	"otto-air/internal/process"
)

type SegmentReport struct {
	Segment int       `json:"segment"`
	Kind    string    `json:"kind"`
	From    gas.State `json:"from"`
	To      gas.State `json:"to"`
	Work    float64   `json:"work"`
	Heat    float64   `json:"heat"`
	Entropy float64   `json:"entropy_change"`
}

func SegmentReports(res Result) []SegmentReport {
	g := gas.Gas{Gamma: res.Gamma}
	in := Input{Mode: res.Mode}
	segments := []process.Segment{
		process.MakeSegment(process.Isentropic, res.States[0], res.States[1]),
		heatingSegment(g, in, res.States),
		process.MakeSegment(process.Isentropic, res.States[2], res.States[3]),
		process.MakeSegment(process.Isochoric, res.States[3], res.States[0]),
	}
	reports := make([]SegmentReport, 0, len(segments))
	for i, seg := range segments {
		heat := 0.0
		switch i {
		case 1:
			heat = res.Qin
		case 3:
			heat = -res.QOut
		}
		reports = append(reports, SegmentReport{
			Segment: i + 1,
			Kind:    seg.Kind.String(),
			From:    seg.From,
			To:      seg.To,
			Work:    process.SegmentWork(g, seg, process.DefaultIsentropicSamples),
			Heat:    heat,
			Entropy: gas.EntropyDelta(g, seg.From, seg.To),
		})
	}
	return reports
}

func SegmentWorkSum(res Result) float64 {
	total := 0.0
	for _, r := range SegmentReports(res) {
		total += r.Work
	}
	return total
}

func SegmentEntropyDrift(res Result) float64 {
	total := 0.0
	for _, r := range SegmentReports(res) {
		if r.Kind == "isentropic" {
			total += absFloat(r.Entropy)
		}
	}
	return total
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
