// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"net"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config is the application configuration.
type Config struct {
	Server  Server
	SMTP    SMTP
	Content Content
}

// Content holds site content configuration.
type Content struct {
	ShowDrafts bool `env:"SHOW_DRAFTS"`
}

// SMTP holds outbound mail configuration.
type SMTP struct {
	Host string `env:"SMTP_HOST"`
	Port int    `env:"SMTP_PORT" envDefault:"587"`
	User string `env:"SMTP_USER"`
	Pass string `env:"SMTP_PASS"`
	From string `env:"SMTP_FROM"`
	To   string `env:"CONTACT_TO"`
}

// Enabled reports whether enough data is present to actually send mail.
func (s SMTP) Enabled() bool {
	return s.Host != "" && s.User != "" && s.Pass != "" && s.To != ""
}

// Sender is the envelope from address, defaulting to the auth user.
func (s SMTP) Sender() string {
	if s.From != "" {
		return s.From
	}
	return s.User
}

// Server holds HTTP server configuration.
type Server struct {
	Host            string        `env:"HOST"             envDefault:"127.0.0.1"`
	Port            int           `env:"PORT"             envDefault:"8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"     envDefault:"5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"    envDefault:"10s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT"     envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

// Addr builds the listen address from host and port.
func (s Server) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

// Load parses Config from environment variables.
func Load() (Config, error) {
	return env.ParseAs[Config]()
}
