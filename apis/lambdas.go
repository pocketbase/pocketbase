package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// LambdaFunctionCreateRequest represents the request for creating a lambda function
type LambdaFunctionCreateRequest struct {
	Name        string                 `json:"name" form:"name"`
	Code        string                 `json:"code" form:"code"`
	Enabled     bool                   `json:"enabled" form:"enabled"`
	Timeout     int                    `json:"timeout" form:"timeout"`
	Triggers    map[string]interface{} `json:"triggers" form:"triggers"`
	EnvVars     map[string]string      `json:"env_vars" form:"env_vars"`
	Description string                 `json:"description" form:"description"`
}

// LambdaFunctionUpdateRequest represents the request for updating a lambda function
type LambdaFunctionUpdateRequest struct {
	Name        string                 `json:"name" form:"name"`
	Code        string                 `json:"code" form:"code"`
	Enabled     *bool                  `json:"enabled" form:"enabled"`
	Timeout     *int                   `json:"timeout" form:"timeout"`
	Triggers    map[string]interface{} `json:"triggers" form:"triggers"`
	EnvVars     map[string]string      `json:"env_vars" form:"env_vars"`
	Description string                 `json:"description" form:"description"`
}

// BindLambdaFunctionRoutes binds the lambda function API routes
func BindLambdaFunctionRoutes(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	api := lambdaFunctionAPI{app}

	subGroup := rg.Group("/lambdas").Bind(RequireSuperuserAuth())
	subGroup.GET("", api.list)
	subGroup.POST("", api.create)
	subGroup.GET("/{id}", api.view)
	subGroup.PATCH("/{id}", api.update)
	subGroup.DELETE("/{id}", api.delete)
	subGroup.POST("/{id}/execute", api.execute)
	subGroup.GET("/{id}/logs", api.logs)
	subGroup.POST("/{id}/enable", api.enable)
	subGroup.POST("/{id}/disable", api.disable)
}

type lambdaFunctionAPI struct {
	app core.App
}

func (api *lambdaFunctionAPI) list(e *core.RequestEvent) error {
	records, err := api.app.FindRecordsByFilter(core.CollectionNameLambdaFunctions, "", "-created", 0, 0)
	if err != nil {
		return e.BadRequestError("Failed to fetch lambda functions", err)
	}

	result := make([]map[string]interface{}, len(records))
	for i, record := range records {
		result[i] = map[string]interface{}{
			"id":          record.Id,
			"name":        record.GetString("name"),
			"enabled":     record.GetBool("enabled"),
			"timeout":     record.GetInt("timeout") / 1000, // Convert ms to seconds for display
			"description": record.GetString("description"),
			"created":     record.GetDateTime("created"),
			"updated":     record.GetDateTime("updated"),
		}
	}

	return e.JSON(200, result)
}

