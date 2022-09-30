package model

import (
	"html/template"

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
