// Package cli implements a function to display prompt to user in a terminal
//
// Example:
//
//	c := cli.New()
//	confirm := c.Confirm("Do you want to proceed with the update?", true)

package cli

import (
	"io"
	"os"
)

type Cli struct {
	writer io.Writer
	reader io.Reader
}

func New() *Cli {
	return &Cli{
		writer: os.Stderr,
		reader: os.Stdin,
	}
}

func NewWithIO(writer io.Writer, reader io.Reader) *Cli {
	return &Cli{
		writer: writer,
		reader: reader,
	}
}
