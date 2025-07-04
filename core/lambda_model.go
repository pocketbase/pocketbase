package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/types"
)

var (
	_ Model      = (*LambdaFunction)(nil)
	_ DBExporter = (*LambdaFunction)(nil)
)

const CollectionNameLambdaFunctions = "lambdas"

const (
	// Trigger types
	TriggerTypeHTTP     = "http"
	TriggerTypeDatabase = "database"
	TriggerTypeCron     = "cron"

	// Database trigger events
	DatabaseEventInsert = "insert"
	DatabaseEventUpdate = "update"
	DatabaseEventDelete = "delete"

	// Default timeout in milliseconds (30 seconds)
	DefaultFunctionTimeout = 30000
	
	// Maximum timeout in milliseconds (5 minutes)
	MaxFunctionTimeout = 300000
)

// LambdaFunction defines a lambda function that can be triggered by various events.
type LambdaFunction struct {
	BaseModel

	Name    string         `db:"name" json:"name"`
	Code    string         `db:"code" json:"code"`
	Enabled bool           `db:"enabled" json:"enabled"`
	Timeout int            `db:"timeout" json:"timeout"` // in milliseconds
	Created types.DateTime `db:"created" json:"created"`
	Updated types.DateTime `db:"updated" json:"updated"`

	// Triggers stores the trigger configurations as JSON
	Triggers types.JSONArray[TriggerConfig] `db:"triggers" json:"triggers"`

	// EnvVars stores environment variables as key-value pairs
	EnvVars types.JSONMap[any] `db:"envVars" json:"envVars"`
}

// TriggerConfig represents a single trigger configuration
type TriggerConfig struct {
	Type   string         `json:"type"`
	Config types.JSONRaw  `json:"config"`
}

// HTTPTriggerConfig represents HTTP route trigger configuration
type HTTPTriggerConfig struct {
	Method string `json:"method"` // GET, POST, PUT, DELETE, etc.
	Path   string `json:"path"`   // e.g., "/api/functions/hello"
}

// DatabaseTriggerConfig represents database event trigger configuration
type DatabaseTriggerConfig struct {
	Collection string   `json:"collection"`
	Events     []string `json:"events"` // insert, update, delete
}

// CronTriggerConfig represents cron schedule trigger configuration
type CronTriggerConfig struct {
	Expression string `json:"expression"` // cron expression
}

// TableName returns the LambdaFunction model SQL table name.
func (m *LambdaFunction) TableName() string {
	return "_pb_functions"
}

// PostValidate performs model validation.
func (m *LambdaFunction) PostValidate(ctx context.Context, app App) error {
	// Validate timeout
	if m.Timeout <= 0 {
		m.Timeout = DefaultFunctionTimeout
	}
	if m.Timeout > MaxFunctionTimeout {
		return fmt.Errorf("timeout cannot exceed %d milliseconds", MaxFunctionTimeout)
	}

	// Validate name
	if m.Name == "" {
		return errors.New("function name is required")
	}

	// Validate code
	if m.Code == "" {
		return errors.New("function code is required")
	}

	// Validate triggers
	for i, trigger := range m.Triggers {
		if err := validateTriggerConfig(trigger); err != nil {
			return fmt.Errorf("invalid trigger at index %d: %w", i, err)
		}
	}

	return nil
}

// DBExport implements the DBExporter interface.
func (m *LambdaFunction) DBExport(app App) (map[string]any, error) {
	result := map[string]any{
		"id":       m.Id,
		"name":     m.Name,
		"code":     m.Code,
		"enabled":  m.Enabled,
		"timeout":  m.Timeout,
		"triggers": m.Triggers,
		"envVars":  m.EnvVars,
	}

	if m.IsNew() {
		result["created"] = types.NowDateTime()
		result["updated"] = types.NowDateTime()
	} else {
		result["created"] = m.Created
		result["updated"] = types.NowDateTime()
	}

	return result, nil
}

// GetHTTPTriggers returns all HTTP trigger configurations
func (m *LambdaFunction) GetHTTPTriggers() ([]HTTPTriggerConfig, error) {
	var configs []HTTPTriggerConfig
	for _, trigger := range m.Triggers {
		if trigger.Type == TriggerTypeHTTP {
			var config HTTPTriggerConfig
			if err := json.Unmarshal(trigger.Config, &config); err != nil {
				return nil, err
			}
			configs = append(configs, config)
		}
	}
	return configs, nil
}

// GetDatabaseTriggers returns all database trigger configurations
func (m *LambdaFunction) GetDatabaseTriggers() ([]DatabaseTriggerConfig, error) {
	var configs []DatabaseTriggerConfig
	for _, trigger := range m.Triggers {
		if trigger.Type == TriggerTypeDatabase {
			var config DatabaseTriggerConfig
			if err := json.Unmarshal(trigger.Config, &config); err != nil {
				return nil, err
			}
			configs = append(configs, config)
		}
	}
	return configs, nil
}

// GetCronTriggers returns all cron trigger configurations
func (m *LambdaFunction) GetCronTriggers() ([]CronTriggerConfig, error) {
	var configs []CronTriggerConfig
	for _, trigger := range m.Triggers {
		if trigger.Type == TriggerTypeCron {
			var config CronTriggerConfig
			if err := json.Unmarshal(trigger.Config, &config); err != nil {
				return nil, err
			}
			configs = append(configs, config)
		}
	}
	return configs, nil
}

func validateTriggerConfig(trigger TriggerConfig) error {
	switch trigger.Type {
	case TriggerTypeHTTP:
		var config HTTPTriggerConfig
		if err := json.Unmarshal(trigger.Config, &config); err != nil {
			return fmt.Errorf("invalid HTTP trigger config: %w", err)
		}
		if config.Method == "" {
			return errors.New("HTTP method is required")
		}
		if config.Path == "" {
			return errors.New("HTTP path is required")
		}
	case TriggerTypeDatabase:
		var config DatabaseTriggerConfig
		if err := json.Unmarshal(trigger.Config, &config); err != nil {
			return fmt.Errorf("invalid database trigger config: %w", err)
		}
		if config.Collection == "" {
			return errors.New("collection name is required")
		}
		if len(config.Events) == 0 {
			return errors.New("at least one event type is required")
		}
		for _, event := range config.Events {
			switch event {
			case DatabaseEventInsert, DatabaseEventUpdate, DatabaseEventDelete:
				// valid event
			default:
				return fmt.Errorf("invalid database event type: %s", event)
			}
		}
	case TriggerTypeCron:
		var config CronTriggerConfig
		if err := json.Unmarshal(trigger.Config, &config); err != nil {
			return fmt.Errorf("invalid cron trigger config: %w", err)
		}
		if config.Expression == "" {
			return errors.New("cron expression is required")
		}
	default:
		return fmt.Errorf("unknown trigger type: %s", trigger.Type)
	}
	return nil
}