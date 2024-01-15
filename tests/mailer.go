package tests

import (
	"github.com/pocketbase/pocketbase/tools/mailer"
)

var _ mailer.Mailer = (*TestMailer)(nil)

// TestMailer is a mock `mailer.Mailer` implementation.
type TestMailer struct {
	TotalSend   int
	LastMessage mailer.Message
}

// Reset clears any previously test collected data.
func (t *TestMailer) Reset() {
	t.TotalSend = 0
	t.LastMessage = mailer.Message{}
}

// Send implements `mailer.Mailer` interface.
func (t *TestMailer) Send(m *mailer.Message) error {
	t.TotalSend++
	t.LastMessage = *m

	return nil
}
