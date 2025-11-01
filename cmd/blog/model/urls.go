package model

import "fmt"

type URLs struct {
	RequestURL      string
	BaseURL         string
	CDN             string
	DisqusShortname string
}

func (u *URLs) Products() string {
	return u.BaseURL + "/products"
}

func (u *URLs) OpenSource() string {
	return u.BaseURL + "/open-source"
}

func (u *URLs) Hire() string {
	return u.BaseURL + "/hire"
}

func (u *URLs) About() string {
	return u.BaseURL + "/about"
}

func (u *URLs) RSSFeed() string {
	return u.BaseURL + "/feed/rss"
}

func (u *URLs) AtomFeed() string {
	return u.BaseURL + "/feed/atom"
}

func (u *URLs) BlogPostURL(blogPostID string) string {
	return fmt.Sprintf("%s/%s", u.BaseURL, blogPostID)
}

func (u *URLs) BlogPostCommentsURL(blogPostID string) string {
	return u.BlogPostURL(blogPostID) + "#disqus_thread"
}

func (u *URLs) TagURL(tagName string) string {
	return fmt.Sprintf("%s/tagged/%s", u.BaseURL, tagName)
}

func (u *URLs) DisqusCountScript() string {
	return fmt.Sprintf("//%s.disqus.com/count.js", u.DisqusShortname)
}

func (u *URLs) Logo() string {
	return u.CDN + "/images/public/logo.png"
}
