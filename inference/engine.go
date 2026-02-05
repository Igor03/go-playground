package inference

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ollama/ollama/api"
)

type EngineOptions struct {
	// EngineOptions holds configuration options for the Engine.
	ModelsLocalRepositoryPath string
}

type Engine struct {
	// Engine represents the inference engine.
	options EngineOptions
}

func NewEngine(options EngineOptions) *Engine {
	return &Engine{
		options: options,
	}
}

func (e *Engine) sha256File(modelName string) (string, error) {

	f, err := os.Open(e.options.ModelsLocalRepositoryPath + "/" + modelName)

	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

func (e *Engine) RunInference(ctx context.Context, modelName string, prompt string) (string, error) {

	oc, err := api.ClientFromEnvironment()

	if err != nil {
		return "", err
	}

	req := &api.GenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		Options: map[string]any{
			"num_predict": 128,
		},
	}

	var response string

	err = oc.Generate(ctx, req, func(r api.GenerateResponse) error {
		response += r.Response
		return nil
	})

	return response, nil
}

func (e *Engine) CreateOllamaModel(ctx context.Context, modelName string) error {

	ggufPath := filepath.Join(e.options.ModelsLocalRepositoryPath, modelName)

	digest, err := e.sha256File(modelName)

	if err != nil {
		panic(err)
	}

	f, err := os.Open(ggufPath)

	if err != nil {
		return err
	}

	defer f.Close()

	oc, err := api.ClientFromEnvironment()

	if err != nil {
		return err
	}

	if err := oc.CreateBlob(ctx, digest, f); err != nil {
		panic(err)
	}

	req := &api.CreateRequest{
		Model: modelName,
		Files: map[string]string{
			filepath.Base(ggufPath): digest,
		},
	}

	err = oc.Create(ctx, req, func(p api.ProgressResponse) error {
		fmt.Println(p.Status)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return nil

}
