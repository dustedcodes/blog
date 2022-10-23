package site

import "fmt"

type URLs struct {
	RequestURL      string
	BaseURL         string
	CDN             string
	DisqusShortname string
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
	return fmt.Sprintf("%s/images/public/dusted-codes-open-graph.jpg", u.CDN)
}
