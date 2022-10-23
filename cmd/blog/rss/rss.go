package rss

import (
	"encoding/xml"
	"fmt"
	"html"
	"time"
)

type Image struct {
	// Required elements
	URL   string `xml:"url,omitempty"`
	Title string `xml:"title,omitempty"`
	Link  string `xml:"link,omitempty"`

	// Optional elements
	Width       int    `xml:"width,omitempty"`
	Height      int    `xml:"height,omitempty"`
	Description string `xml:"description,omitempty"`
}

// NewImage specifies a GIF, JPEG or PNG image that can be displayed with the channel.
// URL is the URL of a GIF, JPEG or PNG image that represents the channel.
// Title describes the image, it's used in the ALT attribute of the HTML <img> tag when the channel is rendered in HTML.
// Link is the URL of the site, when the channel is rendered, the image is a link to the site.
func NewImage(url, title, link string) *Image {
	return &Image{
		URL:   url,
		Title: title,
		Link:  link,
	}
}

// SetWidth specifies the width of the image in pixels.
func (i *Image) SetWidth(width int) *Image {
	i.Width = width
	return i
}

// SetHeight specifies the height of the image in pixels.
func (i *Image) SetHeight(height int) *Image {
	i.Height = height
	return i
}

// SetDescription specifies the text that is included in the TITLE attribute of the link formed around the image in the HTML rendering.
func (i *Image) SetDescription(description string) *Image {
	i.Description = description
	return i
}

type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type Enclosure struct {
	URL      string `xml:"url,attr"`
	Length   int    `xml:"length,attr"`
	MimeType string `xml:"type,attr"`
}

type Category struct {
	Domain string `xml:"domain,attr"`
	Value  string `xml:",chardata"`
}

type Item struct {
	// All elements of an item are optional,
	// however at least one of title or
	// description must be present.
	Title       string `xml:"title,omitempty"`
	Description string `xml:"description,omitempty"`

	// Optional elements
	Link       string     `xml:"link,omitempty"`
	Author     string     `xml:"author,omitempty"`
	Comments   string     `xml:"comments,omitempty"`
	PubDate    string     `xml:"pubDate,omitempty"`
	GUID       *GUID      `xml:"guid,omitempty"`
	Enclosure  *Enclosure `xml:"enclosure,omitempty"`
	Categories []Category `xml:"category,omitempty"`
}

// NewItemWithTitle creates an item which represents a "story" -- much like a story in a newspaper or magazine; if so its description is a synopsis of the story, and the link points to the full story. An item may also be complete in itself, if so, the description contains the text (entity-encoded HTML is allowed), and the link and title may be omitted. All elements of an item are optional, however at least one of title or description must be present.
func NewItemWithTitle(title string) *Item {
	return &Item{
		Title: title,
	}
}

// NewItemWithDescription creates an item which represents a "story" -- much like a story in a newspaper or magazine; if so its description is a synopsis of the story, and the link points to the full story. An item may also be complete in itself, if so, the description contains the text (entity-encoded HTML is allowed), and the link and title may be omitted. All elements of an item are optional, however at least one of title or description must be present.
func NewItemWithDescription(description string) *Item {
	return &Item{
		Description: html.EscapeString(description),
	}
}

// SetTitle sets the title of the item.
func (i *Item) SetTitle(title string) *Item {
	i.Title = title
	return i
}

// SetDescription sets the item synopsis.
func (i *Item) SetDescription(description string) *Item {
	i.Description = html.EscapeString(description)
	return i
}

// SetLink sets the URL of the item.
func (i *Item) SetLink(link string) *Item {
	i.Link = link
	return i
}

// SetAuthor sets the email address of the author of the item.
func (i *Item) SetAuthor(email, name string) *Item {
	author := email
	if name != "" {
		author += " (" + name + ")"
	}
	i.Author = author
	return i
}

// SetComments sets the URL of a page for comments relating to the item.
func (i *Item) SetComments(comments string) *Item {
	i.Comments = comments
	return i
}

// SetGUID sets a string that uniquely identifies the item.
func (i *Item) SetGUID(value string, isPermaLink bool) *Item {
	i.GUID = &GUID{
		Value:       value,
		IsPermaLink: isPermaLink,
	}
	return i
}

// SetEnclosure describes a media object that is attached to the item.
// It has three required attributes. url says where the enclosure is located, length says how big it is in bytes, and type says what its type is, a standard MIME type.
//
// Example:
//
// <enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" type="audio/mpeg" />
func (i *Item) SetEnclosure(url string, length int, mimeType string) *Item {
	i.Enclosure = &Enclosure{
		URL:      url,
		Length:   length,
		MimeType: mimeType,
	}
	return i
}

// SetPubDate indicates when the item was published.
// If it's a date in the future, aggregators may choose to not display the item until that date.
func (i *Item) SetPubDate(pubDate time.Time) *Item {
	i.PubDate = pubDate.Format(time.RFC822)
	return i
}

