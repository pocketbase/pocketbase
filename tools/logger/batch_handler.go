package logger

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ slog.Handler = (*BatchHandler)(nil)

// BatchOptions are options for the BatchHandler.
type BatchOptions struct {
	// WriteFunc processes the batched logs.
	WriteFunc func(ctx context.Context, logs []*Log) error

	// BeforeAddFunc is optional function that is invoked every time
	// before a new log is added to the batch queue.
	//
	// Return false to skip adding the log into the batch queue.
	BeforeAddFunc func(ctx context.Context, log *Log) bool

	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler

	// BatchSize specifies how many logs to accumulate before calling WriteFunc.
	// If not set or 0, fallback to 100 by default.
	BatchSize int
}

// NewBatchHandler creates a slog compatible handler that writes JSON
// logs on batches (default to 100), using the given options.
//
// Panics if [BatchOptions.WriteFunc] is not defined.
//
// Example:
//
//	l := slog.New(logger.NewBatchHandler(logger.BatchOptions{
//	    WriteFunc: func(ctx context.Context, logs []*Log) error {
//	        for _, l := range logs {
//	            fmt.Println(l.Level, l.Message, l.Data)
//	        }
//	        return nil
//	    }
//	}))
//	l.Info("Example message", "title", "lorem ipsum")
func NewBatchHandler(options BatchOptions) *BatchHandler {
	h := &BatchHandler{
		mux:     &sync.Mutex{},
		options: &options,
	}

	if h.options.WriteFunc == nil {
		panic("options.WriteFunc must be set")
	}

	if h.options.Level == nil {
		h.options.Level = slog.LevelInfo
	}

	if h.options.BatchSize == 0 {
		h.options.BatchSize = 100
	}

	h.logs = make([]*Log, 0, h.options.BatchSize)

	return h
}

// BatchHandler is a slog handler that writes records on batches.
//
// The log records attributes are formatted in JSON.
//
// Requires the [BatchOptions.WriteFunc] option to be defined.
type BatchHandler struct {
	mux     *sync.Mutex
	parent  *BatchHandler
	options *BatchOptions
	group   string
	attrs   []slog.Attr
	logs    []*Log
}

// Enabled reports whether the handler handles records at the given level.
//
// The handler ignores records whose level is lower.
func (h *BatchHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

// WithGroup returns a new BatchHandler that starts a group.
//
// All logger attributes will be resolved under the specified group name.
func (h *BatchHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return &BatchHandler{
		parent:  h,
		mux:     h.mux,
		options: h.options,
		group:   name,
	}
}

// WithAttrs returns a new BatchHandler loaded with the specified attributes.
func (h *BatchHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	return &BatchHandler{
		parent:  h,
		mux:     h.mux,
		options: h.options,
		attrs:   attrs,
	}
}

// Handle formats the slog.Record argument as JSON object and adds it
// to the batch queue.
//
// If the batch queue threshold has been reached, the WriteFunc option
// is invoked with the accumulated logs which in turn will reset the batch queue.
func (h *BatchHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.group != "" {
		h.mux.Lock()
		attrs := make([]any, 0, len(h.attrs)+r.NumAttrs())
		for _, a := range h.attrs {
			attrs = append(attrs, a)
		}
		h.mux.Unlock()

		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
			return true
		})

		r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
		r.AddAttrs(slog.Group(h.group, attrs...))
	} else if len(h.attrs) > 0 {
		r = r.Clone()

		h.mux.Lock()
		r.AddAttrs(h.attrs...)
		h.mux.Unlock()
	}

	if h.parent != nil {
		return h.parent.Handle(ctx, r)
	}

	data := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		if err := h.resolveAttr(data, a); err != nil {
			return false
		}
		return true
	})

	log := &Log{
		Time:    r.Time,
		Level:   r.Level,
		Message: r.Message,
		Data:    types.JSONMap[any](data),
	}

	if h.options.BeforeAddFunc != nil && !h.options.BeforeAddFunc(ctx, log) {
		return nil
	}

	h.mux.Lock()
	h.logs = append(h.logs, log)
	totalLogs := len(h.logs)
	h.mux.Unlock()

	if totalLogs >= h.options.BatchSize {
		if err := h.WriteAll(ctx); err != nil {
			return err
		}
	}

	return nil
}

// SetLevel updates the handler options level to the specified one.
func (h *BatchHandler) SetLevel(level slog.Level) {
	h.mux.Lock()
	h.options.Level = level
	h.mux.Unlock()
}

// WriteAll writes all accumulated Log entries and resets the batch queue.
func (h *BatchHandler) WriteAll(ctx context.Context) error {
	if h.parent != nil {
		// invoke recursively the parent level handler since the most
		// top level one is holding the logs queue.
		return h.parent.WriteAll(ctx)
	}

	h.mux.Lock()

	totalLogs := len(h.logs)

	// no logs to write
	if totalLogs == 0 {
		h.mux.Unlock()
		return nil
	}

	// create a copy of the logs slice to prevent blocking during write
	logs := make([]*Log, totalLogs)
	copy(logs, h.logs)
	h.logs = h.logs[:0] // reset

	h.mux.Unlock()

	return h.options.WriteFunc(ctx, logs)
}

// resolveAttr writes attr into data.
func (h *BatchHandler) resolveAttr(data map[string]any, attr slog.Attr) error {
	// ensure that the attr value is resolved before doing anything else
	attr.Value = attr.Value.Resolve()

	if attr.Equal(slog.Attr{}) {
		return nil // ignore empty attrs
	}

	switch attr.Value.Kind() {
	case slog.KindGroup:
		attrs := attr.Value.Group()
		if len(attrs) == 0 {
			return nil // ignore empty groups
		}

		// create a submap to wrap the resolved group attributes
		groupData := make(map[string]any, len(attrs))

		for _, subAttr := range attrs {
			h.resolveAttr(groupData, subAttr)
		}

		if len(groupData) > 0 {
			data[attr.Key] = groupData
		}
	default:
		data[attr.Key] = normalizeLogAttrValue(attr.Value.Any())
	}

	return nil
}

func normalizeLogAttrValue(rawAttrValue any) any {
	switch attrV := rawAttrValue.(type) {
	case validation.Errors:
		out := make(map[string]any, len(attrV))
		for k, v := range attrV {
			out[k] = serializeLogError(v)
		}
		return out
	case map[string]validation.Error:
		out := make(map[string]any, len(attrV))
		for k, v := range attrV {
			out[k] = serializeLogError(v)
		}
		return out
	case map[string]error:
		out := make(map[string]any, len(attrV))
		for k, v := range attrV {
			out[k] = serializeLogError(v)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(attrV))
		for k, v := range attrV {
			switch vv := v.(type) {
			case error:
				out[k] = serializeLogError(vv)
			default:
				out[k] = normalizeLogAttrValue(vv)
			}
		}
		return out
	case error:
		// check for wrapped validation.Errors
		var ve validation.Errors
		if errors.As(attrV, &ve) {
			out := make(map[string]any, len(ve))
			for k, v := range ve {
				out[k] = serializeLogError(v)
			}
			return map[string]any{
				"data": out,
				"raw":  serializeLogError(attrV),
			}
		}
		return serializeLogError(attrV)
	default:
		return attrV
	}
}

func serializeLogError(err error) any {
	if err == nil {
		return nil
	}

	// prioritize a json structured format (e.g. validation.Errors)
	jsonErr, ok := err.(json.Marshaler)
	if ok {
		return jsonErr
	}

	// fallback to its original string representation
	return err.Error()
}
