package handler

import (
	_ "embed"
	"net/http"

	"github.com/Pertsaa/pi-radio/static"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	w.Write(static.IndexHTML)
	return nil
}

func (h *Handler) CSSHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/css")
	w.Write(static.IndexCSS)
	return nil
}

func (h *Handler) FaviconHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(""))
	return nil
}
