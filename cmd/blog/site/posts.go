package site

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1" // nolint: gosec // used for cache invalidation
	"encoding/hex"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/dusted-go/diagnostic/v3/dlog"
	"github.com/dusted-go/fault/fault"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	DefaultBlogPostPath = "dist/posts"
)

type BlogPost struct {
	ID          string
	Title       string
	PublishDate time.Time
	Tags        []string
	HashCode    string
	markdown    string
	html        string
}

func (b *BlogPost) Validate() error {
	return nil
}

func (b *BlogPost) Year() int {
	return b.PublishDate.Year()
}

func (b *BlogPost) Excerpt() string {
	return "ToDo"
}

func (b *BlogPost) HTML() (template.HTML, error) {
	if len(b.html) > 0 {
		// nolint: gosec // This is safe content
		return template.HTML(b.html), nil
	}

	parser := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		), goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		))

	var buf bytes.Buffer
	if err := parser.Convert([]byte(b.markdown), &buf); err != nil {
		return template.HTML(""),
			fault.SystemWrap(err, "error converting Markdown into HTML")
	}

	// nolint: gosec // string was already escaped before
	return template.HTML(buf.Bytes()), nil
}

func parseBlogPost(
	blogPostID string,
	publishDate time.Time,
	buffer []byte,
) (
	*BlogPost,
	error,
) {
	bufferWithoutBOM := bytes.TrimLeft(buffer, "\xef\xbb\xbf")

	var metadata []string
	var title string
	body := strings.Builder{}

	readMeta := false
	readBody := false

	scanner := bufio.NewScanner(bytes.NewReader(bufferWithoutBOM))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "<!--") {
			readMeta = true
			continue
		} else if strings.HasPrefix(line, "-->") {
			readMeta = false
			readBody = true
			continue
		}

		if readMeta {
			metadata = append(metadata, line)
		} else if strings.TrimSpace(line) == "" && body.Len() == 0 {
			continue
		} else if strings.HasPrefix(line, "# ") && body.Len() == 0 {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
		} else if readBody {
			body.WriteString(line)
			body.WriteString("\n")
		}
	}

	if len(title) == 0 {
		return nil, fault.System("blog post title is missing")
	}
	content := body.String()

	isHTML := false
	var tags []string
	for _, meta := range metadata {
		metaParts := strings.SplitN(meta, ":", 2)
		key := strings.ToLower(strings.TrimSpace(metaParts[0]))
		if key == "tags" {
			tags = strings.Split(strings.TrimSpace(metaParts[1]), " ")
		} else if key == "type" {
			isHTML = strings.ToLower(strings.TrimSpace(metaParts[1])) == "html"
		} else {
			return nil, fault.Systemf("unknown blog post metadata key: %s", key)
		}
	}

	valueToHash := title + content + publishDate.String()
	for _, tag := range tags {
		valueToHash = valueToHash + tag
	}

	// nolint: gosec // hash used for caching, not security
	hash := sha1.New()
	hash.Write([]byte(valueToHash))
	hashCode := hex.EncodeToString(hash.Sum(nil))

	blogPost := &BlogPost{
		ID:          blogPostID,
		Title:       title,
		PublishDate: publishDate,
		Tags:        tags,
		HashCode:    hashCode,
	}
	if isHTML {
		blogPost.html = content
	} else {
		blogPost.markdown = content
	}

	return blogPost, nil
}

func ReadBlogPost(ctx context.Context, path string, blogPostID string) (*BlogPost, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fault.SystemWrapf(err, "error reading files from directory '%s'", path)
	}

	fileName := ""
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, blogPostID+".md") {
			fileName = name
			break
		}
	}

	if len(fileName) == 0 {
		return nil, fault.Systemf("blog post with ID '%s' not found", blogPostID)
	}

	fileNameParts := strings.SplitN(fileName, "-", 2)
	publishDate, err := time.Parse("2006_01_02", fileNameParts[0])
	if err != nil {
		return nil, fault.SystemWrapf(err, "error parsing date from file name '%s'", fileName)
	}

	fileBuffer, err := os.ReadFile(path + "/" + fileName)
	if err != nil {
		return nil, fault.SystemWrap(err, "error reading blog post file")
	}

	blogPost, err := parseBlogPost(blogPostID, publishDate, fileBuffer)
	if err != nil {
		return nil, fault.SystemWrapf(err, "error parsing blog post '%s'", fileName)
	}

	return blogPost, nil
}

func ReadBlogPosts(ctx context.Context, path string) ([]*BlogPost, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fault.SystemWrapf(err, "error reading files from directory '%s'", path)
	}
	blogPosts := []*BlogPost{}

	for _, f := range files {
		fileName := f.Name()
		if f.IsDir() {
			dlog.New(ctx).Warning().Fmt("Cannot read blog post because it is a directory: %s", fileName)
			continue
		}
		if !strings.HasSuffix(fileName, ".md") {
			dlog.New(ctx).Warning().Fmt("Skipping file %s because it doesn't appear to be a Markdown file.", fileName)
			continue
		}
		fileNameParts := strings.SplitN(fileName, "-", 2)
		blogPostID := strings.TrimSuffix(fileNameParts[1], ".md")
		publishDate, err := time.Parse("2006_01_02", fileNameParts[0])
		if err != nil {
			dlog.New(ctx).Error().Err(err).Msg("An error occurred when parsing the blog post's date.")
			continue
		}

		fileBuffer, err := os.ReadFile(path + "/" + fileName)
		if err != nil {
			dlog.New(ctx).Error().Err(err).Msg("An error occurred when reading a .md file from disk.")
			continue
		}

		blogPost, err := parseBlogPost(blogPostID, publishDate, fileBuffer)
		if err != nil {
			dlog.New(ctx).Warning().Fmt("Skipping blog post '%s': %s", fileName, err.Error())
			continue
		}

		blogPosts = append(blogPosts, blogPost)
		dlog.New(ctx).Info().Fmt("Successfully parsed: %s", fileName)
	}

	dlog.New(ctx).Info().Fmt("Total blog posts parsed: %d", len(blogPosts))
	return blogPosts, nil
}
