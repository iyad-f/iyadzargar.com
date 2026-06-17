// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/httprate"

	"github.com/iyad-f/iyadf.com/internal/mailer"
	"github.com/iyad-f/iyadf.com/web/views"
)

// contactForm is a submission to validate before it is sent as mail.
type contactForm struct {
	Name    string `validate:"required,max=100"`
	Email   string `validate:"required,email,max=254"`
	Subject string `validate:"required,max=150"`
	Message string `validate:"required,max=5000"`
}

// contactRateLimit throttles submissions per client, answering an over limit
// request with the usual toast or redirect.
func (h *Handler) contactRateLimit(next http.Handler) http.Handler {
	return httprate.Limit(5, time.Hour,
		httprate.WithKeyFuncs(httprate.KeyByRealIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			h.contactRespond(w, r, "ratelimited")
		}),
	)(next)
}

// contactPage renders the contact page. A status query from an earlier post
// drives the banner shown on load.
func (h *Handler) contactPage(w http.ResponseWriter, r *http.Request) {
	if err := views.Contact(h.site, r.URL.Query().Get("status")).Render(r.Context(), w); err != nil {
		h.logger.Error("render contact page", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// contactSubmit validates a submission and sends it as mail.
func (h *Handler) contactSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.contactRespond(w, r, "error")
		return
	}

	// Honeypot. Bots fill hidden fields humans never see, so treat a filled
	// one as success without sending anything.
	if strings.TrimSpace(r.PostForm.Get("company")) != "" {
		h.contactRespond(w, r, "sent")
		return
	}

	form := contactForm{
		Name:    strings.TrimSpace(r.PostForm.Get("name")),
		Email:   strings.TrimSpace(r.PostForm.Get("email")),
		Subject: strings.TrimSpace(r.PostForm.Get("subject")),
		Message: strings.TrimSpace(r.PostForm.Get("message")),
	}
	if err := h.validate.Struct(form); err != nil {
		h.contactRespond(w, r, "invalid")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	msg := mailer.Message{
		ReplyTo: form.Email,
		Subject: "[iyadf.com] " + form.Subject,
		Body: fmt.Sprintf(
			"New message from the iyadf.com contact form.\n\nName     %s\nEmail    %s\n\n%s\n",
			form.Name, form.Email, form.Message,
		),
	}

	if err := h.sender.Send(ctx, msg); err != nil {
		h.logger.Error("send contact message", "err", err)
		h.contactRespond(w, r, "error")
		return
	}

	h.logger.Info("contact message sent", "reply_to", form.Email)
	h.contactRespond(w, r, "sent")
}

// contactRespond writes the outcome of a submission. An htmx request gets a
// toast, a plain form post a redirect carrying the status.
func (h *Handler) contactRespond(w http.ResponseWriter, r *http.Request, status string) {
	if r.Header.Get("HX-Request") == "true" {
		_ = views.Toast(status).Render(r.Context(), w)
		return
	}
	http.Redirect(w, r, "/contact?status="+status+"#send", http.StatusSeeOther)
}
