package core

import (
	"context"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
)

// LambdaFunctionContext represents the execution context for a lambda function.
type LambdaFunctionContext struct {
	// Base context for cancellation and timeout
	Context context.Context

	// App reference
	App App

	// Function being executed
	Function *LambdaFunction

	// Trigger information
	TriggerType   string
	TriggerConfig types.JSONRaw

	// HTTP-specific context (for HTTP triggers)
	HTTPRequest  *http.Request
	HTTPResponse http.ResponseWriter

	// Database-specific context (for database triggers)
	Collection     *Collection
	Record         *Record
	OldRecord      *Record // for update events
	DatabaseEvent  string  // insert, update, delete

	// Cron-specific context (for cron triggers)
	ScheduledTime time.Time

	// Execution metadata
	RequestID   string
	StartTime   time.Time
	Environment map[string]string // merged from function EnvVars and system env
}

// LambdaFunctionResult represents the result of a lambda function execution.
type LambdaFunctionResult struct {
	// Success indicates whether the function executed successfully
	Success bool `json:"success"`

	// Error message if the function failed
	Error string `json:"error,omitempty"`

	// Response data for HTTP triggers
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`

	// Execution metadata
	Duration  time.Duration `json:"duration"`
	RequestID string        `json:"requestId"`

	// Logs captured during execution
	Logs []LambdaFunctionLog `json:"logs,omitempty"`
}

// LambdaFunctionLog represents a log entry from function execution.
type LambdaFunctionLog struct {
	Level     string    `json:"level"`     // debug, info, warn, error
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

// NewLambdaFunctionContext creates a new lambda function execution context.
func NewLambdaFunctionContext(app App, function *LambdaFunction) *LambdaFunctionContext {
	return &LambdaFunctionContext{
		Context:     context.Background(),
		App:         app,
		Function:    function,
		StartTime:   time.Now(),
		RequestID:   generateRequestID(),
		Environment: mergeEnvironment(function.EnvVars, getSystemEnv()),
	}
}

// WithTimeout sets a timeout for the function execution.
func (ctx *LambdaFunctionContext) WithTimeout(timeout time.Duration) *LambdaFunctionContext {
	ctx.Context, _ = context.WithTimeout(ctx.Context, timeout)
	return ctx
}

// WithHTTPTrigger configures the context for an HTTP trigger.
func (ctx *LambdaFunctionContext) WithHTTPTrigger(r *http.Request, w http.ResponseWriter, config types.JSONRaw) *LambdaFunctionContext {
	ctx.TriggerType = TriggerTypeHTTP
	ctx.TriggerConfig = config
	ctx.HTTPRequest = r
	ctx.HTTPResponse = w
	return ctx
}

// WithDatabaseTrigger configures the context for a database trigger.
func (ctx *LambdaFunctionContext) WithDatabaseTrigger(collection *Collection, record, oldRecord *Record, event string, config types.JSONRaw) *LambdaFunctionContext {
	ctx.TriggerType = TriggerTypeDatabase
	ctx.TriggerConfig = config
	ctx.Collection = collection
	ctx.Record = record
	ctx.OldRecord = oldRecord
	ctx.DatabaseEvent = event
	return ctx
}

// WithCronTrigger configures the context for a cron trigger.
func (ctx *LambdaFunctionContext) WithCronTrigger(scheduledTime time.Time, config types.JSONRaw) *LambdaFunctionContext {
	ctx.TriggerType = TriggerTypeCron
	ctx.TriggerConfig = config
	ctx.ScheduledTime = scheduledTime
	return ctx
}

// Helper functions

func generateRequestID() string {
	// Simple implementation - in production you might want to use a UUID
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func mergeEnvironment(functionEnv types.JSONMap[any], systemEnv map[string]string) map[string]string {
	result := make(map[string]string)
	
	// Copy system environment
	for k, v := range systemEnv {
		result[k] = v
	}
	
	// Override with function-specific environment
	for k, v := range functionEnv {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	
	return result
}

func getSystemEnv() map[string]string {
	// Return a filtered set of system environment variables
	// that are safe to expose to lambda functions
	return map[string]string{
		"PB_VERSION": "lambda-functions-v1",
		"RUNTIME":    "javascript",
	}
}