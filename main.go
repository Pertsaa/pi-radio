package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Pertsaa/pi-radio/handler"
	"github.com/Pertsaa/pi-radio/middleware"
	"github.com/Pertsaa/pi-radio/radio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./pi-radio <audio_dir>")
		os.Exit(1)
	}

	audioDir := os.Args[1]

	radio := radio.New(audioDir)
	err := radio.ScanAudioFiles()
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()

	r := http.NewServeMux()

	h := handler.NewHandler(ctx, radio)

	r.HandleFunc("GET /", handler.Make(h.IndexHandler))
	r.HandleFunc("GET /index.css", handler.Make(h.CSSHandler))
	r.HandleFunc("GET /favicon.png", handler.Make(h.FaviconHandler))

	r.HandleFunc("GET /api/radio/status", handler.Make(h.RadioStatusHandler))

	r.HandleFunc("GET /api/radio/files", handler.Make(h.RadioFileListHandler))
	r.HandleFunc("POST /api/radio/files", handler.Make(h.RadioFileUploadHandler))
	r.HandleFunc("DELETE /api/radio/files/{fileID}", handler.Make(h.RadioFileDeleteHandler))

	r.HandleFunc("POST /api/radio/play", handler.Make(h.RadioPlayHandler))
	r.HandleFunc("POST /api/radio/pause", handler.Make(h.RadioPauseHandler))
	r.HandleFunc("POST /api/radio/volume", handler.Make(h.RadioVolumeHandler))
	r.HandleFunc("POST /api/radio/stop", handler.Make(h.RadioStopHandler))

	r.HandleFunc("/", handler.Make(h.NotFoundHandler))

	stack := middleware.CreateStack(
		middleware.Log,
		middleware.CORS,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(r),
	}

	fmt.Println("Server listening on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
