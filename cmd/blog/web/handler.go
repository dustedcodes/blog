package web

import (
	"net/http"
	"strings"

	"github.com/dusted-go/http/v6/htmlview"
	"github.com/dusted-go/http/v6/route"

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
	masterFiles := []string{
		"dist/templates/components/branding.html",
		"dist/templates/components/nav.html",
		"dist/templates/components/footer.html",
		"dist/templates/pages/_master.html",
		"dist/templates/svgs/logos/buymeacoffee.svg",
		"dist/templates/svgs/logos/github.svg",
		"dist/templates/svgs/logos/instagram.svg",
		"dist/templates/svgs/logos/linkedin.svg",
		"dist/templates/svgs/logos/paypal.svg",
		"dist/templates/svgs/logos/rssfeed.svg",
		"dist/templates/svgs/logos/x.svg",
		"dist/templates/svgs/logos/youtube.svg",
		"dist/templates/svgs/logos/logo.svg",
	}

	templateFiles := map[string][]string{
		"index": append(masterFiles,
			"dist/templates/svgs/illustrations/dustin-tshirt.svg",
			"dist/templates/pages/index.html",
		),
		"blog": append(masterFiles,
			"dist/templates/pages/_page.html",
			"dist/templates/svgs/illustrations/blogging.svg",
			"dist/templates/pages/blog.html",
		),
		"tagged": append(masterFiles,
			"dist/templates/pages/_page.html",
			"dist/templates/pages/tagged.html",
			"dist/templates/components/tags.html",
		),
		"404": append(masterFiles,
			"dist/templates/svgs/illustrations/404.svg",
			"dist/templates/pages/404.html",
		),
		"blogPost": append(masterFiles,
			"dist/templates/pages/_page.html",
			"dist/templates/pages/article.html",
			"dist/templates/components/tags.html",
		),
		"products": append(masterFiles,
			"dist/templates/pages/_page.html",
			"dist/templates/svgs/illustrations/link.svg",
			"dist/templates/svgs/logos/msgdrop.svg",
			"dist/templates/svgs/logos/dotnet-ketchup.svg",
			"dist/templates/pages/products.html",
		),
		"oss": append(masterFiles,
			"dist/templates/pages/_page.html",
			"dist/templates/svgs/illustrations/link.svg",
			"dist/templates/pages/oss.html",
		),
		"hire": append(masterFiles,
			"dist/templates/pages/_page.html",
			"dist/templates/svgs/illustrations/dustin-macbook.svg",
			"dist/templates/svgs/illustrations/collaboration.svg",
			"dist/templates/svgs/illustrations/devops.svg",
			"dist/templates/svgs/illustrations/opensource.svg",
			"dist/templates/svgs/illustrations/security.svg",
			"dist/templates/svgs/illustrations/tea.svg",
			"dist/templates/svgs/illustrations/training.svg",
			"dist/templates/pages/hire.html",
		),
		"about": append(masterFiles,
			"dist/templates/pages/_page.html",
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
	path := r.URL.Path

	if verb != "GET" && verb != "HEAD" {
		h.notFound(w, r)
		return
	}

	if path == "/" {
		h.index(w, r)
		return
	}

	if path == "/version" {
		h.version(w, r)
		return
	}

	if path == "/ping" {
		h.ping(w, r)
		return
	}

	if path == "/panic" && !h.config.IsProduction() {
		h.panic(w, r)
		return
	}

	if path == "/blog" {
		h.blog(w, r)
		return
	}

	if path == "/products" {
		h.products(w, r)
		return
	}

	if path == "/open-source" {
		h.oss(w, r)
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

	if path == "/feed/rss" {
		h.rss(w, r)
		return
	}

	if path == "/feed/atom" {
		h.atom(w, r)
		return
	}

	if path == "/sitemap.xml" {
		h.sitemap(w, r)
		return
	}

	if path == "/robots.txt" {
		h.robots(w, r)
		return
	}

	head, tail := route.ShiftPath(path)
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
}
