package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) AudioFileListHandler(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, h.radio.AudioFiles)
}

func (h *Handler) AudioPlayHandler(w http.ResponseWriter, r *http.Request) error {
	audioFileID := r.PathValue("audioFileID")

	if audioFileID == "" {
		return fmt.Errorf("invalid file ID: %s", audioFileID)
	}

	err := h.radio.Play(audioFileID)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) AudioPauseHandler(w http.ResponseWriter, r *http.Request) error {
	h.radio.SetPaused(false)

	return writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) AudioVolumeHandler(w http.ResponseWriter, r *http.Request) error {
	h.radio.SetVolume(0)

	return writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) AudioStopHandler(w http.ResponseWriter, r *http.Request) error {
	h.radio.Stop()

	return writeJSON(w, http.StatusOK, nil)
}
