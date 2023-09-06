package web

import (
	"net/http"
	"strings"

	"github.com/dusted-go/http/v5/htmlview"
	"github.com/dusted-go/http/v5/route"

	"github.com/dustedcodes/blog/cmd/blog/model"
	"github.com/dustedcodes/blog/internal/blog"
	"github.com/dustedcodes/blog/internal/config"
)

type Handler struct {
	config     *config.Config
	assets     *model.Assets
	viewWriter *htmlview.Writer
	blogPosts  []*blog.Post
}

func NewHandler(
	config *config.Config,
	assets *model.Assets,
	blobPosts []*blog.Post,
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
		"tagged": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/pages/tagged.html",
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
		"oss": append(socialSVGs,
			"dist/templates/pages/_layout.html",
			"dist/templates/svgs/link.svg",
			"dist/templates/pages/oss.html",
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
	viewWriter := htmlview.NewWriter(
		config.HotReload(),
		"layout",
		templateFiles)

	return &Handler{
		config:     config,
		assets:     assets,
		viewWriter: viewWriter,
		blogPosts:  blobPosts,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	verb := r.Method
	p := r.URL.Path

	if verb == "GET" || verb == "HEAD" {

		if p == "/" {
			h.index(w, r)
			return
		}

		if p == "/version" {
			h.version(w, r)
			return
		}

		if p == "/ping" {
			h.ping(w, r)
			return
		}

		if p == "/panic" && !h.config.IsProduction() {
			h.panic(w, r)
			return
		}

		if p == "/projects" {
			h.projects(w, r)
			return
		}

		if p == "/open-source" {
			h.oss(w, r)
			return
		}

		if p == "/hire" {
			h.hire(w, r)
			return
		}

		if p == "/about" {
			h.about(w, r)
			return
		}

		if p == "/feed/rss" {
			h.rss(w, r)
			return
		}

		if p == "/feed/atom" {
			h.atom(w, r)
			return
		}

		if p == "/sitemap.xml" {
			h.sitemap(w, r)
			return
		}

		if p == "/robots.xml" {
			h.robots(w, r)
			return
		}

		head, tail := route.ShiftPath(p)
		if head == "tagged" {
			tagName := strings.TrimLeft(tail, "/")
			h.tagged(w, r, tagName)
			return
		}

		// Support for legacy URLs:
		if head == "demystifying-aspnet-mvc-5-error-pages" {
			http.Redirect(
				w, r,
				"/demystifying-aspnet-mvc-5-error-pages-and-error-logging",
				http.StatusMovedPermanently)
			return
		}

		h.blogPost(w, r)
		return
	}

	h.notFound(w, r)
}
