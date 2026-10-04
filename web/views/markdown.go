// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package views

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"

	"github.com/iyad-f/iyadzargar.com/internal/content"
)

// codeStyle is the Chroma style that colors highlighted code.
const codeStyle = "gruvbox"

// markdown renders post bodies. Raw HTML passes through since posts are trusted.
// Code is highlighted with CSS classes, colored by codeStyles.
var markdown = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
		highlighting.NewHighlighting(
			highlighting.WithStyle(codeStyle),
			highlighting.WithFormatOptions(chromahtml.WithClasses(true)),
			highlighting.WithWrapperRenderer(codeBlockWrapper),
		),
	),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

// codeCSS is the stylesheet Chroma generates for codeStyle.
var codeCSS = func() string {
	var b strings.Builder
	if err := chromahtml.New(chromahtml.WithClasses(true)).WriteCSS(&b, styles.Get(codeStyle)); err != nil {
		panic(fmt.Sprintf("generate %s code css: %v", codeStyle, err))
	}
	return b.String()
}()

// codeStyles inlines the code highlighting stylesheet.
func codeStyles() templ.Component {
	return templ.Raw("<style>" + codeCSS + "</style>")
}

// postBody renders a post's markdown body.
func postBody(p content.Post) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		return markdown.Convert([]byte(p.Body), w)
	})
}

// codeBlockWrapper frames each fenced code block with a header naming its
// language, where app.js adds a copy button. Chroma writes its own pre and code
// tags, so they are only written here for blocks it left unhighlighted.
func codeBlockWrapper(w util.BufWriter, c highlighting.CodeBlockContext, entering bool) {
	if !entering {
		if !c.Highlighted() {
			_, _ = w.WriteString("</code></pre>")
		}
		_, _ = w.WriteString("</div>")
		return
	}

	lang, _ := c.Language()
	_, _ = w.WriteString(`<div class="code-block"><div class="code-head"><span>`)
	_, _ = w.Write(util.EscapeHTML(lang))
	_, _ = w.WriteString("</span></div>")
	if !c.Highlighted() {
		_, _ = w.WriteString("<pre><code>")
	}
}
