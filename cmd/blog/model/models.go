package model

import (
	"html/template"
	"sort"

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

type Empty struct {
	Base Base
}

type UserMessage struct {
	Base     Base
	Messages []template.HTML
}

type Index struct {
	Base        Base
	Catalogue   map[int][]string
	SortedYears []int
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
	catalogue := map[int][]string{}
	years := []int{}
	for _, post := range blogPosts {
		year := post.Year()
		if !array.Contains(years, year) {
			years = append(years, year)
		}

		catalogue[year] = array.Prepend(
			catalogue[year],
			post.Title)
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
