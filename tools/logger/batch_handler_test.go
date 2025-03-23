package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func TestNewBatchHandlerPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected to panic.")
		}
	}()

	NewBatchHandler(BatchOptions{})
}

func TestNewBatchHandlerDefaults(t *testing.T) {
	h := NewBatchHandler(BatchOptions{
		WriteFunc: func(ctx context.Context, logs []*Log) error {
			return nil
		},
	})

	if h.options.BatchSize != 100 {
		t.Fatalf("Expected default BatchSize %d, got %d", 100, h.options.BatchSize)
	}

	if h.options.Level != slog.LevelInfo {
		t.Fatalf("Expected default Level Info, got %v", h.options.Level)
	}

	if h.options.BeforeAddFunc != nil {
		t.Fatal("Expected default BeforeAddFunc to be nil")
	}

	if h.options.WriteFunc == nil {
		t.Fatal("Expected default WriteFunc to be set")
	}

	if h.group != "" {
		t.Fatalf("Expected empty group, got %s", h.group)
	}

	if len(h.attrs) != 0 {
		t.Fatalf("Expected empty attrs, got %v", h.attrs)
	}

	if len(h.logs) != 0 {
		t.Fatalf("Expected empty logs queue, got %v", h.logs)
	}
}

func TestBatchHandlerEnabled(t *testing.T) {
	h := NewBatchHandler(BatchOptions{
		Level: slog.LevelWarn,
		WriteFunc: func(ctx context.Context, logs []*Log) error {
			return nil
		},
	})

	l := slog.New(h)

	scenarios := []struct {
		level    slog.Level
		expected bool
	}{
		{slog.LevelDebug, false},
		{slog.LevelInfo, false},
		{slog.LevelWarn, true},
		{slog.LevelError, true},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("Level %v", s.level), func(t *testing.T) {
			result := l.Enabled(context.Background(), s.level)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestBatchHandlerSetLevel(t *testing.T) {
	h := NewBatchHandler(BatchOptions{
		Level: slog.LevelWarn,
		WriteFunc: func(ctx context.Context, logs []*Log) error {
			return nil
		},
	})

	if h.options.Level != slog.LevelWarn {
		t.Fatalf("Expected the initial level to be %d, got %d", slog.LevelWarn, h.options.Level)
	}

	h.SetLevel(slog.LevelDebug)

	if h.options.Level != slog.LevelDebug {
		t.Fatalf("Expected the updated level to be %d, got %d", slog.LevelDebug, h.options.Level)
	}
}

func TestBatchHandlerWithAttrsAndWithGroup(t *testing.T) {
	h0 := NewBatchHandler(BatchOptions{
		WriteFunc: func(ctx context.Context, logs []*Log) error {
			return nil
		},
	})

	h1 := h0.WithAttrs([]slog.Attr{slog.Int("test1", 1)}).(*BatchHandler)
	h2 := h1.WithGroup("h2_group").(*BatchHandler)
	h3 := h2.WithAttrs([]slog.Attr{slog.Int("test2", 2)}).(*BatchHandler)

	scenarios := []struct {
		name           string
		handler        *BatchHandler
		expectedParent *BatchHandler
		expectedGroup  string
		expectedAttrs  int
	}{
		{
			"h0",
			h0,
			nil,
			"",
			0,
		},
		{
			"h1",
			h1,
			h0,
			"",
			1,
		},
		{
			"h2",
			h2,
			h1,
			"h2_group",
			0,
		},
		{
			"h3",
			h3,
			h2,
			"",
			1,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if s.handler.group != s.expectedGroup {
				t.Fatalf("Expected group %q, got %q", s.expectedGroup, s.handler.group)
			}

			if s.handler.parent != s.expectedParent {
				t.Fatalf("Expected parent %v, got %v", s.expectedParent, s.handler.parent)
			}

			if totalAttrs := len(s.handler.attrs); totalAttrs != s.expectedAttrs {
				t.Fatalf("Expected %d attrs, got %d", s.expectedAttrs, totalAttrs)
			}
		})
	}
}

func TestBatchHandlerHandle(t *testing.T) {
	ctx := context.Background()

	beforeLogs := []*Log{}
	writeLogs := []*Log{}

	h := NewBatchHandler(BatchOptions{
		BatchSize: 3,
		BeforeAddFunc: func(_ context.Context, log *Log) bool {
			beforeLogs = append(beforeLogs, log)

			// skip test2 log
			return log.Message != "test2"
		},
		WriteFunc: func(_ context.Context, logs []*Log) error {
			writeLogs = logs
			return nil
		},
	})

	h.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "test1", 0))
	h.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "test2", 0))
	h.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "test3", 0))

	// no batch write
	{
		checkLogMessages([]string{"test1", "test2", "test3"}, beforeLogs, t)

		checkLogMessages([]string{"test1", "test3"}, h.logs, t)

		// should be empty because no batch write has happened yet
		if totalWriteLogs := len(writeLogs); totalWriteLogs != 0 {
			t.Fatalf("Expected %d writeLogs, got %d", 0, totalWriteLogs)
		}
	}

	// add one more log to trigger the batch write
	{
		h.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "test4", 0))

		// should be empty after the batch write
		checkLogMessages([]string{}, h.logs, t)

		checkLogMessages([]string{"test1", "test3", "test4"}, writeLogs, t)
	}
}

