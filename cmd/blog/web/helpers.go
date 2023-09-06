package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dusted-go/logging/stackdriver"

	"github.com/dustedcodes/blog/cmd/blog/model"
)

func (h *Handler) getURLs(r *http.Request) *model.URLs {
	return &model.URLs{
		RequestURL:      r.URL.Redacted(),
		BaseURL:         h.config.BaseURL,
		CDN:             h.config.CDN,
		DisqusShortname: h.config.DisqusShortname,
	}
}

func (h *Handler) newBaseModel(r *http.Request) model.Base {
	return model.Base{
		Title:           "Dusted Codes",
		SubTitle:        "Programming Adventures",
		Year:            time.Now().Year(),
		Assets:          h.assets,
		URLs:            h.getURLs(r),
		DisqusShortname: h.config.DisqusShortname,
	}
}

func (h *Handler) writeText(w http.ResponseWriter, r *http.Request, statusCode int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := fmt.Fprint(w, text)
	if err != nil {
		stackdriver.GetLogger(r.Context()).Error(
			"Failed to write text/plain response.",
			"error", err,
			"response", text)
	}
}

func (h *Handler) internalError(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.writeText(w, r,
		http.StatusInternalServerError,
		"Oops, something went wrong. The server encountered an internal error or misconfiguration and was unable to complete your request.")
}

func (h *Handler) renderView(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	viewKey string,
	viewModel any,
) {
	err := h.viewWriter.WriteView(
		w,
		statusCode,
		viewKey,
		viewModel)
	if err != nil {
		stackdriver.GetLogger(r.Context()).Error(
			"Failed to write html view response.",
			"error", err,
			"view", viewKey)
		h.internalError(w, r)
	}
}

func (h *Handler) renderUserMessages(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	title string,
	messages ...template.HTML,
) {
	model := h.newBaseModel(r).WithTitle(title).UserMessages(messages...)
	h.renderView(
		w, r,
		statusCode,
		"message",
		model)
}

func (h *Handler) handleErr(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) bool {
	if err == nil {
		return false
	}
	stackdriver.GetLogger(r.Context()).Error(
		"An unexpected error occurred.",
		"error", err)
	h.internalError(w, r)
	return true
}

func (h *Handler) notFound(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.renderUserMessages(
		w, r,
		http.StatusNotFound,
		"Page not found",
		"Sorry, the page you have requested may have been moved or deleted.")
}

func (h *Handler) setCacheDirective(
	w http.ResponseWriter,
	cacheDuration int,
	eTag string,
) {
	cacheDirective := fmt.Sprintf("public, max-age=%d", cacheDuration)
	w.Header().Add("Cache-Control", cacheDirective)
	w.Header().Add("ETag", fmt.Sprintf("\"%s\"", eTag))
}

func (h *Handler) Recover(recovered any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("Panic: %+v", recovered)
		stackdriver.GetLogger(r.Context()).Error(
			"Application panicked.",
			"error", err)
		h.internalError(w, r)
	}
}
