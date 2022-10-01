package site

import "fmt"

type URLs struct {
	RequestURL string
	BaseURL    string
	CDN        string
}

func (u *URLs) BlogPostPermalink(blogPostID string) string {
	return fmt.Sprintf("%s/%s", u.BaseURL, blogPostID)
}
