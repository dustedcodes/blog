package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dusted-go/http/v3/response"
	"github.com/dustedcodes/blog/cmd/blog/site"
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
	model := h.newBaseModel(r, "Dusted Codes").Index(h.blogPosts)
	h.renderView(w, r, 200, "index", model)
}

func (h *Handler) renderBlogPost(
	w http.ResponseWriter,
	r *http.Request,
	b *site.BlogPost,
) {

	// Parse Markdown:
	// ---
	parsed, err := b.ParsedMarkdown()
	if h.handleErr(w, r, err) {
		return
	}

	// Set Cache directive:
	// ---
	// 60sec * 60 * 24 = 1 day
	cacheDuration := 60 * 60 * 24
	cacheDirective := fmt.Sprintf("public, max-age=%d", cacheDuration)
	w.Header().Add("Cache-Control", cacheDirective)

	// Respond with view:
	// ---
	m := h.newBaseModel(r, b.Title).BlogPost(b.ID, parsed, b.PublishDate, b.Tags)
	h.renderView(
		w, r, 200, "blogPost", m)
}

func (h *Handler) blogPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	blogPostID := strings.TrimPrefix(r.URL.Path, "/")

	for _, blogPost := range h.blogPosts {
		if blogPost.ID == blogPostID {
			h.renderBlogPost(w, r, blogPost)
			return
		}
	}

	h.notFound(w, r)
	return
}
