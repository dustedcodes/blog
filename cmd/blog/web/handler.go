package web

import (
	"net/http"

	"github.com/dusted-go/http/v3/server"
	"github.com/dustedcodes/blog/cmd/blog/site"
)

type Handler struct {
	settings    *site.Settings
	assets      *site.Assets
	viewHandler *server.ViewHandler
}

func NewHandler(
	settings *site.Settings,
	assets *site.Assets,
) *Handler {
	templateFiles := map[string][]string{
		"index": {
			"dist/templates/_layout.html",
			"dist/templates/index.html",
		},
		"message": {
			"dist/templates/_layout.html",
			"dist/templates/message.html",
		},
	}
	viewHandler := server.NewViewHandler(
		settings.HotReload(),
		"layout",
		templateFiles)

	return &Handler{
		settings:    settings,
		assets:      assets,
		viewHandler: viewHandler,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	verb := r.Method
	path := r.URL.Path

	if verb == "GET" || verb == "HEAD" {

		if path == "/" {
			h.index(w, r)
			return
		}

		if path == "/version" {
			h.version(w, r)
			return
		}

		if path == "/panic" && !h.settings.IsProduction() {
			h.panic(w, r)
			return
		}
	}

	h.notFound(w, r)
}
