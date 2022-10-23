package sitemap

import (
	"encoding/xml"
	"fmt"
	"time"
)

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

func NewURL(loc string) *URL {
	return &URL{Loc: loc}
}

// SetLastMod sets the last modification date of the URL.
func (u *URL) SetLastMod(lastMod time.Time) *URL {
	u.LastMod = lastMod.Format("2006-01-02")
	return u
}

// SetChangeFreq sets the change frequency of the URL.
// Valid values are: always, hourly, daily, weekly, monthly, yearly, never.
func (u *URL) SetChangeFreq(changeFreq string) *URL {
	u.ChangeFreq = changeFreq
	return u
}

// SetPriority sets the priority of the URL.
// Valid values are: 0.0 to 1.0.
func (u *URL) SetPriority(priority string) *URL {
	u.Priority = priority
	return u
}

type URLSet struct {
	XMLName string `xml:"urlset"`
	XMLNS   string `xml:"xmlns,attr"`
	URLs    []*URL `xml:"url"`
}

// NewURLSet creates a new URLSet with the standard Sitemap namespace.
func NewURLSet() *URLSet {
	return &URLSet{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
}

func (u *URLSet) AddURL(url *URL) *URLSet {
	u.URLs = append(u.URLs, url)
	return u
}

func (u *URLSet) toXML(indent bool) ([]byte, error) {
	if indent {
		// nolint: wrapcheck
		return xml.MarshalIndent(u, "", "  ")
	}
	// nolint: wrapcheck
	return xml.Marshal(u)
}

// ToXML returns the XML representation of the feed.
func (u *URLSet) ToXML(indent, includeHeader bool) ([]byte, error) {
	bytes, err := u.toXML(indent)
	if err != nil {
		return nil, fmt.Errorf("error XML encoding Sitemap: %w", err)
	}
	if includeHeader {
		return append([]byte(xml.Header), bytes...), nil
	}
	return bytes, nil
}
