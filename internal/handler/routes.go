// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/iyad-f/iyadzargar.com/web/views"
)

// Register mounts the site's routes on mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.Handle("GET /{$}", templ.Handler(views.Home(h.site)))
	mux.Handle("GET /projects", templ.Handler(views.Projects(h.site)))
	mux.Handle("GET /about", templ.Handler(views.About(h.site)))
	mux.HandleFunc("GET /contact", h.contactPage)
	mux.Handle("POST /contact", h.contactRateLimit(http.HandlerFunc(h.contactSubmit)))
}
