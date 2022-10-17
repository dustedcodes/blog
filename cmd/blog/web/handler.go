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
	socialSVGs := []string{
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
	}

	templateFiles := map[string][]string{
		"index": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/index.html",
		),
		"message": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/message.html",
		),
		"blogPost": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/blogPost.html",
		),
		"projects": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/svgs/link.svg",
			"dist/templates/pages/projects.html",
		),
		"hire": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/hire.html",
		),
		"about": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/about.html",
		),
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

		if path == "/projects" {
			h.projects(w, r)
			return
		}

		if path == "/hire" {
			h.hire(w, r)
			return
		}

		if path == "/about" {
			h.about(w, r)
			return
		}

		h.blogPost(w, r)
		return
	}

	h.notFound(w, r)
}
