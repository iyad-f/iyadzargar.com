// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"log/slog"

	"github.com/go-playground/validator/v10"

	"github.com/iyad-f/iyadf.com/internal/content"
	"github.com/iyad-f/iyadf.com/internal/mailer"
)

// Handler holds the dependencies shared by the site's routes.
type Handler struct {
	logger   *slog.Logger
	site     content.Site
	sender   mailer.Sender
	validate *validator.Validate
}

// NewHandler returns a Handler.
func NewHandler(logger *slog.Logger, site content.Site, sender mailer.Sender) *Handler {
	return &Handler{
		logger:   logger,
		site:     site,
		sender:   sender,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}
