package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dusted-go/diagnostic/v3/dlog"
	"github.com/dusted-go/http/v3/server"
	"github.com/dusted-go/utils/array"
	"github.com/dustedcodes/blog/cmd/blog/rss"
	"github.com/dustedcodes/blog/cmd/blog/site"
	"github.com/dustedcodes/blog/cmd/blog/sitemap"
)

func (h *Handler) setCacheDirective(
	w http.ResponseWriter,
	cacheDuration int,
) {
	cacheDirective := fmt.Sprintf("public, max-age=%d", cacheDuration)
	w.Header().Add("Cache-Control", cacheDirective)
}

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
	err := server.WritePlaintext(
		w,
		http.StatusOK,
		h.settings.ApplicationVersion)
	h.handleErr(w, r, err)
}

func (h *Handler) ping(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := server.WritePlaintext(
		w,
		http.StatusOK,
		"pong")
	h.handleErr(w, r, err)
}

func (h *Handler) index(
	w http.ResponseWriter,
	r *http.Request,
) {
	model := h.newBaseModel(r).Index(h.blogPosts)
	h.setCacheDirective(w, 60*60*4)
	h.renderView(w, r, 200, "index", model)
}

func (h *Handler) tagged(
	w http.ResponseWriter,
	r *http.Request,
	tagName string,
) {
	filtered := []*site.BlogPost{}
	for _, b := range h.blogPosts {
		if array.Contains(b.Tags, tagName) {
			filtered = append(filtered, b)
		}
	}
	model := h.newBaseModel(r).WithTitle(fmt.Sprintf("Tagged with '%s'", tagName)).Tagged(filtered)
	h.setCacheDirective(w, 60*60*4)
	h.renderView(w, r, 200, "tagged", model)
}

func (h *Handler) renderBlogPost(
	w http.ResponseWriter,
	r *http.Request,
	b *site.BlogPost,
) {
	// Parse content:
	// ---
	content, err := b.HTML()
	if h.handleErr(w, r, err) {
		return
	}

	// Respond with view:
	// ---
	m := h.newBaseModel(r).WithTitle(b.Title).BlogPost(b.ID, content, b.PublishDate, b.Tags)
	h.setCacheDirective(w, 60*60*4)
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
}

func (h *Handler) projects(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24)
	h.renderView(w, r, 200, "projects", h.newBaseModel(r).WithTitle("Projects").Empty())
}

func (h *Handler) oss(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24)
	h.renderView(w, r, 200, "oss", h.newBaseModel(r).WithTitle("Open Source").Empty())
}

func (h *Handler) hire(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24)
	h.renderView(w, r, 200, "hire", h.newBaseModel(r).WithTitle("Hire").Empty())
}

func (h *Handler) about(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24)
	h.renderView(w, r, 200, "about", h.newBaseModel(r).WithTitle("About").Empty())
}

func (h *Handler) rss(
	w http.ResponseWriter,
	r *http.Request,
) {
	urls := h.settings.URLs(r)
	latestPost := h.blogPosts[0]
	rssFeed := rss.NewFeed(
		rss.NewChannel(
			"Dusted Codes",
			urls.BaseURL,
			"Programming Adventures").
			SetLanguage("en-gb").
			SetWebMaster("dustin@dusted.codes", "Dustin Moris Gorski").
			SetManagingEditor("dustin@dusted.codes", "Dustin Moris Gorski").
			SetCopyright(fmt.Sprintf("Copyright %d, Dustin Moris Gorski", time.Now().Year())).
			SetLastBuildDate(latestPost.PublishDate).
			SetPubDate(latestPost.PublishDate),
	)

	for _, b := range h.blogPosts {
		permalink := urls.BlogPostURL(b.ID)
		comments := urls.BlogPostCommentsURL(b.ID)
		htmlContent, err := b.HTML()
		if h.handleErr(w, r, err) {
			return
		}
		rssItem := rss.NewItemWithTitle(b.Title).
			SetLink(permalink).
			SetGUID(permalink, true).
			SetPubDate(b.PublishDate).
			SetAuthor("dustin@dusted.codes", "Dustin Moris Gorski").
			SetComments(comments).
			SetDescription(string(htmlContent))
		for _, t := range b.Tags {
			rssItem.AddCategory(t, urls.TagURL(t))
		}
		rssFeed.Channel.AddItem(rssItem)
	}

	bytes, err := rssFeed.ToXML(true, true)
	if h.handleErr(w, r, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/rss+xml")
	_, err = w.Write(bytes)
	if err != nil {
		dlog.New(r.Context()).Critical().Err(err).Msg("Error writing rss feed to response body.")
	}
}

func (h *Handler) sitemap(
	w http.ResponseWriter,
	r *http.Request,
) {
	urls := h.settings.URLs(r)
	urlset := sitemap.NewURLSet().
		AddURL(
			sitemap.
				NewURL(urls.BaseURL).
				SetPriority("1").
				SetChangeFreq("monthly")).
		AddURL(
			sitemap.
				NewURL(urls.Projects()).
				SetPriority("0.9").
				SetChangeFreq("monthly")).
		AddURL(
			sitemap.
				NewURL(urls.OpenSource()).
				SetPriority("0.9").
				SetChangeFreq("monthly")).
		AddURL(
			sitemap.
				NewURL(urls.Hire()).
				SetPriority("0.9").
				SetChangeFreq("monthly")).
		AddURL(
			sitemap.
				NewURL(urls.About()).
				SetPriority("0.9").
				SetChangeFreq("monthly"))

	for _, blogPost := range h.blogPosts {
		urlset.AddURL(
			sitemap.
				NewURL(urls.BlogPostURL(blogPost.ID)).
				SetPriority("0.9").
				SetChangeFreq("monthly").
				SetLastMod(blogPost.PublishDate))
	}

	bytes, err := urlset.ToXML(true, true)
	if h.handleErr(w, r, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/xml; charset=UTF-8")
	_, err = w.Write(bytes)
	if err != nil {
		dlog.New(r.Context()).Critical().Err(err).Msg("Error writing sitemap to response body.")
	}
}

func (h *Handler) robots(
	w http.ResponseWriter,
	r *http.Request,
) {
	contents := fmt.Sprintf("Sitemap: %s/sitemap.xml\n", h.settings.BaseURL)
	err := server.WritePlaintext(
		w,
		http.StatusOK,
		contents)
	h.handleErr(w, r, err)
}
