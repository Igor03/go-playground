package main

import (
	"context"
	"fmt"
	"playground/api"
	file "playground/internal/fs"
	"strconv"
	"strings"
)

// @title        Cache API
// @version      1.0
// @description  Dummy API for uploading artifacts
// @BasePath     /
func main() {

	ctx := context.Background()

	path := "static/inputfile.txt"

	output, _ := file.ReadFile(path)
	hash, _ := file.CreateHashForFile(path)

	fmt.Printf("Hash for file %s: %s\n", path, hash)

	api := api.NewApi(api.ApiSettings{
		BaseUrl:   "http://localhost:8080",
		Port:      8080,
		UploadDir: "./uploads",
	})

	go func() {
		if err := api.Start(ctx); err != nil {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	for _, line := range output {
		filedata := strings.Split(line, " ")

		fSize, _ := strconv.ParseInt(filedata[0], 10, 64)

		fType, fName := filedata[1], "dummyfile"
		file.CreateDummyFile(fSize, fType, fName, "static/")
	}

	// Waiting for the CTRL+C signal to stop the server
	<-ctx.Done()
}