func (api *lambdaFunctionAPI) create(e *core.RequestEvent) error {
	form := &LambdaFunctionCreateRequest{}
	if err := e.BindBody(form); err != nil {
		return e.BadRequestError("Invalid request data", err)
	}

	// Validate required fields
	if form.Name == "" {
		return e.BadRequestError("Name is required", nil)
	}
	if form.Code == "" {
		return e.BadRequestError("Code is required", nil)
	}

	// Validate function name
	if !isValidFunctionName(form.Name) {
		return e.BadRequestError("Invalid function name. Must be alphanumeric with underscores and hyphens", nil)
	}

	// Check if function name already exists
	existing, err := api.app.FindFirstRecordByFilter(core.CollectionNameLambdaFunctions, "name = {:name}", map[string]any{
		"name": form.Name,
	})
	if err == nil && existing != nil {
		return e.BadRequestError("Function with this name already exists", nil)
	}

	// Validate timeout
	if form.Timeout == 0 {
		form.Timeout = 30 // Default 30 seconds
	}
	if form.Timeout < 1 || form.Timeout > 300 {
		return e.BadRequestError("Timeout must be between 1 and 300 seconds", nil)
	}
	
	// Convert seconds to milliseconds for storage
	timeoutMs := form.Timeout * 1000

	// Validate triggers
	if err := validateTriggers(form.Triggers); err != nil {
		return e.BadRequestError("Invalid trigger configuration", err)
	}

	// Create record
	collection, err := api.app.FindCollectionByNameOrId(core.CollectionNameLambdaFunctions)
	if err != nil {
		return e.BadRequestError("Functions collection not found", err)
	}

	record := core.NewRecord(collection)
	record.Set("name", form.Name)
	record.Set("code", form.Code)
	record.Set("enabled", form.Enabled)
	record.Set("timeout", timeoutMs)
	record.Set("description", form.Description)

	// Convert triggers to JSON
	triggersJSON, _ := json.Marshal(form.Triggers)
	record.Set("triggers", string(triggersJSON))

	// Convert env vars to JSON
	envVarsJSON, _ := json.Marshal(form.EnvVars)
	record.Set("envVars", string(envVarsJSON))

	if err := api.app.Save(record); err != nil {
		return e.BadRequestError("Failed to create lambda function", err)
	}

	return e.JSON(201, map[string]interface{}{
		"id":          record.Id,
		"name":        record.GetString("name"),
		"enabled":     record.GetBool("enabled"),
		"timeout":     record.GetInt("timeout") / 1000, // Convert ms to seconds for display
		"description": record.GetString("description"),
		"created":     record.GetDateTime("created"),
		"updated":     record.GetDateTime("updated"),
	})
}

func (api *lambdaFunctionAPI) view(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	record, err := api.app.FindRecordById(core.CollectionNameLambdaFunctions, id)
	if err != nil {
		return e.NotFoundError("Lambda function not found", err)
	}

	return e.JSON(200, map[string]interface{}{
		"id":          record.Id,
		"name":        record.GetString("name"),
		"code":        record.GetString("code"),
		"enabled":     record.GetBool("enabled"),
		"timeout":     record.GetInt("timeout") / 1000, // Convert ms to seconds for display
		"triggers":    record.GetString("triggers"),
		"env_vars":    record.GetString("envVars"),
		"description": record.GetString("description"),
		"created":     record.GetDateTime("created"),
		"updated":     record.GetDateTime("updated"),
	})
}

func (api *lambdaFunctionAPI) update(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	record, err := api.app.FindRecordById(core.CollectionNameLambdaFunctions, id)
	if err != nil {
		return e.NotFoundError("Lambda function not found", err)
	}

	form := &LambdaFunctionUpdateRequest{}
	if err := e.BindBody(form); err != nil {
		return e.BadRequestError("Invalid request data", err)
	}

	// Update fields if provided
	if form.Name != "" {
		if !isValidFunctionName(form.Name) {
			return e.BadRequestError("Invalid function name", nil)
		}
		record.Set("name", form.Name)
	}

	if form.Code != "" {
		record.Set("code", form.Code)
	}

	if form.Enabled != nil {
		record.Set("enabled", *form.Enabled)
	}

	if form.Timeout != nil {
		if *form.Timeout < 1 || *form.Timeout > 300 {
			return e.BadRequestError("Timeout must be between 1 and 300 seconds", nil)
		}
		// Convert seconds to milliseconds for storage
		timeoutMs := *form.Timeout * 1000
		record.Set("timeout", timeoutMs)
	}

	if form.Description != "" {
		record.Set("description", form.Description)
	}

	if form.Triggers != nil {
		if err := validateTriggers(form.Triggers); err != nil {
			return e.BadRequestError("Invalid trigger configuration", err)
		}
		triggersJSON, _ := json.Marshal(form.Triggers)
		record.Set("triggers", string(triggersJSON))
	}

	if form.EnvVars != nil {
		envVarsJSON, _ := json.Marshal(form.EnvVars)
		record.Set("envVars", string(envVarsJSON))
	}

	if err := api.app.Save(record); err != nil {
		return e.BadRequestError("Failed to update lambda function", err)
	}

	return e.JSON(200, map[string]interface{}{
		"id":          record.Id,
		"name":        record.GetString("name"),
		"enabled":     record.GetBool("enabled"),
		"timeout":     record.GetInt("timeout") / 1000, // Convert ms to seconds for display
		"description": record.GetString("description"),
		"created":     record.GetDateTime("created"),
		"updated":     record.GetDateTime("updated"),
	})
}

