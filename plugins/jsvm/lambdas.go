package jsvm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/buffer"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/template"
)

// LambdaFunctionPluginConfig defines the configuration for the lambda function plugin
type LambdaFunctionPluginConfig struct {
	// PoolSize specifies how many goja.Runtime instances to prewarm
	// for lambda function execution
	PoolSize int

	// MaxExecutionTime specifies the maximum execution time for lambda functions
	MaxExecutionTime time.Duration

	// MaxMemory specifies the maximum memory usage for lambda functions (in bytes)
	MaxMemory int64

	// OnInit allows custom initialization of the JS runtime
	OnInit func(vm *goja.Runtime)
}

// LambdaFunctionPlugin manages lambda function execution
type LambdaFunctionPlugin struct {
	app           core.App
	config        LambdaFunctionPluginConfig
	executors     *vmsPool
	scheduler     *cron.Cron
	httpRoutes    sync.Map // map[string]*LambdaFunctionHTTPRoute
	dbTriggers    sync.Map // map[string][]*LambdaFunctionDBTrigger
	cronJobs      sync.Map // map[string]*LambdaFunctionCronJob
	templateRegistry *template.Registry
	requireRegistry  *require.Registry
}

// LambdaFunctionHTTPRoute represents an HTTP route for an lambda function
type LambdaFunctionHTTPRoute struct {
	FunctionID string
	Method     string
	Path       string
	Handler    func(*core.RequestEvent) error
}

// LambdaFunctionDBTrigger represents a database trigger for an lambda function
type LambdaFunctionDBTrigger struct {
	FunctionID string
	Collection string
	Event      string // "create", "update", "delete"
}

// LambdaFunctionCronJob represents a cron job for an lambda function
type LambdaFunctionCronJob struct {
	FunctionID string
	Schedule   string
	JobID      string
}

// LambdaFunctionExecutionContext provides context for lambda function execution
type LambdaFunctionExecutionContext struct {
	FunctionID   string
	TriggerType  string
	Request      *http.Request
	Response     http.ResponseWriter
	Record       interface{}
	OldRecord    interface{}
	Environment  map[string]string
	StartTime    time.Time
}

// LambdaFunctionExecutionResult represents the result of lambda function execution
type LambdaFunctionExecutionResult struct {
	Success   bool
	Output    interface{}
	Error     string
	Duration  time.Duration
	Memory    int64
}

// RegisterLambdaFunctionPlugin registers the lambda function plugin with the app
func RegisterLambdaFunctionPlugin(app core.App, config LambdaFunctionPluginConfig) (*LambdaFunctionPlugin, error) {
	if config.MaxExecutionTime == 0 {
		config.MaxExecutionTime = 30 * time.Second
	}
	if config.MaxMemory == 0 {
		config.MaxMemory = 128 * 1024 * 1024 // 128MB
	}

	plugin := &LambdaFunctionPlugin{
		app:              app,
		config:           config,
		scheduler:        cron.New(),
		templateRegistry: template.NewRegistry(),
		requireRegistry:  new(require.Registry),
	}

	// Initialize VM pool
	plugin.executors = newPool(config.PoolSize, plugin.createVM)

	// Register app lifecycle hooks
	plugin.registerLifecycleHooks()

	// Load existing lambda functions after database is ready
	plugin.app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return plugin.loadLambdaFunctions()
	})

	return plugin, nil
}

// createVM creates a new goja.Runtime instance for lambda function execution
func (p *LambdaFunctionPlugin) createVM() *goja.Runtime {
	vm := goja.New()

	// Enable Node.js compatibility
	p.requireRegistry.Enable(vm)
	console.Enable(vm)
	process.Enable(vm)
	buffer.Enable(vm)

	// Add PocketBase bindings
	baseBinds(vm)
	dbxBinds(vm)
	filesystemBinds(vm)
	securityBinds(vm)
	osBinds(vm)
	filepathBinds(vm)
	httpClientBinds(vm)
	formsBinds(vm)
	apisBinds(vm)
	mailsBinds(vm)

	// Add lambda function specific bindings
	vm.Set("$app", p.app)
	vm.Set("$template", p.templateRegistry)

	// Custom initialization
	if p.config.OnInit != nil {
		p.config.OnInit(vm)
	}

	return vm
}

