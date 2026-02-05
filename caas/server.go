package caas

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type ServerSettings struct {
	BaseUrl   string
	Port      int
	UploadDir string
}

type Server struct {
	Settings ServerSettings
	srv      *http.Server
}

func NewServer(settings ServerSettings) *Server {
	return &Server{
		Settings: settings,
	}
}

func (s *Server) Serve(ctx context.Context) error {

	mux := http.NewServeMux()
	mux.HandleFunc("/artifact", cacheArtifactHandler)

	addr := fmt.Sprintf(":%d", s.Settings.Port)

	s.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v\n", err)
		}
	}()

	log.Printf("listening on: %s\n", addr)

	// Wait for shutdown signal
	<-ctx.Done()

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(shutdownCtx)
}

func cacheArtifactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "expected multipart/form-data", http.StatusBadRequest)
		return
	}

	var (
		key string
		// If we want to check the consistency of the generated Hash based on the uploaded file
		// computedSHA  string
	)

	for {
		part, err := mr.NextPart()

		if errors.Is(err, io.EOF) {
			break
		}

		defer part.Close()

		switch part.FormName() {
		case "key":
			b, _ := io.ReadAll(io.LimitReader(part, 1024))
			key = strings.TrimSpace(string(b))
			if key == "" {
				http.Error(w, "missing key", http.StatusBadRequest)
				return
			}
			key = filepath.Base(key)
		case "file":
			if key == "" {
				http.Error(w, "send key before file", http.StatusBadRequest)
				return
			}

			// This file is supposed to be chunked and sent to some caching manager, such as CaaS Agent
			// file, _, _ := r.FormFile("file")
		default:
			return
		}

		w.WriteHeader(http.StatusOK)

	}

}
