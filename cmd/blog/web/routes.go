package web

import (
	"net/http"

	"github.com/dusted-go/http/v3/response"
)

func (h *Handler) panic(
	_ http.ResponseWriter,
	_ *http.Request,
) {
	panic("boom crash burn")
}

func (h *Handler) version(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := response.Plaintext(
		true,
		http.StatusOK,
		h.settings.ApplicationVersion,
		w, r)
	h.handleErr(w, r, err)
}

func (h *Handler) index(
	w http.ResponseWriter,
	r *http.Request,
) {
	model := h.newBaseModel(r, "Dusted Codes").Empty()
	h.renderView(w, r, 200, "index", model)
}
