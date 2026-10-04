// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package content

import (
	_ "embed"
	"fmt"
	"slices"

	"github.com/pelletier/go-toml/v2"
)

//go:embed site.toml
var raw []byte

// Language is a project's primary implementation language.
type Language string

// Recognized languages.
const (
	Go         Language = "Go"
	Rust       Language = "Rust"
	Zig        Language = "Zig"
	Python     Language = "Python"
	JavaScript Language = "JavaScript"
)

// Project is one project shown on the site.
type Project struct {
	Name      string   `toml:"name"`
	Desc      string   `toml:"desc"`
	Lang      Language `toml:"lang"`
	Tags      []string `toml:"tags"`
	URL       string   `toml:"url"`
	LinkLabel string   `toml:"link_label"`
	Featured  bool     `toml:"featured"`
}

// Owner is the site owner's identity and contact details.
type Owner struct {
	Name           string `toml:"name"`
	URL            string `toml:"url"`
	Email          string `toml:"email"`
	GitHubUsername string `toml:"github_username"`
	Location       string `toml:"location"`
}

// MailtoURL is the mailto link for the owner's email.
func (o Owner) MailtoURL() string { return "mailto:" + o.Email }

// GitHubURL is the full URL to the owner's GitHub profile.
func (o Owner) GitHubURL() string { return "https://github.com/" + o.GitHubUsername }

// GitHubLabel is the display form of the owner's GitHub profile.
func (o Owner) GitHubLabel() string { return "github.com/" + o.GitHubUsername }

// SkillGroup is a category of related skills.
type SkillGroup struct {
	Category string   `toml:"category"`
	Items    []string `toml:"items"`
}

// About is the owner's background.
type About struct {
	Bio         []string     `toml:"bio"`
	SkillGroups []SkillGroup `toml:"skills"`
}

// Site is the site's content.
type Site struct {
	Owner    Owner     `toml:"owner"`
	About    About     `toml:"about"`
	Projects []Project `toml:"projects"`
	Posts    []Post    `toml:"-"`
}

// Featured returns the projects marked as featured.
func (s Site) Featured() []Project {
	var out []Project
	for _, p := range s.Projects {
		if p.Featured {
			out = append(out, p)
		}
	}
	return out
}

// Post returns the post with the given slug.
func (s Site) Post(slug string) (Post, bool) {
	i := slices.IndexFunc(s.Posts, func(p Post) bool { return p.Slug == slug })
	if i < 0 {
		return Post{}, false
	}
	return s.Posts[i], true
}

// Load parses the embedded site content and posts. Draft posts are kept only
// when drafts is set.
func Load(drafts bool) (Site, error) {
	var s Site
	if err := toml.Unmarshal(raw, &s); err != nil {
		return Site{}, fmt.Errorf("parse site.toml: %w", err)
	}
	posts, err := loadPosts(drafts)
	if err != nil {
		return Site{}, fmt.Errorf("load posts: %w", err)
	}
	s.Posts = posts
	return s, nil
}
