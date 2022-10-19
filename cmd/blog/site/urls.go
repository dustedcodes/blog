package site

import "fmt"

type URLs struct {
	RequestURL      string
	BaseURL         string
	CDN             string
	DisqusShortname string
}

func (u *URLs) BlogPostPermalink(blogPostID string) string {
	return fmt.Sprintf("%s/%s", u.BaseURL, blogPostID)
}

func (u *URLs) DisqusCountScript() string {
	return fmt.Sprintf("//%s.disqus.com/count.js", u.DisqusShortname)
}