func (api *lambdaFunctionAPI) delete(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	record, err := api.app.FindRecordById(core.CollectionNameLambdaFunctions, id)
	if err != nil {
		return e.NotFoundError("Lambda function not found", err)
	}

	if err := api.app.Delete(record); err != nil {
		return e.BadRequestError("Failed to delete lambda function", err)
	}

	return e.NoContent(204)
}

func (api *lambdaFunctionAPI) execute(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	record, err := api.app.FindRecordById(core.CollectionNameLambdaFunctions, id)
	if err != nil {
		return e.NotFoundError("Lambda function not found", err)
	}

	if !record.GetBool("enabled") {
		return e.BadRequestError("Lambda function is disabled", nil)
	}

	// Get the JavaScript code
	code := record.GetString("code")
	functionName := record.GetString("name")
	timeoutMs := record.GetInt("timeout")

	// Execute the function
	startTime := time.Now()
	
	result, err := executeLambdaFunction(api.app, code, functionName, timeoutMs, map[string]interface{}{
		"request": map[string]interface{}{
			"method": e.Request.Method,
			"url":    e.Request.URL.String(),
			"headers": extractHeaders(e.Request),
		},
	})

	duration := time.Since(startTime)

	// Log the execution
	success := err == nil
	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}

	logExecution(api.app, record.Id, functionName, "manual", success, result, errorMsg, duration, map[string]interface{}{
		"request": map[string]interface{}{
			"method": e.Request.Method,
			"url":    e.Request.URL.String(),
		},
	})

	if err != nil {
		return e.JSON(200, map[string]interface{}{
			"success":     false,
			"error":       err.Error(),
			"duration_ms": duration.Milliseconds(),
			"timestamp":   time.Now(),
		})
	}

	return e.JSON(200, map[string]interface{}{
		"success":     true,
		"output":      result,
		"duration_ms": duration.Milliseconds(),
		"timestamp":   time.Now(),
	})
}

func (api *lambdaFunctionAPI) logs(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	_, err := api.app.FindRecordById(core.CollectionNameLambdaFunctions, id)
	if err != nil {
		return e.NotFoundError("Lambda function not found", err)
	}

	// Get logs from database
	logRecords, err := api.app.FindRecordsByFilter("lambda_logs", 
		"function_id = {:id}", "-created", 50, 0, map[string]any{"id": id})
	if err != nil {
		return e.BadRequestError("Failed to fetch logs", err)
	}

	logs := make([]map[string]interface{}, len(logRecords))
	for i, record := range logRecords {
		logs[i] = map[string]interface{}{
			"id":           record.Id,
			"function_id":  record.GetString("function_id"),
			"function_name": record.GetString("function_name"),
			"trigger_type": record.GetString("trigger_type"),
			"success":      record.GetBool("success"),
			"output":       record.Get("output"),
			"error":        record.GetString("error"),
			"duration_ms":  record.GetInt("duration_ms"),
			"context":      record.Get("context"),
			"timestamp":    record.GetDateTime("created"),
		}
	}

	return e.JSON(200, map[string]interface{}{
		"logs":  logs,
		"total": len(logs),
	})
}

func (api *lambdaFunctionAPI) enable(e *core.RequestEvent) error {
	return api.toggleEnabled(e, true)
}

