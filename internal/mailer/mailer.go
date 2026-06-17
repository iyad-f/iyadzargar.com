// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

package mailer

import (
	"context"
	"fmt"

	"github.com/wneessen/go-mail"
)

// Message is an email to send.
type Message struct {
	ReplyTo string
	Subject string
	Body    string
}

// Sender delivers a Message.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// SMTP sends mail through an authenticated server over STARTTLS.
type SMTP struct {
	client *mail.Client
	err    error
	from   string
	to     string
}

// NewSMTP returns a sender for the given server and account. The client is
// built once and reused. If it cannot be created, the error is returned from
// Send rather than here.
func NewSMTP(host string, port int, user, pass, from, to string) *SMTP {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithTLSPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(user),
		mail.WithPassword(pass),
	)
	return &SMTP{client: client, err: err, from: from, to: to}
}

// Send delivers msg, honoring ctx for the dial and the exchange. The reused
// client opens and closes a fresh connection per call.
func (s *SMTP) Send(ctx context.Context, msg Message) error {
	if s.err != nil {
		return fmt.Errorf("smtp client: %w", s.err)
	}

	m := mail.NewMsg()
	if err := m.From(s.from); err != nil {
		return fmt.Errorf("set from: %w", err)
	}
	if err := m.To(s.to); err != nil {
		return fmt.Errorf("set to: %w", err)
	}
	if msg.ReplyTo != "" {
		if err := m.ReplyTo(msg.ReplyTo); err != nil {
			return fmt.Errorf("set reply-to: %w", err)
		}
	}
	m.Subject(msg.Subject)
	m.SetBodyString(mail.TypeTextPlain, msg.Body)

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}
	return nil
}
