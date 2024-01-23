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
func (tm *TestMailer) Reset() {
	tm.TotalSend = 0
	tm.LastMessage = mailer.Message{}
	tm.SentMessages = nil
}

// Send implements `mailer.Mailer` interface.
func (tm *TestMailer) Send(m *mailer.Message) error {
	tm.TotalSend++
	tm.LastMessage = *m
	tm.SentMessages = append(tm.SentMessages, tm.LastMessage)

	return nil
}
