package tests

import (
	"github.com/pocketbase/pocketbase/tools/mailer"
)

var _ mailer.Mailer = (*TestMailer)(nil)

// TestMailer is a mock `mailer.Mailer` implementation.
type TestMailer struct {
	TotalSend   int
	LastMessage mailer.Message

	// @todo consider deprecating the above 2 fields?
	SentMessages []mailer.Message
}

// Reset clears any previously test collected data.
func (m *TestMailer) Reset() {
	m.TotalSend = 0
	m.LastMessage = mailer.Message{}
	m.SentMessages = nil
}

// Send implements `mailer.Mailer` interface.
func (c *TestMailer) Send(m *mailer.Message) error {
	c.TotalSend++
	c.LastMessage = *m
	c.SentMessages = append(c.SentMessages, c.LastMessage)

	return nil
}
