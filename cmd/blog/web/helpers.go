package web

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dusted-go/diagnostic/v3/dlog"
	"github.com/dusted-go/fault/fault"
	"github.com/dusted-go/fault/stack"
	"github.com/dusted-go/http/v3/request"
	"github.com/dusted-go/http/v3/server"
	"github.com/dustedcodes/blog/cmd/blog/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertErrorsToHTML(errorMessages []string) template.HTML {
	out := "<div class=\"error\"><p>We've encountered some errors with your request:</p>"
	for _, msg := range errorMessages {
		out = out + fmt.Sprintf("<p>%s</p>", msg)
	}
	out = out + "</div>"
	// nolint: gosec // System generated error messages
	return template.HTML(out)
}

func (h *Handler) newBaseModel(r *http.Request) model.Base {
	return model.Base{
		Title:    "Dusted Codes",
		SubTitle: "Programming Adventures",
		Year:     time.Now().Year(),
		Assets:   h.assets,
		URLs:     h.settings.URLs(r),
	}
}

func (h *Handler) internalError(
	w http.ResponseWriter,
	r *http.Request,
) {
	server.ClearHeaders(w)
	err := server.WritePlaintext(
		w,
		http.StatusInternalServerError,
		"Oops, something went wrong. The server encountered an internal error or misconfiguration and was unable to complete your request.")
	if err != nil {
		dlog.New(r.Context()).
			Err(err).
			Critical().
			Msg("Error sending 'Internal Server Error' response.")
	}
}

func (h *Handler) renderView(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	viewKey string,
	viewModel any,
) {
	err := h.viewHandler.WriteView(
		w,
		statusCode,
		viewKey,
		viewModel)
	if err != nil {
		dlog.New(r.Context()).
			Critical().
			Err(err).
			Fmt("Failed to render view with key '%s'", viewKey)
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
	server.ClearHeaders(w)
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
	if errors.Is(err, context.Canceled) {
		return true
	}
	if grpcStatus, ok := fault.As(err, status.FromError); ok {
		if grpcStatus != nil && grpcStatus.Code() == codes.Canceled {
			return true
		}
	}

	var userErr *fault.UserError
	if errors.As(err, &userErr) {
		h.renderUserMessages(
			w, r, 400,
			"Bad Request",
			convertErrorsToHTML(userErr.ErrorMessages()))
		return true
	}

	dlog.New(r.Context()).
		Critical().
		Err(err).
		Msg("An unexpected error occurred.")
	h.internalError(w, r)
	return true
}

func (h *Handler) notFound(
	w http.ResponseWriter,
	r *http.Request,
) {
	dlog.New(r.Context()).
		Debug().
		Fmt("Not Found: %s", request.FullURL(r))
	server.ClearHeaders(w)
	h.renderUserMessages(
		w, r,
		http.StatusNotFound,
		"Page not found",
		"Sorry, the page you have requested may have been moved or deleted.")
}

// func (h *Handler) methodNotAllowed(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {
// 	server.ClearHeaders(w)
// 	err := server.WritePlaintext(
// 		w,
// 		http.StatusMethodNotAllowed,
// 		fmt.Sprintf("The HTTP method '%s' is not allowed on this path.", r.Method))
// 	if err != nil {
// 		dlog.New(r.Context()).
// 			Err(err).
// 			Critical().
// 			Msg("Error sending 'Method Not Allowed' response.")
// 	}
// }

func (h *Handler) Recover(recovered any, stackTrace stack.Trace) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("%v\n\n%v", recovered, stackTrace.String())
		dlog.New(r.Context()).Critical().Fmt("Application panicked with error:\n\n%v", msg)
		h.internalError(w, r)
	}
}
