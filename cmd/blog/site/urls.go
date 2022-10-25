package site

import "fmt"

type URLs struct {
	RequestURL      string
	BaseURL         string
	CDN             string
	DisqusShortname string
}

func (u *URLs) Projects() string {
	return fmt.Sprintf("%s/projects", u.BaseURL)
}

func (u *URLs) OpenSource() string {
	return fmt.Sprintf("%s/open-source", u.BaseURL)
}

func (u *URLs) Hire() string {
	return fmt.Sprintf("%s/hire", u.BaseURL)
}

func (u *URLs) About() string {
	return fmt.Sprintf("%s/about", u.BaseURL)
}

func (u *URLs) RSSFeed() string {
	return fmt.Sprintf("%s/feed/rss", u.BaseURL)
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

func (u *URLs) OpenGraphImage() string {
	return fmt.Sprintf("%s/images/public/dusted-codes-open-graph-small.jpg", u.CDN)
}

func (u *URLs) Logo() string {
	return fmt.Sprintf("%s/images/public/logo.png", u.CDN)
}
