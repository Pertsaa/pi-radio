package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pertsaa/pi-radio/radio"
)

func (h *Handler) RadioFileListHandler(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, h.radio.AudioFiles)
}

func (h *Handler) RadioFileUploadHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(100)

	for _, fileHeader := range r.MultipartForm.File["files"] {
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		fileName := fileHeader.Filename
		fileExt := filepath.Ext(fileName)
		fileNameWithoutExt := fileName[:len(fileName)-len(fileExt)]
		lowerCaseExt := strings.ToLower(fileExt)
		newFileName := fileNameWithoutExt + lowerCaseExt

		dst, err := os.Create(filepath.Join(h.radio.AudioDir, newFileName))
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			return err
		}
	}

	h.radio.ScanAudioFiles()

	return writeJSON(w, http.StatusOK, h.radio.AudioFiles)
}

func (h *Handler) RadioFileDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	fileID := r.PathValue("fileID")

	var audioFile radio.AudioFile
	for _, file := range h.radio.AudioFiles {
		if file.ID == fileID {
			audioFile = file
		}
	}
	if audioFile.ID == "" {
		return fmt.Errorf("audio file not found: %s", fileID)
	}

	if err := os.Remove(fmt.Sprintf("%s/%s", h.radio.AudioDir, audioFile.Name)); err != nil {
		return err
	}

	h.radio.ScanAudioFiles()

	return writeJSON(w, http.StatusOK, h.radio.AudioFiles)
}

type AudioPlayBody struct {
	AudioFileID string `json:"audio_file_id"`
}

func (h *Handler) RadioPlayHandler(w http.ResponseWriter, r *http.Request) error {
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

func (h *Handler) RadioPauseHandler(w http.ResponseWriter, r *http.Request) error {
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

func (h *Handler) RadioVolumeHandler(w http.ResponseWriter, r *http.Request) error {
	body, err := parseBody[AudioVolumeBody](r)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, err)
	}

	h.radio.SetVolume(body.Volume)

	return writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) RadioStopHandler(w http.ResponseWriter, r *http.Request) error {
	h.radio.Stop()

	return writeJSON(w, http.StatusOK, nil)
}
