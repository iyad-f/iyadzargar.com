// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package content

import (
	"bytes"
	"cmp"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// postFiles holds the posts folder. A pattern that matches nothing fails the
// build, so posts/.gitkeep keeps the folder in git and the embed non-empty
// while there are no posts. Delete it once a real post exists.
//
//go:embed all:posts
var postFiles embed.FS

// Post is one blog post, a markdown file under posts/ that opens with a TOML
// front matter block fenced by +++ lines. Its file name minus .md is its slug.
type Post struct {
	Slug        string    `toml:"-"`
	Title       string    `toml:"title"`
	Date        time.Time `toml:"date"`
	Summary     string    `toml:"summary"`
	Tags        []string  `toml:"tags"`
	Draft       bool      `toml:"draft"`
	Body        string    `toml:"-"` // markdown
	ReadingTime int       `toml:"-"` // minutes, at least 1
}

// wordsPerMinute is the reading speed behind ReadingTime.
const wordsPerMinute = 200

var (
	fence       = []byte("+++\n")
	slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

// loadPosts parses every embedded post, newest first. Drafts are dropped
// unless drafts is set.
func loadPosts(drafts bool) ([]Post, error) {
	names, err := fs.Glob(postFiles, "posts/*.md")
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, name := range names {
		slug := strings.TrimSuffix(path.Base(name), ".md")
		if !slugPattern.MatchString(slug) {
			return nil, fmt.Errorf("%s: name must be lowercase words joined by hyphens", name)
		}

		src, err := postFiles.ReadFile(name)
		if err != nil {
			return nil, err
		}

		p, err := parsePost(slug, src)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		if p.Draft && !drafts {
			continue
		}

		posts = append(posts, p)
	}

	slices.SortFunc(posts, func(a, b Post) int {
		if c := b.Date.Compare(a.Date); c != 0 {
			return c
		}
		return cmp.Compare(a.Slug, b.Slug)
	})
	return posts, nil
}

// parsePost splits src into front matter and markdown body.
func parsePost(slug string, src []byte) (Post, error) {
	src = bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))

	rest, ok := bytes.CutPrefix(src, fence)
	if !ok {
		return Post{}, errors.New("missing opening +++ fence")
	}
	meta, body, ok := bytes.Cut(rest, append([]byte("\n"), fence...))
	if !ok {
		return Post{}, errors.New("missing closing +++ fence")
	}

	p := Post{Slug: slug}
	if err := toml.NewDecoder(bytes.NewReader(meta)).DisallowUnknownFields().Decode(&p); err != nil {
		return Post{}, fmt.Errorf("parse front matter: %w", err)
	}
	switch {
	case p.Title == "":
		return Post{}, errors.New("front matter needs a title")
	case p.Date.IsZero():
		return Post{}, errors.New("front matter needs a date")
	case p.Summary == "":
		return Post{}, errors.New("front matter needs a summary")
	}

	p.Body = string(body)

	words := len(strings.Fields(p.Body))
	p.ReadingTime = max(1, int(math.Ceil(float64(words)/wordsPerMinute)))
	return p, nil
}