func TestBatchHandlerWriteAll(t *testing.T) {
	ctx := context.Background()

	beforeLogs := []*Log{}
	writeLogs := []*Log{}

	h := NewBatchHandler(BatchOptions{
		BatchSize: 3,
		BeforeAddFunc: func(_ context.Context, log *Log) bool {
			beforeLogs = append(beforeLogs, log)
			return true
		},
		WriteFunc: func(_ context.Context, logs []*Log) error {
			writeLogs = logs
			return nil
		},
	})

	h.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "test1", 0))
	h.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "test2", 0))

	checkLogMessages([]string{"test1", "test2"}, beforeLogs, t)
	checkLogMessages([]string{"test1", "test2"}, h.logs, t)
	checkLogMessages([]string{}, writeLogs, t) // empty because the batch size hasn't been reached

	// force trigger the batch write
	h.WriteAll(ctx)

	checkLogMessages([]string{"test1", "test2"}, beforeLogs, t)
	checkLogMessages([]string{}, h.logs, t) // reset
	checkLogMessages([]string{"test1", "test2"}, writeLogs, t)
}

func TestBatchHandlerAttrsFormat(t *testing.T) {
	ctx := context.Background()

	beforeLogs := []*Log{}

	h0 := NewBatchHandler(BatchOptions{
		BeforeAddFunc: func(_ context.Context, log *Log) bool {
			beforeLogs = append(beforeLogs, log)
			return true
		},
		WriteFunc: func(_ context.Context, logs []*Log) error {
			return nil
		},
	})

	h1 := h0.WithAttrs([]slog.Attr{slog.Int("a", 1), slog.String("b", "123")})

	h2 := h1.WithGroup("sub").WithAttrs([]slog.Attr{
		slog.Int("c", 3),
		slog.Any("d", map[string]any{"d.1": 1}),
		slog.Any("e", errors.New("example error")),
	})

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)
	record.AddAttrs(slog.String("name", "test"))

	h0.Handle(ctx, record)
	h1.Handle(ctx, record)
	h2.Handle(ctx, record)

	// errors serialization checks
	errorsRecord := slog.NewRecord(time.Now(), slog.LevelError, "details", 0)
	errorsRecord.Add("validation.Errors", validation.Errors{
		"a": validation.NewError("validation_code", "validation_message"),
		"b": errors.New("plain"),
	})
	errorsRecord.Add("wrapped_validation.Errors", fmt.Errorf("wrapped: %w", validation.Errors{
		"a": validation.NewError("validation_code", "validation_message"),
		"b": errors.New("plain"),
	}))
	errorsRecord.Add("map[string]any", map[string]any{
		"a": validation.NewError("validation_code", "validation_message"),
		"b": errors.New("plain"),
		"c": "test_any",
		"d": map[string]any{
			"nestedA": validation.NewError("nested_code", "nested_message"),
			"nestedB": errors.New("nested_plain"),
		},
	})
	errorsRecord.Add("map[string]error", map[string]error{
		"a": validation.NewError("validation_code", "validation_message"),
		"b": errors.New("plain"),
	})
	errorsRecord.Add("map[string]validation.Error", map[string]validation.Error{
		"a": validation.NewError("validation_code", "validation_message"),
		"b": nil,
	})
	errorsRecord.Add("plain_error", errors.New("plain"))
	h0.Handle(ctx, errorsRecord)

	expected := []string{
		`{"name":"test"}`,
		`{"a":1,"b":"123","name":"test"}`,
		`{"a":1,"b":"123","sub":{"c":3,"d":{"d.1":1},"e":"example error","name":"test"}}`,
		`{"map[string]any":{"a":"validation_message","b":"plain","c":"test_any","d":{"nestedA":"nested_message","nestedB":"nested_plain"}},"map[string]error":{"a":"validation_message","b":"plain"},"map[string]validation.Error":{"a":"validation_message","b":null},"plain_error":"plain","validation.Errors":{"a":"validation_message","b":"plain"},"wrapped_validation.Errors":{"data":{"a":"validation_message","b":"plain"},"raw":"wrapped: a: validation_message; b: plain."}}`,
	}

	if len(beforeLogs) != len(expected) {
		t.Fatalf("Expected %d logs, got %d", len(expected), len(beforeLogs))
	}

	for i, data := range expected {
		t.Run(fmt.Sprintf("log handler %d", i), func(t *testing.T) {
			log := beforeLogs[i]
			raw, _ := log.Data.MarshalJSON()
			if string(raw) != data {
				t.Fatalf("Expected \n%s \ngot \n%s", data, raw)
			}
		})
	}
}

func checkLogMessages(expected []string, logs []*Log, t *testing.T) {
	if len(logs) != len(expected) {
		t.Fatalf("Expected %d batched logs, got %d (expected: %v)", len(expected), len(logs), expected)
	}

	for _, message := range expected {
		exists := false
		for _, l := range logs {
			if l.Message == message {
				exists = true
				continue
			}
		}
		if !exists {
			t.Fatalf("Missing %q log message", message)
		}
	}
}
