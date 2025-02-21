package subscriptions

import (
	"io"
)

// Message defines a client's channel data.
type Message struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

// WriteSSE writes the current message in a SSE format into the provided writer.
//
// For example, writing to a router.Event:
//
//	m := Message{Name: "users/create", Data: []byte{...}}
//	m.Write(e.Response, "yourEventId")
//	e.Flush()
func (m *Message) WriteSSE(w io.Writer, eventId string) error {
	parts := [][]byte{
		[]byte("id:" + eventId + "\n"),
		[]byte("event:" + m.Name + "\n"),
		[]byte("data:"),
		m.Data,
		[]byte("\n\n"),
	}

	for _, part := range parts {
		_, err := w.Write(part)
		if err != nil {
			return err
		}
	}

	return nil
}
