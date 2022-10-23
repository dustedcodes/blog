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
	Title           string
	SubTitle        string
	Year            int
	Assets          *site.Assets
	URLs            *site.URLs
	DisqusShortname string
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

type Tag struct {
	Value string
	URL   string
}

type BlogPostLink struct {
	Title       string
	Permalink   string
	PublishDate time.Time
	Tags        []Tag
}

func (b BlogPostLink) PublishedOn() string {
	return b.PublishDate.Format("02 Jan 2006")
}

type Index struct {
	Base        Base
	Catalogue   map[int][]BlogPostLink
	SortedYears []int
}

type BlogPost struct {
	Base             Base
	ID               string
	PublishDate      time.Time
	Content          template.HTML
	Tags             []Tag
	EncodedTitle     string
	Permalink        string
	EncodedPermalink string
}

type Tagged struct {
	Base      Base
	BlogPosts []BlogPostLink
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
		tags := []Tag{}
		for _, tag := range post.Tags {
			tags = append(tags, Tag{
				Value: tag,
				URL:   b.URLs.TagURL(tag),
			})
		}
		catalogue[year] = append(
			catalogue[year],
			BlogPostLink{
				Title:       post.Title,
				Permalink:   b.URLs.BlogPostURL(post.ID),
				PublishDate: post.PublishDate,
				Tags:        tags,
			})
	}
	sort.Slice(years, func(i, j int) bool {
		return years[i] > years[j]
	})

	for _, year := range years {
		sort.Slice(catalogue[year], func(i, j int) bool {
			return catalogue[year][i].PublishDate.After(catalogue[year][j].PublishDate)
		})
	}

	return Index{
		Base:        b,
		Catalogue:   catalogue,
		SortedYears: years,
	}
}

func (b Base) Tagged(blogPosts []*site.BlogPost) Tagged {
	blogPostLinks := []BlogPostLink{}

	for _, post := range blogPosts {
		tags := []Tag{}
		for _, tag := range post.Tags {
			tags = append(tags, Tag{
				Value: tag,
				URL:   b.URLs.TagURL(tag),
			})
		}
		blogPostLinks = append(blogPostLinks, BlogPostLink{
			Title:       post.Title,
			Permalink:   b.URLs.BlogPostURL(post.ID),
			PublishDate: post.PublishDate,
			Tags:        tags,
		})
	}

	sort.Slice(blogPostLinks, func(i, j int) bool {
		return blogPostLinks[i].PublishDate.After(blogPostLinks[j].PublishDate)
	})

	return Tagged{
		Base:      b,
		BlogPosts: blogPostLinks,
	}
}

func (b Base) BlogPost(
	blogPostID string,
	content template.HTML,
	publishDate time.Time,
	tags []string,
) BlogPost {
	permalink := b.URLs.BlogPostURL(blogPostID)
	t := []Tag{}
	for _, tag := range tags {
		t = append(t, Tag{
			Value: tag,
			URL:   b.URLs.TagURL(tag),
		})
	}
	return BlogPost{
		Base:             b,
		ID:               blogPostID,
		PublishDate:      publishDate,
		Content:          content,
		Tags:             t,
		EncodedTitle:     url.QueryEscape(b.Title),
		Permalink:        permalink,
		EncodedPermalink: url.QueryEscape(permalink),
	}
}
