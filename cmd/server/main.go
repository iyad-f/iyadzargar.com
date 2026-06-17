// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log/slog"
	"os"

	"github.com/iyad-f/iyadf.com/internal/config"
	"github.com/iyad-f/iyadf.com/internal/content"
	"github.com/iyad-f/iyadf.com/internal/mailer"
	"github.com/iyad-f/iyadf.com/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "err", err)
		os.Exit(1)
	}

	site, err := content.Load()
	if err != nil {
		logger.Error("load site", "err", err)
		os.Exit(1)
	}

	sender := mailer.NewSMTP(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.Sender(), cfg.SMTP.To)
	if !cfg.SMTP.Enabled() {
		logger.Warn("smtp not configured, mailer wont work")
	}

	if err := server.New(cfg.Server, logger, site, sender).Run(); err != nil {
		logger.Error("server startup failed", "err", err)
		os.Exit(1)
	}
}
