package rss

import (
	"testing"

	"github.com/dusted-go/utils/assert"
)

func Test_ToXML_WithEmptyFeed(t *testing.T) {
	feed := NewFeed(nil)

	bytes, err := feed.ToXML(false, false)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<rss version="2.0"></rss>`

	assert.Equal(t, expected, string(bytes))
}

func Test_ToXML_WithEmptyFeedAndHeader(t *testing.T) {
	feed := NewFeed(nil)

	bytes, err := feed.ToXML(false, true)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"></rss>`

	assert.Equal(t, expected, string(bytes))
}

func Test_ToXML_WithEmptyChannel(t *testing.T) {
	feed := NewFeed(NewChannel("", "", ""))

	bytes, err := feed.ToXML(false, false)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<rss version="2.0"><channel><title></title><link></link><description></description></channel></rss>`

	assert.Equal(t, expected, string(bytes))
}

func Test_ToXML_WithTwoItemsAndCategories(t *testing.T) {
	feed := NewFeed(
		NewChannel("title", "link", "description").
			SetLanguage("en-us").
			SetManagingEditor("foo@bar.com", "Foo Bar").
			AddCategory("red", "https://example.org/red").
			AddCategory("blue", "https://example.org/blue").
			AddItem(
				NewItemWithTitle("article A").
					SetGUID("https://example.org/article-a", true).
					SetDescription("bla blah")).AddItem(
			NewItemWithTitle("article B").
				SetGUID("https://example.org/article-b", true).
				SetDescription("yada yada")))

	bytes, err := feed.ToXML(true, true)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>title</title>
    <link>link</link>
    <description>description</description>
    <language>en-us</language>
    <managingEditor>foo@bar.com (Foo Bar)</managingEditor>
    <category domain="https://example.org/red">red</category>
    <category domain="https://example.org/blue">blue</category>
    <item>
      <title>article A</title>
      <description>bla blah</description>
      <guid isPermaLink="true">https://example.org/article-a</guid>
    </item>
    <item>
      <title>article B</title>
      <description>yada yada</description>
      <guid isPermaLink="true">https://example.org/article-b</guid>
    </item>
  </channel>
</rss>`

	assert.Equal(t, expected, string(bytes))
}
