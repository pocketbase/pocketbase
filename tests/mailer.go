package tests

import (
	"slices"
	"sync"

	"github.com/pocketbase/pocketbase/tools/mailer"
)

var _ mailer.Mailer = (*TestMailer)(nil)

// TestMailer is a mock [mailer.Mailer] implementation.
type TestMailer struct {
	mux      sync.Mutex
	messages []*mailer.Message
}

// Send implements [mailer.Mailer] interface.
func (tm *TestMailer) Send(m *mailer.Message) error {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	tm.messages = append(tm.messages, m)
	return nil
}

// Reset clears any previously test collected data.
func (tm *TestMailer) Reset() {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	tm.messages = nil
}

// TotalSend returns the total number of sent messages.
func (tm *TestMailer) TotalSend() int {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	return len(tm.messages)
}

// Messages returns a shallow copy of all of the collected test messages.
func (tm *TestMailer) Messages() []*mailer.Message {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	return slices.Clone(tm.messages)
}

// FirstMessage returns a shallow copy of the first sent message.
//
// Returns an empty mailer.Message struct if there are no sent messages.
func (tm *TestMailer) FirstMessage() mailer.Message {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	var m mailer.Message

	if len(tm.messages) > 0 {
		return *tm.messages[0]
	}

	return m
}

// LastMessage returns a shallow copy of the last sent message.
//
// Returns an empty mailer.Message struct if there are no sent messages.
func (tm *TestMailer) LastMessage() mailer.Message {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	var m mailer.Message

	if len(tm.messages) > 0 {
		return *tm.messages[len(tm.messages)-1]
	}

	return m
}
