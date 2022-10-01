package site

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/dusted-go/diagnostic/v3/dlog"
	"github.com/dusted-go/fault/fault"
)

type BlogPost struct {
	ID          string
	Title       string
	PublishDate time.Time
	Tags        []string
	Markdown    string
	HashCode    string
}

func (b *BlogPost) Validate() error {
	return nil
}

func (b *BlogPost) Year() int {
	return b.PublishDate.Year()
}

func (b *BlogPost) URLEncodedTitle() string {
	return "ToDo"
}

func (b *BlogPost) Excerpt() string {
	return "ToDo"
}

func ReadBlogPosts(ctx context.Context, path string) ([]*BlogPost, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fault.SystemWrap(err, "site", "ReadBlogPosts", "error reading blogs posts")
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
		publishDate, err := time.Parse("2006_01_02", fileNameParts[0])
		if err != nil {
			dlog.New(ctx).Error().Err(err).Msg("An error occurred when parsing the blog post's date.")
			continue
		}

		file, err := os.Open(path + "/" + fileName)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		if err != nil {
			dlog.New(ctx).Error().Err(err).Msg("An error occurred when reading blog post.")
			continue
		}

		var metadata []string
		var body []string
		readMeta := false
		readBody := false

		scanner := bufio.NewScanner(file)
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
			} else if readBody {
				body = append(body, line)
			}
		}

		var tags []string
		for _, meta := range metadata {
			metaParts := strings.SplitN(meta, ":", 2)
			key := strings.ToLower(strings.TrimSpace(metaParts[0]))
			if key != "tags" {
				dlog.New(ctx).Warning().Fmt("Skipping unknown meta data key '%s' for blog post %s.", key, fileName)
				continue
			}
			tags = strings.Split(strings.TrimSpace(metaParts[1]), " ")
		}

		var title string
		content := strings.Builder{}
		skipBlogPost := false
		for _, line := range body {
			if content.Len() == 0 {
				if strings.TrimSpace(line) == "" {
					continue
				}
				if strings.HasPrefix(line, "# ") {
					title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
				}
			}

			if len(title) == 0 {
				dlog.New(ctx).Warning().Fmt("Skipping blog post without title: %s", fileName)
				skipBlogPost = true
				break
			}

			content.WriteString(line)
			content.WriteString("\n")
		}
		if skipBlogPost {
			continue
		}

		markdown := content.String()

		valueToHash := title + markdown + publishDate.String()
		for _, tag := range tags {
			valueToHash = valueToHash + tag
		}

		// nolint: gosec // hash used for caching, not security
		hash := sha1.New()
		hash.Write([]byte(valueToHash))
		hashCode := hex.EncodeToString(hash.Sum(nil))

		blogPost := &BlogPost{
			ID:          fileNameParts[1],
			Title:       title,
			PublishDate: publishDate,
			Tags:        tags,
			Markdown:    markdown,
			HashCode:    hashCode,
		}

		blogPosts = append(blogPosts, blogPost)
		dlog.New(ctx).Info().Fmt("Successfully parsed: %s", fileName)
	}

	dlog.New(ctx).Info().Fmt("Total blog posts parsed: %d", len(blogPosts))
	return blogPosts, nil
}
