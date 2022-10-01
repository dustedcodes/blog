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
	blogPosts   []*site.BlogPost
}

func NewHandler(
	settings *site.Settings,
	assets *site.Assets,
	blobPosts []*site.BlogPost,
) *Handler {
	templateFiles := map[string][]string{
		"index": {
			"dist/templates/svgs/buymeacoffee.svg",
			"dist/templates/svgs/docker.svg",
			"dist/templates/svgs/github.svg",
			"dist/templates/svgs/instagram.svg",
			"dist/templates/svgs/linkedin.svg",
			"dist/templates/svgs/paypal.svg",
			"dist/templates/svgs/rssfeed.svg",
			"dist/templates/svgs/stackoverflow.svg",
			"dist/templates/svgs/twitter.svg",
			"dist/templates/svgs/youtube.svg",
			"dist/templates/svgs/logo.svg",
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/index.html",
		},
		"message": {
			"dist/templates/svgs/buymeacoffee.svg",
			"dist/templates/svgs/docker.svg",
			"dist/templates/svgs/github.svg",
			"dist/templates/svgs/instagram.svg",
			"dist/templates/svgs/linkedin.svg",
			"dist/templates/svgs/paypal.svg",
			"dist/templates/svgs/rssfeed.svg",
			"dist/templates/svgs/stackoverflow.svg",
			"dist/templates/svgs/twitter.svg",
			"dist/templates/svgs/youtube.svg",
			"dist/templates/svgs/logo.svg",
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/message.html",
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
		blogPosts:   blobPosts,
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
