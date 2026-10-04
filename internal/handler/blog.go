// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"net/http"

	"github.com/iyad-f/iyadzargar.com/web/views"
)

// blogPost renders a single post, or a 404 when no post has the slug.
func (h *Handler) blogPost(w http.ResponseWriter, r *http.Request) {
	p, ok := h.site.Post(r.PathValue("slug"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := views.Post(h.site.Owner, p).Render(r.Context(), w); err != nil {
		h.logger.Error("render post", "slug", p.Slug, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
