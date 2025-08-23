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
	r.HandleFunc("GET /favicon.ico", handler.Make(h.FaviconHandler))

	r.HandleFunc("GET /api/audio_files", handler.Make(h.AudioFileListHandler))
	r.HandleFunc("POST /api/audio_files/{audioFileID}/play", handler.Make(h.AudioPlayHandler))
	r.HandleFunc("POST /api/audio/pause", handler.Make(h.AudioPauseHandler))
	r.HandleFunc("POST /api/audio/volume", handler.Make(h.AudioVolumeHandler))
	r.HandleFunc("POST /api/audio/stop", handler.Make(h.AudioStopHandler))

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
