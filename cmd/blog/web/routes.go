package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dusted-go/http/v5/atom"
	"github.com/dusted-go/http/v5/rss"
	"github.com/dusted-go/http/v5/sitemap"

	"github.com/dustedcodes/blog/internal/array"
	"github.com/dustedcodes/blog/internal/blog"
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
	h.writeText(
		w, r,
		http.StatusOK,
		h.config.ApplicationVersion)
}

func (h *Handler) ping(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.writeText(
		w, r,
		http.StatusOK,
		"pong")
}

func (h *Handler) index(
	w http.ResponseWriter,
	r *http.Request,
) {
	model := h.newBaseModel(r).Index(h.blogPosts)
	h.setCacheDirective(w, 60*60, h.config.ApplicationVersion)
	h.renderView(w, r, 200, "index", model)
}

func (h *Handler) tagged(
	w http.ResponseWriter,
	r *http.Request,
	tagName string,
) {
	filtered := []*blog.Post{}
	for _, b := range h.blogPosts {
		if array.Contains(b.Tags, tagName) {
			filtered = append(filtered, b)
		}
	}
	model := h.newBaseModel(r).WithTitle(fmt.Sprintf("Tagged with '%s'", tagName)).Tagged(filtered)
	h.setCacheDirective(w, 60*60*4, h.config.ApplicationVersion)
	h.renderView(w, r, 200, "tagged", model)
}

func (h *Handler) renderBlogPost(
	w http.ResponseWriter,
	r *http.Request,
	b *blog.Post,
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
	h.setCacheDirective(w, 60*60*4, b.HashCode)
	h.renderView(
		w, r, 200, "blogPost", m)
}

func (h *Handler) blogPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	blogPostID := strings.TrimPrefix(r.URL.Path, "/")

	if !h.config.IsProduction() {
		blogPost, err := blog.ReadPost(blog.DefaultBlogPostPath, blogPostID)
		if h.handleErr(w, r, err) {
			return
		}
		h.renderBlogPost(w, r, blogPost)
		return
	}

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
	h.setCacheDirective(w, 60*60*24, h.config.ApplicationVersion)
	h.renderView(w, r, 200, "projects", h.newBaseModel(r).WithTitle("Projects").Empty())
}

func (h *Handler) oss(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24, h.config.ApplicationVersion)
	h.renderView(w, r, 200, "oss", h.newBaseModel(r).WithTitle("Open Source").Empty())
}

func (h *Handler) hire(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24, h.config.ApplicationVersion)
	h.renderView(w, r, 200, "hire", h.newBaseModel(r).WithTitle("Hire").Empty())
}

func (h *Handler) about(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.setCacheDirective(w, 60*60*24, h.config.ApplicationVersion)
	h.renderView(w, r, 200, "about", h.newBaseModel(r).WithTitle("About").Empty())
}

func (h *Handler) rss(
	w http.ResponseWriter,
	r *http.Request,
) {
	loc, err := time.LoadLocation("Europe/London")
	if h.handleErr(w, r, err) {
		return
	}
	urls := h.getURLs(r)
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
			SetLastBuildDate(latestPost.PublishDate, loc).
			SetPubDate(latestPost.PublishDate, loc).
			SetImage(rss.NewImage(urls.Logo(), "Dusted Codes", urls.BaseURL)),
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
			SetPubDate(b.PublishDate, loc).
			SetAuthor("dustin@dusted.codes", "Dustin Moris Gorski").
			SetComments(comments).
			SetDescription(string(htmlContent)).
			SetEnclosure(urls.OpenGraphImage(), 28*1024, "image/png")
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
	_, _ = w.Write(bytes)
}

func (h *Handler) atom(
	w http.ResponseWriter,
	r *http.Request,
) {
	urls := h.getURLs(r)
	latestPost := h.blogPosts[0]
	author := atom.NewPerson(
		"Dustin Moris Gorski").
		SetEmail("dustin@dusted.codes").
		SetURI(urls.BaseURL)
	atomFeed := atom.NewFeed(
		urls.BaseURL,
		atom.NewText("Dusted Codes"),
		latestPost.PublishDate).
		SetSubtitle(atom.NewText("Programming Adventures")).
		SetIcon(urls.Logo()).
		SetAuthor(author).
		AddLink(atom.NewLink(urls.AtomFeed()).SetRel("self")).
		AddLink(atom.NewLink(urls.BaseURL).SetRel("alternate")).
		SetRights(
			atom.NewText(
				fmt.Sprintf("Copyright © %d, Dusted Codes Limited", time.Now().Year())))

	for _, b := range h.blogPosts {
		permalink := urls.BlogPostURL(b.ID)
		htmlContent, err := b.HTML()
		if h.handleErr(w, r, err) {
			return
		}
		entry := atom.NewEntry(
			permalink,
			atom.NewText(b.Title),
			b.PublishDate).
			SetAuthor(author).
			AddLink(atom.NewLink(permalink).SetRel("alternate")).
			AddLink(atom.NewLink(urls.BlogPostCommentsURL(b.ID)).SetRel("related")).
			AddLink(atom.NewLink(urls.OpenGraphImage()).SetRel("enclosure").SetLength(28 * 1024)).
			SetPublished(b.PublishDate).
			SetContent(atom.NewHTML(string(htmlContent)))

		for _, t := range b.Tags {
			entry.AddCategory(atom.NewCategory(t).
				SetLabel(t).
				SetScheme(urls.TagURL(t)))
		}
		atomFeed.AddEntry(entry)
	}

	bytes, err := atomFeed.ToXML(true, true)
	if h.handleErr(w, r, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/atom+xml")
	_, _ = w.Write(bytes)
}

func (h *Handler) sitemap(
	w http.ResponseWriter,
	r *http.Request,
) {
	urls := h.getURLs(r)
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
	_, _ = w.Write(bytes)
}

func (h *Handler) robots(
	w http.ResponseWriter,
	r *http.Request,
) {
	contents := fmt.Sprintf("Sitemap: %s/sitemap.xml\n", h.config.BaseURL)
	h.writeText(
		w, r,
		http.StatusOK,
		contents)
}
