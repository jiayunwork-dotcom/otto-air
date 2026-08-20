package web

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"otto-air/internal/cycle"
)

func newTestServer() http.Handler {
	return NewServer(Assets{
		WebFS: nil,
		Examples: map[string][]byte{
			"r8": []byte(`{"r":8,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000,"mode":"otto"}`),
		},
	})
}

func validBody() string {
	return `{"r":8,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000,"mode":"otto"}`
}

func TestCycleEndpointSuccess(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/cycle", strings.NewReader(validBody()))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var rep cycle.CycleReport
	if err := json.Unmarshal(rec.Body.Bytes(), &rep); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(rep.States) != 4 {
		t.Errorf("states = %d, want 4", len(rep.States))
	}
	want := cycle.ClosedFormEfficiency(8, 1.4)
	if math.Abs(rep.Eta-want) > 1e-9 {
		t.Errorf("eta = %v, want %v", rep.Eta, want)
	}
	if math.Abs(rep.EtaClose-rep.Eta) > 1e-9 {
		t.Errorf("eta_closed = %v, want equal to eta %v", rep.EtaClose, rep.Eta)
	}
}

func TestPvEndpointClosedCurve(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/pv", strings.NewReader(validBody()))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var rep cycle.PVReport
	if err := json.Unmarshal(rec.Body.Bytes(), &rep); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(rep.Curve) < 3 {
		t.Fatalf("curve points = %d, want >= 3", len(rep.Curve))
	}
	first := rep.Curve[0]
	last := rep.Curve[len(rep.Curve)-1]
	if first.P != last.P || first.V != last.V {
		t.Errorf("curve not closed: first(p=%v v=%v) last(p=%v v=%v)", first.P, first.V, last.P, last.V)
	}
	if math.Abs(rep.Area-rep.WNet) > 5e-3*math.Abs(rep.WNet) {
		t.Errorf("area = %v, w_net = %v, want close", rep.Area, rep.WNet)
	}
}

func TestCycleEndpointErrorJSON(t *testing.T) {
	srv := newTestServer()
	bodies := []string{
		`{"r":1,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000}`,
		`{"r":8,"gamma":1.0,"intake":{"t":300,"p":101325},"qin":1200000}`,
		`{"r":8,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":-5}`,
		`{"r":8,"gamma":1.4,"intake":{"t":0,"p":101325},"qin":1200000}`,
		`{"r":8,"gamma":1.4,"intake":{"t":300,"p":0},"qin":1200000}`,
		`{"r":8,"gamma":1.4,"intake":{"t":300,"p":101325}}`,
		`{"r":8,"gamma":1.4,"intake":{"t":300},"qin":1200000}`,
		`{"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000}`,
		`not-json`,
		`{"r":8,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000,"mode":"hybrid"}`,
	}
	for _, body := range bodies {
		req := httptest.NewRequest(http.MethodPost, "/api/cycle", strings.NewReader(body))
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body %q: status = %d, want 400", body, rec.Code)
			continue
		}
		var eb errorBody
		if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
			t.Errorf("body %q: error response not JSON: %v", body, err)
			continue
		}
		if eb.Error == "" {
			t.Errorf("body %q: empty error message", body)
		}
	}
}

func TestPvEndpointErrorJSON(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/pv", strings.NewReader(`{"r":1,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000}`))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Fatalf("error response not JSON: %v", err)
	}
	if eb.Error == "" {
		t.Error("empty error message")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/cycle", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Errorf("405 body not JSON: %v", err)
	}
	if eb.Error == "" {
		t.Error("empty error message on 405")
	}
}

func TestExamplesEndpoint(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/examples", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var rep examplesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &rep); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(rep.R8) == 0 {
		t.Error("example r8 empty")
	}
}

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var rep healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &rep); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if rep.Status != "ok" {
		t.Errorf("status = %q, want ok", rep.Status)
	}
}

func TestDieselEndpointSuccess(t *testing.T) {
	srv := newTestServer()
	body := `{"r":8,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000,"mode":"diesel"}`
	req := httptest.NewRequest(http.MethodPost, "/api/cycle", strings.NewReader(body))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var rep cycle.CycleReport
	if err := json.Unmarshal(rec.Body.Bytes(), &rep); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if rep.Mode != "diesel" {
		t.Errorf("mode = %q, want diesel", rep.Mode)
	}
	if rep.States[1].V == rep.States[2].V {
		t.Error("diesel heating should expand specific volume between points 2 and 3")
	}
}