// registerLifecycleHooks registers the necessary app lifecycle hooks
func (p *LambdaFunctionPlugin) registerLifecycleHooks() {
	// Register HTTP routes on serve
	p.app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		p.registerHTTPRoutes(e)
		return e.Next()
	})

	// Register database triggers
	p.registerDatabaseTriggers()

	// Start cron scheduler
	p.app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		p.scheduler.Start()
		return e.Next()
	})

	// Stop cron scheduler on termination
	p.app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		p.scheduler.Stop()
		return e.Next()
	})

	// Handle lambda function CRUD operations
	p.app.OnRecordCreate("lambdas").BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return p.handleFunctionCreated(e.Record)
	})

	p.app.OnRecordUpdate("lambdas").BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return p.handleFunctionUpdated(e.Record)
	})

	p.app.OnRecordDelete("lambdas").BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return p.handleFunctionDeleted(e.Record)
	})
}

// loadLambdaFunctions loads all existing lambda functions from the database
func (p *LambdaFunctionPlugin) loadLambdaFunctions() error {
	// Check if the collection exists first
	_, err := p.app.FindCollectionByNameOrId("lambdas")
	if err != nil {
		// Collection doesn't exist yet, this is fine
		p.app.Logger().Debug("Lambda functions collection not found, skipping loading")
		return nil
	}

	functions, err := p.app.FindRecordsByFilter("lambdas", "enabled = true", "", 0, 0)
	if err != nil {
		return fmt.Errorf("failed to load lambda functions: %w", err)
	}

	for _, function := range functions {
		if err := p.registerFunction(function); err != nil {
			// Log error but continue loading other functions
			p.app.Logger().Error("Failed to register lambda function", "function", function.GetString("name"), "error", err)
		}
	}

	return nil
}

// registerFunction registers triggers for a specific lambda function
func (p *LambdaFunctionPlugin) registerFunction(function *core.Record) error {
	if !function.GetBool("enabled") {
		return nil
	}

	functionID := function.Id
	triggers := function.GetString("triggers")

	var triggerConfig map[string]interface{}
	if err := json.Unmarshal([]byte(triggers), &triggerConfig); err != nil {
		return fmt.Errorf("invalid trigger configuration: %w", err)
	}

	// Register HTTP triggers
	if httpTriggers, ok := triggerConfig["http"].([]interface{}); ok {
		for _, trigger := range httpTriggers {
			if httpTrigger, ok := trigger.(map[string]interface{}); ok {
				method := strings.ToUpper(httpTrigger["method"].(string))
				path := httpTrigger["path"].(string)
				p.registerHTTPTrigger(functionID, method, path)
			}
		}
	}

	// Register database triggers
	if dbTriggers, ok := triggerConfig["database"].([]interface{}); ok {
		for _, trigger := range dbTriggers {
			if dbTrigger, ok := trigger.(map[string]interface{}); ok {
				collection := dbTrigger["collection"].(string)
				event := dbTrigger["event"].(string)
				p.registerDatabaseTrigger(functionID, collection, event)
			}
		}
	}

	// Register cron triggers
	if cronTriggers, ok := triggerConfig["cron"].([]interface{}); ok {
		for _, trigger := range cronTriggers {
			if cronTrigger, ok := trigger.(map[string]interface{}); ok {
				schedule := cronTrigger["schedule"].(string)
				p.registerCronTrigger(functionID, schedule)
			}
		}
	}

	return nil
}

// registerHTTPTrigger registers an HTTP trigger for an lambda function
func (p *LambdaFunctionPlugin) registerHTTPTrigger(functionID, method, path string) {
	routeKey := fmt.Sprintf("%s:%s", method, path)
	route := &LambdaFunctionHTTPRoute{
		FunctionID: functionID,
		Method:     method,
		Path:       path,
		Handler:    p.createHTTPHandler(functionID),
	}
	p.httpRoutes.Store(routeKey, route)
}

// registerDatabaseTrigger registers a database trigger for an lambda function
func (p *LambdaFunctionPlugin) registerDatabaseTrigger(functionID, collection, event string) {
	trigger := &LambdaFunctionDBTrigger{
		FunctionID: functionID,
		Collection: collection,
		Event:      event,
	}

	key := fmt.Sprintf("%s:%s", collection, event)
	triggers, _ := p.dbTriggers.LoadOrStore(key, []*LambdaFunctionDBTrigger{})
	updatedTriggers := append(triggers.([]*LambdaFunctionDBTrigger), trigger)
	p.dbTriggers.Store(key, updatedTriggers)
}

