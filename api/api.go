package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	_ "playground/swagger"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	swagger "github.com/swaggo/http-swagger"
)

type ApiSettings struct {
	BaseUrl   string
	Port      int
	UploadDir string
}

type Api struct {
	Settings ApiSettings
	srv      *http.Server
}

func NewApi(settings ApiSettings) *Api {
	return &Api{
		Settings: settings,
	}
}

func (s *Api) Start(ctx context.Context) error {

	r := chi.NewRouter()

	r.Post("/artifact", cacheArtifactHandler)
	r.Get("/swagger/*", swagger.WrapHandler)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/artifact", cacheArtifactHandler)
	// mux.HandleFunc("/swagger/", swagger.WrapHandler)

	addr := fmt.Sprintf(":%d", s.Settings.Port)

	s.srv = &http.Server{
		Addr:    addr,
		Handler: r,
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

type DefaultOuputput struct {
	ElapsedTime string `json:"elapsed_time"`
}

// cacheArtifactHandler godoc
// @Summary      Upload an artifact
// @Description  Receives multipart form with key + file
// @Accept       multipart/form-data
// @Produce      plain
// @Param        key   formData  string  true  "hash/key"
// @Param        file  formData  file    true  "artifact file"
// @Success      200   {string}  string  "ok"
// @Failure      400   {string}  string  "bad request"
// @Router       /artifact [post]
func cacheArtifactHandler(w http.ResponseWriter, r *http.Request) {

	// With Chi we dont need this
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 	// log.Println("error here")
	// 	return
	// }

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
