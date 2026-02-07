package server

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/actiometa/gopherwatch/internal/processor"
	"github.com/actiometa/gopherwatch/web"
)

type ApiServer struct {
	engine *processor.Engine
}

func NewApiServer(e *processor.Engine) *ApiServer {
	return &ApiServer{engine: e}
}

func (s *ApiServer) Start(port string) error {
	content, err := fs.Sub(web.Assets, "dist")
	if err != nil {
		log.Printf("Error serving static files: %v", err)
	}
	http.Handle("/", http.FileServer(http.FS(content)))

	http.HandleFunc("/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		s.handleStats(w, r)
	})

	log.Printf("Dashboard en http://localhost:%s", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *ApiServer) handleStats(w http.ResponseWriter, r *http.Request) {
	var stats []processor.ContainerStats

	// Convertimos el sync.Map a un Slice
	s.engine.Containers.Range(func(key, value interface{}) bool {
		stats = append(stats, value.(processor.ContainerStats))
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