// AddCategory includes the item in one or more categories.
// More than one category may be specified by including multiple category elements.
//
// It has one optional attribute, domain, a string that identifies a categorization taxonomy.
//
// The value of the element is a forward-slash-separated string that identifies a hierarchic location in the indicated taxonomy. Processors may establish conventions for the interpretation of categories.
func (i *Item) AddCategory(value, domain string) *Item {
	i.Categories = append(i.Categories, Category{
		Domain: domain,
		Value:  value,
	})
	return i
}

type Channel struct {
	// Required elements
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	// Optional elements
	Language       string     `xml:"language,omitempty"`
	Copyright      string     `xml:"copyright,omitempty"`
	ManagingEditor string     `xml:"managingEditor,omitempty"`
	WebMaster      string     `xml:"webMaster,omitempty"`
	PubDate        string     `xml:"pubDate,omitempty"`
	LastBuildDate  string     `xml:"lastBuildDate,omitempty"`
	Generator      string     `xml:"generator,omitempty"`
	Docs           string     `xml:"docs,omitempty"`
	Categories     []Category `xml:"category,omitempty"`
	TTL            int        `xml:"ttl,omitempty"`
	Image          *Image     `xml:"image,omitempty"`
	Items          []*Item    `xml:"item,omitempty"`
}

// NewChannel creates a new channel.
//
// Title represents the name of the channel. It's how people refer to your service. If you have an HTML website that contains the same information as your RSS file, the title of your channel should be the same as the title of your website.
//
// Link is the URL to the HTML website corresponding to the channel.
//
// Description is the phrase or sentence describing the channel.
func NewChannel(title, link, description string) *Channel {
	return &Channel{
		Title:       title,
		Link:        link,
		Description: description,
	}
}

// SetLanguage sets the language the channel is written in.
//
// Examples: en-us, it-it, en-gb, en-ca, fr-fr, de-de, ja-jp, zh-cn, zh-tw, etc.
func (c *Channel) SetLanguage(language string) *Channel {
	c.Language = language
	return c
}

// SetCopyright sets the copyright notice for content in the channel.
func (c *Channel) SetCopyright(copyright string) *Channel {
	c.Copyright = copyright
	return c
}

// SetManagingEditor sets the email address for person responsible for editorial content.
func (c *Channel) SetManagingEditor(email, name string) *Channel {
	managingEditor := email
	if name != "" {
		managingEditor += " (" + name + ")"
	}
	c.ManagingEditor = managingEditor
	return c
}

// SetWebMaster sets the email address for person responsible for technical issues relating to channel.
func (c *Channel) SetWebMaster(email, name string) *Channel {
	webMaster := email
	if name != "" {
		webMaster += " (" + name + ")"
	}
	c.WebMaster = webMaster
	return c
}

// SetPubDate indicates when the channel was last updated.
func (c *Channel) SetPubDate(pubDate time.Time) *Channel {
	c.PubDate = pubDate.Format(time.RFC822)
	return c
}

// SetLastBuildDate indicates the last time the content of the channel changed.
func (c *Channel) SetLastBuildDate(lastBuildDate time.Time) *Channel {
	c.LastBuildDate = lastBuildDate.Format(time.RFC822)
	return c
}

// SetGenerator indicates the program used to generate the channel.
func (c *Channel) SetGenerator(generator string) *Channel {
	c.Generator = generator
	return c
}

// SetDocs indicates the documentation of the format used in the RSS file.
func (c *Channel) SetDocs(docs string) *Channel {
	c.Docs = docs
	return c
}

// AddCategory specifies one or more categories that the channel belongs to.
// It follows the same rules as the <item>-level category element.
func (c *Channel) AddCategory(value, domain string) *Channel {
	c.Categories = append(c.Categories, Category{
		Domain: domain,
		Value:  value,
	})
	return c
}

// SetTTL specifies the number of minutes that indicates how long a channel can be cached before refreshing from the source.
func (c *Channel) SetTTL(ttl int) *Channel {
	c.TTL = ttl
	return c
}

// SetImage specifies a GIF, JPEG or PNG image that can be displayed with the channel.
func (c *Channel) SetImage(image *Image) *Channel {
	c.Image = image
	return c
}

// AddItem adds an item to the channel.
func (c *Channel) AddItem(item *Item) *Channel {
	c.Items = append(c.Items, item)
	return c
}

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
}

// NewFeed creates a new RSS feed.
func NewFeed(channel *Channel) *Feed {
	return &Feed{
		Version: "2.0",
		Channel: channel,
	}
}

func (f *Feed) toXML(indent bool) ([]byte, error) {
	if indent {
		// nolint: wrapcheck
		return xml.MarshalIndent(f, "", "  ")
	}
	// nolint: wrapcheck
	return xml.Marshal(f)
}

// ToXML returns the XML representation of the feed.
func (f *Feed) ToXML(indent, includeHeader bool) ([]byte, error) {
	bytes, err := f.toXML(indent)
	if err != nil {
		return nil, fmt.Errorf("error XML encoding RSS feed: %w", err)
	}
	if includeHeader {
		return append([]byte(xml.Header), bytes...), nil
	}
	return bytes, nil
}
