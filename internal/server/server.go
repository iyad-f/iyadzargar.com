// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/iyad-f/iyadzargar.com/internal/config"
	"github.com/iyad-f/iyadzargar.com/internal/content"
	"github.com/iyad-f/iyadzargar.com/internal/handler"
	"github.com/iyad-f/iyadzargar.com/internal/mailer"
	"github.com/iyad-f/iyadzargar.com/web"
)

// Server is the application's HTTP server.
type Server struct {
	cfg    config.Server
	logger *slog.Logger
	site   content.Site
	sender mailer.Sender
}

// New returns a Server.
func New(cfg config.Server, logger *slog.Logger, site content.Site, sender mailer.Sender) *Server {
	return &Server{cfg: cfg, logger: logger, site: site, sender: sender}
}

func (s *Server) addRoutes(mux *http.ServeMux) {
	handler.NewHandler(s.logger, s.site, s.sender).Register(mux)

	// The embed path is fixed, so fs.Sub cannot fail.
	static, _ := fs.Sub(web.Static, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
}

// Run serves until SIGINT/SIGTERM triggers a graceful shutdown.
func (s *Server) Run() error {
	mux := http.NewServeMux()
	s.addRoutes(mux)

	srv := &http.Server{
		Addr:         s.cfg.Addr(),
		Handler:      mux,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	}

	s.logger.Info("starting server", "addr", srv.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-quit:
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("server is now ready to exit, bye bye...")
	return <-errCh
}
