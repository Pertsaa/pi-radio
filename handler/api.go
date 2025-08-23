package handler

import (
	"net/http"
)

func (h *Handler) AudioFileListHandler(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, h.radio.AudioFiles)
}

type AudioPlayBody struct {
	AudioFileID string `json:"audio_file_id"`
}

func (h *Handler) AudioPlayHandler(w http.ResponseWriter, r *http.Request) error {
	body, err := parseBody[AudioPlayBody](r)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, err)
	}

	err = h.radio.Play(body.AudioFileID)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, nil)
}

type AudioPauseBody struct {
	Paused bool `json:"paused"`
}

func (h *Handler) AudioPauseHandler(w http.ResponseWriter, r *http.Request) error {
	body, err := parseBody[AudioPauseBody](r)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, err)
	}

	h.radio.SetPaused(body.Paused)

	return writeJSON(w, http.StatusOK, nil)
}

type AudioVolumeBody struct {
	Volume float64 `json:"volume"`
}

func (h *Handler) AudioVolumeHandler(w http.ResponseWriter, r *http.Request) error {
	body, err := parseBody[AudioVolumeBody](r)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, err)
	}

	h.radio.SetVolume(body.Volume)

	return writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) AudioStopHandler(w http.ResponseWriter, r *http.Request) error {
	h.radio.Stop()

	return writeJSON(w, http.StatusOK, nil)
}