func (api *lambdaFunctionAPI) disable(e *core.RequestEvent) error {
	return api.toggleEnabled(e, false)
}

func (api *lambdaFunctionAPI) toggleEnabled(e *core.RequestEvent, enabled bool) error {
	id := e.Request.PathValue("id")
	record, err := api.app.FindRecordById(core.CollectionNameLambdaFunctions, id)
	if err != nil {
		return e.NotFoundError("Lambda function not found", err)
	}

	record.Set("enabled", enabled)
	if err := api.app.Save(record); err != nil {
		return e.BadRequestError("Failed to update lambda function", err)
	}

	return e.JSON(200, map[string]interface{}{
		"id":      record.Id,
		"name":    record.GetString("name"),
		"enabled": record.GetBool("enabled"),
	})
}

// Helper functions

func isValidFunctionName(name string) bool {
	if len(name) == 0 || len(name) > 50 {
		return false
	}
	
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '_' || char == '-') {
			return false
		}
	}
	
	return true
}

func validateTriggers(triggers map[string]interface{}) error {
	if triggers == nil {
		return nil
	}

	// Validate HTTP triggers
	if httpTriggers, ok := triggers["http"]; ok {
		if httpList, ok := httpTriggers.([]interface{}); ok {
			for _, trigger := range httpList {
				if httpTrigger, ok := trigger.(map[string]interface{}); ok {
					method, hasMethod := httpTrigger["method"]
					path, hasPath := httpTrigger["path"]
					
					if !hasMethod || !hasPath {
						return fmt.Errorf("HTTP trigger must have method and path")
					}
					
					if !isValidHTTPMethod(method.(string)) {
						return fmt.Errorf("Invalid HTTP method: %s", method)
					}
					
					if !isValidPath(path.(string)) {
						return fmt.Errorf("Invalid path: %s", path)
					}
				}
			}
		}
	}

	// Validate database triggers
	if dbTriggers, ok := triggers["database"]; ok {
		if dbList, ok := dbTriggers.([]interface{}); ok {
			for _, trigger := range dbList {
				if dbTrigger, ok := trigger.(map[string]interface{}); ok {
					_, hasCollection := dbTrigger["collection"]
					event, hasEvent := dbTrigger["event"]
					
					if !hasCollection || !hasEvent {
						return fmt.Errorf("Database trigger must have collection and event")
					}
					
					if !isValidDBEvent(event.(string)) {
						return fmt.Errorf("Invalid database event: %s", event)
					}
				}
			}
		}
	}

	// Validate cron triggers
	if cronTriggers, ok := triggers["cron"]; ok {
		if cronList, ok := cronTriggers.([]interface{}); ok {
			for _, trigger := range cronList {
				if cronTrigger, ok := trigger.(map[string]interface{}); ok {
					schedule, hasSchedule := cronTrigger["schedule"]
					
					if !hasSchedule {
						return fmt.Errorf("Cron trigger must have schedule")
					}
					
					if !isValidCronSchedule(schedule.(string)) {
						return fmt.Errorf("Invalid cron schedule: %s", schedule)
					}
				}
			}
		}
	}

	return nil
}

func isValidHTTPMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	method = strings.ToUpper(method)
	
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	
	return false
}

func isValidPath(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}
	
	// Basic path validation - can be enhanced
	return len(path) > 0 && len(path) <= 100
}

func isValidDBEvent(event string) bool {
	validEvents := []string{"create", "update", "delete"}
	
	for _, validEvent := range validEvents {
		if event == validEvent {
			return true
		}
	}
	
	return false
}

func isValidCronSchedule(schedule string) bool {
	// Basic cron validation - should be enhanced with proper cron parser
	parts := strings.Fields(schedule)
	return len(parts) == 5 || len(parts) == 6
}

