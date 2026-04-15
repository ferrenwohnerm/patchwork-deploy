package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event holds a single notification entry.
type Event struct {
	Timestamp   time.Time
	Environment string
	Level       Level
	Message     string
}

// Notifier writes deployment events to an output sink.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier writing to out. Pass nil to use os.Stdout.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Send writes a formatted event to the output sink.
func (n *Notifier) Send(env string, level Level, msg string) Event {
	e := Event{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Level:       level,
		Message:     msg,
	}
	fmt.Fprintf(n.out, "%s [%s] (%s) %s\n",
		e.Timestamp.Format(time.RFC3339),
		strings.ToUpper(string(e.Level)),
		e.Environment,
		e.Message,
	)
	return e
}

// Info is a convenience wrapper for LevelInfo events.
func (n *Notifier) Info(env, msg string) Event {
	return n.Send(env, LevelInfo, msg)
}

// Warn is a convenience wrapper for LevelWarn events.
func (n *Notifier) Warn(env, msg string) Event {
	return n.Send(env, LevelWarn, msg)
}

// Error is a convenience wrapper for LevelError events.
func (n *Notifier) Error(env, msg string) Event {
	return n.Send(env, LevelError, msg)
}