// registerCronTrigger registers a cron trigger for an lambda function
func (p *LambdaFunctionPlugin) registerCronTrigger(functionID, schedule string) {
	jobID := fmt.Sprintf("lambda_function_%s", functionID)
	job := &LambdaFunctionCronJob{
		FunctionID: functionID,
		Schedule:   schedule,
		JobID:      jobID,
	}

	p.scheduler.MustAdd(jobID, schedule, func() {
		p.executeFunctionForCron(functionID)
	})

	p.cronJobs.Store(functionID, job)
}

// registerHTTPRoutes registers HTTP routes with the PocketBase router
func (p *LambdaFunctionPlugin) registerHTTPRoutes(e *core.ServeEvent) {
	p.httpRoutes.Range(func(key, value interface{}) bool {
		route := value.(*LambdaFunctionHTTPRoute)
		fullPath := "/api/functions" + route.Path
		e.Router.Route(route.Method, fullPath, route.Handler)
		return true
	})
}

// registerDatabaseTriggers registers database event triggers
func (p *LambdaFunctionPlugin) registerDatabaseTriggers() {
	// Register for record creation
	p.app.OnRecordCreate().BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return p.executeFunctionForDBEvent(e.Record, nil, "create")
	})

	// Register for record updates
	p.app.OnRecordUpdate().BindFunc(func(e *core.RecordEvent) error {
		oldRecord := e.Record.Original()
		if err := e.Next(); err != nil {
			return err
		}
		return p.executeFunctionForDBEvent(e.Record, oldRecord, "update")
	})

	// Register for record deletion
	p.app.OnRecordDelete().BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return p.executeFunctionForDBEvent(e.Record, nil, "delete")
	})
}

// createHTTPHandler creates an HTTP handler for an lambda function
func (p *LambdaFunctionPlugin) createHTTPHandler(functionID string) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		ctx := &LambdaFunctionExecutionContext{
			FunctionID:  functionID,
			TriggerType: "http",
			Request:     e.Request,
			Response:    e.Response,
			StartTime:   time.Now(),
		}

		result := p.executeFunction(ctx)
		
		if !result.Success {
			return e.InternalServerError("Lambda function execution failed", fmt.Errorf(result.Error))
		}

		// If the function returned a response object, handle it
		if response, ok := result.Output.(map[string]interface{}); ok {
			if status, ok := response["status"].(float64); ok {
				e.Response.WriteHeader(int(status))
			}
			if headers, ok := response["headers"].(map[string]interface{}); ok {
				for key, value := range headers {
					e.Response.Header().Set(key, fmt.Sprintf("%v", value))
				}
			}
			if body, ok := response["body"]; ok {
				if bodyStr, ok := body.(string); ok {
					return e.String(http.StatusOK, bodyStr)
				}
				return e.JSON(http.StatusOK, body)
			}
		}

		return e.JSON(http.StatusOK, result.Output)
	}
}

// executeFunctionForDBEvent executes functions triggered by database events
func (p *LambdaFunctionPlugin) executeFunctionForDBEvent(record, oldRecord *core.Record, event string) error {
	collection := record.Collection().Name
	key := fmt.Sprintf("%s:%s", collection, event)

	if triggers, ok := p.dbTriggers.Load(key); ok {
		for _, trigger := range triggers.([]*LambdaFunctionDBTrigger) {
			ctx := &LambdaFunctionExecutionContext{
				FunctionID:  trigger.FunctionID,
				TriggerType: "database",
				Record:      record,
				OldRecord:   oldRecord,
				StartTime:   time.Now(),
			}

			// Execute async to not block database operations
			go func(ctx *LambdaFunctionExecutionContext) {
				result := p.executeFunction(ctx)
				if !result.Success {
					p.app.Logger().Error("Lambda function execution failed", 
						"function", ctx.FunctionID, 
						"error", result.Error)
				}
			}(ctx)
		}
	}

	return nil
}

// executeFunctionForCron executes functions triggered by cron
func (p *LambdaFunctionPlugin) executeFunctionForCron(functionID string) {
	ctx := &LambdaFunctionExecutionContext{
		FunctionID:  functionID,
		TriggerType: "cron",
		StartTime:   time.Now(),
	}

	result := p.executeFunction(ctx)
	if !result.Success {
		p.app.Logger().Error("Lambda function cron execution failed", 
			"function", functionID, 
			"error", result.Error)
	}
}

