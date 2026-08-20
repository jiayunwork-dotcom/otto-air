package web

import (
	"io/fs"
	"net/http"
)

type Assets struct {
	WebFS    fs.FS
	Examples map[string][]byte
}

type Server struct {
	assets Assets
	mux    *http.ServeMux
}

func NewServer(assets Assets) http.Handler {
	s := &Server{assets: assets, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/cycle", s.handleCycle)
	s.mux.HandleFunc("/api/pv", s.handlePV)
	s.mux.HandleFunc("/api/examples", s.handleExamples)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.Handle("/", staticHandler(assets.WebFS))
	return limitBody(withStatusLog(s.mux))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
