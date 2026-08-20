package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"otto-air/internal/web"
)

//go:embed web
var webFS embed.FS

//go:embed example/r8.json
var r8JSON []byte

func main() {
	httpAddr := flag.String("http", ":8080", "serve the web console on this address (e.g. :8080)")
	flag.Parse()
	handler := web.NewServer(web.Assets{
		WebFS: webFS,
		Examples: map[string][]byte{
			"r8": r8JSON,
		},
	})
	log.Printf("otto-air web console on http://localhost%s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatal(err)
	}
}