// executeLambdaFunction executes JavaScript code using Goja runtime
func executeLambdaFunction(app core.App, code, functionName string, timeoutMs int, ctxData map[string]interface{}) (interface{}, error) {
	vm := goja.New()
	
	// Set up timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()
	
	// Set up console.log functionality
	console := vm.NewObject()
	console.Set("log", func(args ...interface{}) {
		fmt.Printf("[%s] ", functionName)
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(arg)
		}
		fmt.Println()
	})
	vm.Set("console", console)
	
	// Set up PocketBase $app object
	pbApp := vm.NewObject()
	
	// Add findRecordById function
	pbApp.Set("findRecordById", func(collection, id string) interface{} {
		record, err := app.FindRecordById(collection, id)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return recordToMap(record)
	})
	
	// Add findRecordByFilter function
	pbApp.Set("findRecordByFilter", func(collection, filter string) interface{} {
		record, err := app.FindFirstRecordByFilter(collection, filter)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return recordToMap(record)
	})
	
	// Add findRecords function
	pbApp.Set("findRecords", func(collection string, args ...interface{}) interface{} {
		var filter string
		var sort string = "-created"
		var limit int = 30
		var offset int = 0
		
		if len(args) > 0 {
			if f, ok := args[0].(string); ok {
				filter = f
			}
		}
		if len(args) > 1 {
			if s, ok := args[1].(string); ok {
				sort = s
			}
		}
		if len(args) > 2 {
			if l, ok := args[2].(int64); ok {
				limit = int(l)
			}
		}
		if len(args) > 3 {
			if o, ok := args[3].(int64); ok {
				offset = int(o)
			}
		}
		
		records, err := app.FindRecordsByFilter(collection, filter, sort, limit, offset)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		
		result := make([]interface{}, len(records))
		for i, record := range records {
			result[i] = recordToMap(record)
		}
		return result
	})
	
	vm.Set("$app", pbApp)
	
	// Set context variables
	for key, value := range ctxData {
		vm.Set(key, value)
	}
	
	// Execute with timeout
	done := make(chan struct{})
	var result goja.Value
	var err error
	
	go func() {
		defer close(done)
		result, err = vm.RunString(code)
	}()
	
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("execution timeout after %dms", timeoutMs)
	case <-done:
		if err != nil {
			return nil, fmt.Errorf("execution error: %v", err)
		}
		
		if result != nil {
			return result.Export(), nil
		}
		return nil, nil
	}
}

// recordToMap converts a PocketBase record to a map for JavaScript access
func recordToMap(record *core.Record) map[string]interface{} {
	result := map[string]interface{}{
		"id":      record.Id,
		"created": record.GetDateTime("created"),
		"updated": record.GetDateTime("updated"),
	}
	
	// Add all field values
	for name := range record.Collection().Fields.AsMap() {
		result[name] = record.Get(name)
	}
	
	return result
}

// logExecution logs function execution to the database
func logExecution(app core.App, functionId, functionName, triggerType string, success bool, output interface{}, errorMsg string, duration time.Duration, context map[string]interface{}) {
	collection, err := app.FindCollectionByNameOrId("lambda_logs")
	if err != nil {
		fmt.Printf("Failed to find logs collection: %v\n", err)
		return
	}

	record := core.NewRecord(collection)
	record.Set("function_id", functionId)
	record.Set("function_name", functionName)
	record.Set("trigger_type", triggerType)
	record.Set("success", success)
	record.Set("duration_ms", int(duration.Milliseconds()))
	record.Set("error", errorMsg)

	// Convert output to JSON
	if output != nil {
		outputJSON, _ := json.Marshal(output)
		record.Set("output", string(outputJSON))
	}

	// Convert context to JSON
	if context != nil {
		contextJSON, _ := json.Marshal(context)
		record.Set("context", string(contextJSON))
	}

	if err := app.Save(record); err != nil {
		fmt.Printf("Failed to save execution log: %v\n", err)
	}
}

// extractHeaders extracts HTTP headers as a map
func extractHeaders(req *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	return headers
}