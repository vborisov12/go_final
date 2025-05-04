package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func Start(port int) error {
	webDir, err := filepath.Abs("./web")
	if err != nil {
		return fmt.Errorf("Filed to get absolute path: %w", err)
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fileServer)

	log.Printf("Starting server on port %d", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		return fmt.Errorf("Filed to start server: %w", err)
	}
	return nil
}