// executeFunction executes an lambda function with the given context
func (p *LambdaFunctionPlugin) executeFunction(ctx *LambdaFunctionExecutionContext) *LambdaFunctionExecutionResult {
	// Load function from database
	function, err := p.app.FindRecordById("lambdas", ctx.FunctionID)
	if err != nil {
		return &LambdaFunctionExecutionResult{
			Success:  false,
			Error:    fmt.Sprintf("Function not found: %v", err),
			Duration: time.Since(ctx.StartTime),
		}
	}

	if !function.GetBool("enabled") {
		return &LambdaFunctionExecutionResult{
			Success:  false,
			Error:    "Function is disabled",
			Duration: time.Since(ctx.StartTime),
		}
	}

	var result *LambdaFunctionExecutionResult

	// Execute with VM from pool
	p.executors.run(func(vm *goja.Runtime) error {
		// Set execution context
		p.setExecutionContext(vm, ctx, function)

		// Execute with timeout
		execCtx, cancel := context.WithTimeout(context.Background(), p.config.MaxExecutionTime)
		defer cancel()

		// Execute the function
		output, err := p.executeWithContext(execCtx, vm, function.GetString("code"))
		
		result = &LambdaFunctionExecutionResult{
			Success:  err == nil,
			Output:   output,
			Error:    p.formatError(err),
			Duration: time.Since(ctx.StartTime),
		}
		
		return nil
	})

	return result
}

// setExecutionContext sets the execution context in the VM
func (p *LambdaFunctionPlugin) setExecutionContext(vm *goja.Runtime, ctx *LambdaFunctionExecutionContext, function *core.Record) {
	// Set environment variables
	env := make(map[string]string)
	if envVars := function.GetString("env_vars"); envVars != "" {
		json.Unmarshal([]byte(envVars), &env)
	}
	vm.Set("$env", env)

	// Set trigger context
	vm.Set("$trigger", map[string]interface{}{
		"type":       ctx.TriggerType,
		"function":   function.GetString("name"),
		"timestamp":  ctx.StartTime.Unix(),
	})

	// Set request context for HTTP triggers
	if ctx.Request != nil {
		vm.Set("$request", map[string]interface{}{
			"method":  ctx.Request.Method,
			"url":     ctx.Request.URL.String(),
			"headers": ctx.Request.Header,
			"body":    p.getRequestBody(ctx.Request),
		})
	}

	// Set record context for database triggers
	if ctx.Record != nil {
		vm.Set("$record", ctx.Record)
		if ctx.OldRecord != nil {
			vm.Set("$oldRecord", ctx.OldRecord)
		}
	}
}

// executeWithContext executes JavaScript code with timeout
func (p *LambdaFunctionPlugin) executeWithContext(ctx context.Context, vm *goja.Runtime, code string) (interface{}, error) {
	done := make(chan struct{})
	var result interface{}
	var err error

	go func() {
		defer close(done)
		result, err = vm.RunString(code)
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("execution timeout")
	case <-done:
		return result, err
	}
}

// getRequestBody extracts request body as string
func (p *LambdaFunctionPlugin) getRequestBody(r *http.Request) string {
	if r.Body == nil {
		return ""
	}
	
	body := make([]byte, 0, 1024)
	if _, err := r.Body.Read(body); err != nil {
		return ""
	}
	
	return string(body)
}

// formatError formats an error for output
func (p *LambdaFunctionPlugin) formatError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// Function lifecycle handlers
func (p *LambdaFunctionPlugin) handleFunctionCreated(record *core.Record) error {
	return p.registerFunction(record)
}

func (p *LambdaFunctionPlugin) handleFunctionUpdated(record *core.Record) error {
	// Remove old registrations
	p.handleFunctionDeleted(record)
	// Register new ones
	return p.registerFunction(record)
}

func (p *LambdaFunctionPlugin) handleFunctionDeleted(record *core.Record) error {
	functionID := record.Id

	// Remove HTTP routes
	p.httpRoutes.Range(func(key, value interface{}) bool {
		route := value.(*LambdaFunctionHTTPRoute)
		if route.FunctionID == functionID {
			p.httpRoutes.Delete(key)
		}
		return true
	})

	// Remove database triggers
	p.dbTriggers.Range(func(key, value interface{}) bool {
		triggers := value.([]*LambdaFunctionDBTrigger)
		filtered := make([]*LambdaFunctionDBTrigger, 0)
		for _, trigger := range triggers {
			if trigger.FunctionID != functionID {
				filtered = append(filtered, trigger)
			}
		}
		if len(filtered) == 0 {
			p.dbTriggers.Delete(key)
		} else {
			p.dbTriggers.Store(key, filtered)
		}
		return true
	})

	// Remove cron jobs
	if job, ok := p.cronJobs.LoadAndDelete(functionID); ok {
		cronJob := job.(*LambdaFunctionCronJob)
		p.scheduler.Remove(cronJob.JobID)
	}

	return nil
}