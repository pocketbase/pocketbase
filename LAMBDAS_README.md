# Lambda Functions Foundation for PocketBase

This implementation provides the foundation for lambda functions in PocketBase, following the existing PocketBase patterns and conventions.

## Created Files

### Core Models
- `/core/edge_function_model.go` - Main LambdaFunction model with validation and database integration
- `/core/edge_function_query.go` - Query helpers for finding and filtering lambda functions
- `/core/edge_function_context.go` - Execution context and result structures
- `/core/edge_function_events.go` - Event hooks integration

### Migration
- `/migrations/1751650699_edge_functions_init.go` - Database migration to create the `_pb_functions` table and collection

## Features

### LambdaFunction Model
- **Name**: Unique function identifier (3-50 chars, alphanumeric with hyphens/underscores)
- **Code**: JavaScript code to execute
- **Enabled**: Boolean flag to enable/disable function
- **Timeout**: Execution timeout in milliseconds (1s-5min)
- **Triggers**: Array of trigger configurations (HTTP, Database, Cron)
- **EnvVars**: Key-value environment variables
- **Created/Updated**: Automatic timestamps

### Trigger Types
1. **HTTP Triggers**: Route-based execution (GET, POST, etc.)
2. **Database Triggers**: Execute on collection events (insert, update, delete)
3. **Cron Triggers**: Schedule-based execution with cron expressions

### Security
- Functions collection (`_pb_functions`) is system-protected
- Only superusers can manage lambda functions
- Environment variables are isolated per function

### Integration Points
- Follows PocketBase model/collection patterns
- Full event hook integration
- Database query helpers
- Validation and error handling

## Usage Examples

### Creating an Lambda Function
```go
func := &core.LambdaFunction{
    Name:    "hello-world",
    Code:    "export default function(ctx) { return { message: 'Hello World!' }; }",
    Enabled: true,
    Timeout: 30000, // 30 seconds
}

// Add HTTP trigger
httpConfig, _ := json.Marshal(core.HTTPTriggerConfig{
    Method: "GET",
    Path:   "/api/functions/hello",
})

func.Triggers = append(func.Triggers, core.TriggerConfig{
    Type:   core.TriggerTypeHTTP,
    Config: httpConfig,
})

err := app.Save(func)
```

### Finding Functions
```go
// Find by name
fn, err := app.FindLambdaFunctionByName("hello-world")

// Find all enabled functions
enabled, err := app.FindEnabledLambdaFunctions()

// Find functions with database triggers for a collection
dbFunctions, err := app.FindLambdaFunctionsByCollection("posts")
```

## Next Steps

To complete the lambda functions implementation:

1. **Execution Engine**: Implement JavaScript runtime (using goja or similar)
2. **HTTP Router Integration**: Wire HTTP triggers into the existing API router
3. **Database Event Integration**: Hook into collection events to trigger functions
4. **Cron Scheduler**: Integrate with the existing cron system
5. **API Endpoints**: Create CRUD endpoints for managing functions
6. **UI Integration**: Add lambda functions management to the admin dashboard

This foundation provides all the data structures, validation, and database integration needed to build a complete lambda functions system.