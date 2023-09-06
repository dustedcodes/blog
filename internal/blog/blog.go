package blog

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1" // nolint: gosec // used for cache invalidation
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/dusted-go/logging/stackdriver"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	syntax "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	DefaultBlogPostPath = "dist/posts"
)

type Post struct {
	ID          string
	Title       string
	PublishDate time.Time
	Tags        []string
	HashCode    string
	markdown    string
	html        string
}

func (p *Post) Year() int {
	return p.PublishDate.Year()
}

func (p *Post) HTML() (template.HTML, error) {
	if len(p.html) > 0 {
		// nolint: gosec // This is safe content
		return template.HTML(p.html), nil
	}

	parser := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			syntax.NewHighlighting(
				syntax.WithStyle("gruvbox"),
				syntax.WithFormatOptions(
					chromahtml.TabWidth(4),
					chromahtml.WithLineNumbers(false),
					chromahtml.PreventSurroundingPre(false),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		), goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		))

	var buf bytes.Buffer
	if err := parser.Convert([]byte(p.markdown), &buf); err != nil {
		return template.HTML(""),
			fmt.Errorf("error converting Markdown into HTML: %w", err)
	}

	// nolint: gosec // string was already escaped before
	return template.HTML(buf.Bytes()), nil
}

func parsePost(
	blogPostID string,
	publishDate time.Time,
	buffer []byte,
) (
	*Post,
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
		return nil, errors.New("blog post title is missing")
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
			return nil, fmt.Errorf("unknown blog post metadata key: %s", key)
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

	blogPost := &Post{
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

func ReadPost(path string, blogPostID string) (*Post, error) {

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading files from directory '%s': %w", path, err)
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
		return nil, fmt.Errorf("blog post with ID '%s' not found", blogPostID)
	}

	fileNameParts := strings.SplitN(fileName, "-", 2)
	publishDate, err := time.Parse("2006_01_02", fileNameParts[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing date from file '%s': %w", fileName, err)
	}

	fileBuffer, err := os.ReadFile(path + "/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading blog post file: %w", err)
	}

	blogPost, err := parsePost(blogPostID, publishDate, fileBuffer)
	if err != nil {
		return nil, fmt.Errorf("error parsing blog post '%s': %w", fileName, err)
	}

	return blogPost, nil
}

func ReadPosts(ctx context.Context, path string) ([]*Post, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading files from directory '%s': %w", path, err)
	}
	blogPosts := []*Post{}

	logger := stackdriver.GetLogger(ctx)

	for _, f := range files {
		fileName := f.Name()
		if f.IsDir() {
			logger.Warn("Cannot read blog post because it is a directory.",
				"filename", fileName)
			continue
		}
		if !strings.HasSuffix(fileName, ".md") {
			logger.Warn("Skipping file because it doesn't appear to be a Markdown file.",
				"filename", fileName)
			continue
		}
		fileNameParts := strings.SplitN(fileName, "-", 2)
		blogPostID := strings.TrimSuffix(fileNameParts[1], ".md")
		publishDate, err := time.Parse("2006_01_02", fileNameParts[0])
		if err != nil {
			logger.Error("An error occurred when parsing the blog post's date.", "error", err)
			continue
		}

		fileBuffer, err := os.ReadFile(path + "/" + fileName)
		if err != nil {
			logger.Error("An error occurred when reading a .md file from disk.", "error", err)
			continue
		}

		blogPost, err := parsePost(blogPostID, publishDate, fileBuffer)
		if err != nil {
			logger.Error("Skipping blog post because of parsing error.",
				"filename", fileName,
				"error", err)
			continue
		}

		blogPosts = append(blogPosts, blogPost)
		logger.Debug("Successfully parsed blog post.",
			"filename", fileName)
	}

	logger.Info("Finished parsing blog posts.",
		"count", len(blogPosts))
	return blogPosts, nil
}
