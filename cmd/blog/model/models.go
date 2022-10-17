package model

import (
	"html/template"
	"net/url"
	"sort"
	"time"

	"github.com/dusted-go/utils/array"
	"github.com/dustedcodes/blog/cmd/blog/site"
)

type Base struct {
	Title    string
	SubTitle string
	Year     int
	Assets   *site.Assets
	URLs     *site.URLs
}

func (b Base) WithTitle(title string) Base {
	b.Title = title
	return b
}

type Empty struct {
	Base Base
}

type UserMessage struct {
	Base     Base
	Messages []template.HTML
}

type BlogPostLink struct {
	Title     string
	Permalink string
}

type Index struct {
	Base        Base
	Catalogue   map[int][]BlogPostLink
	SortedYears []int
}

type BlogPost struct {
	Base             Base
	PublishDate      time.Time
	Content          template.HTML
	Tags             []string
	EncodedTitle     string
	Permalink        string
	EncodedPermalink string
}

func (b BlogPost) PublishedOn() string {
	return b.PublishDate.Format("02 Jan 2006")
}

func (b BlogPost) PublishedOnMachineReadable() string {
	return b.PublishDate.Format("2006-01-02T15:04:05")
}

func (b Base) Empty() Empty {
	return Empty{Base: b}
}

func (b Base) UserMessage(msg template.HTML) UserMessage {
	return UserMessage{
		Base:     b,
		Messages: []template.HTML{msg},
	}
}

func (b Base) UserMessages(msgs ...template.HTML) UserMessage {
	return UserMessage{
		Base:     b,
		Messages: msgs,
	}
}

func (b Base) Index(blogPosts []*site.BlogPost) Index {
	catalogue := map[int][]BlogPostLink{}
	years := []int{}

	for _, post := range blogPosts {
		year := post.Year()
		if !array.Contains(years, year) {
			years = append(years, year)
		}
		catalogue[year] = array.Prepend(
			catalogue[year],
			BlogPostLink{
				Title:     post.Title,
				Permalink: b.URLs.BlogPostPermalink(post.ID),
			})
	}
	sort.Slice(years, func(i, j int) bool {
		return years[i] > years[j]
	})

	return Index{
		Base:        b,
		Catalogue:   catalogue,
		SortedYears: years,
	}
}

func (b Base) BlogPost(
	blogPostID string,
	content template.HTML,
	publishDate time.Time,
	tags []string,
) BlogPost {
	permalink := b.URLs.BlogPostPermalink(blogPostID)
	return BlogPost{
		Base:             b,
		PublishDate:      publishDate,
		Content:          content,
		Tags:             tags,
		EncodedTitle:     url.QueryEscape(b.Title),
		Permalink:        permalink,
		EncodedPermalink: url.QueryEscape(permalink),
	}
}
