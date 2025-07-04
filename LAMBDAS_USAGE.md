# PocketBase Lambda Functions Usage Guide

## Overview

PocketBase Lambda Functions provide a Supabase-like serverless function execution environment that integrates seamlessly with your PocketBase application. Functions are written in JavaScript and can be triggered by HTTP requests, database events, or cron schedules.

## Features

- **JavaScript Runtime**: Powered by Goja engine with Node.js-like APIs
- **Multiple Triggers**: HTTP routes, database events, cron schedules
- **Web Management**: Full admin UI for creating and managing functions
- **Environment Variables**: Per-function environment configuration
- **Execution Logs**: Monitor function performance and debug issues
- **Live Testing**: Execute functions directly from the admin interface

## Getting Started

### 1. Enable Lambda Functions

Update your PocketBase main.go or use the example:

```go
jsvm.MustRegister(app, jsvm.Config{
    HooksDir:      hooksDir,
    HooksWatch:    hooksWatch,
    HooksPoolSize: hooksPool,
    LambdaFunctions: &jsvm.LambdaFunctionPluginConfig{
        PoolSize: 5,
    },
})
```

### 2. Access the Admin Interface

1. Start PocketBase: `./pocketbase serve`
2. Open admin UI: `http://localhost:8090/_/`
3. Log in as a superuser
4. Navigate to "Lambda Functions" in the sidebar

### 3. Create Your First Function

Click "New function" and configure:

- **Name**: `hello-world`
- **Description**: `A simple hello world function`
- **Code**:
```javascript
// Basic function structure
if ($trigger.type === 'http') {
    return {
        status: 200,
        body: {
            message: "Hello from PocketBase Lambda Functions!",
            timestamp: new Date().toISOString(),
            trigger: $trigger
        }
    };
}

console.log("Function executed via", $trigger.type);
```

- **HTTP Trigger**: Add GET `/api/hello`
- **Timeout**: 30 seconds
- **Enabled**: ✓

## Function Structure

### Available Globals

- `$app` - PocketBase app instance
- `$trigger` - Trigger context information
- `$env` - Environment variables
- `$request` - HTTP request object (for HTTP triggers)
- `$record` - Record data (for database triggers)
- `$oldRecord` - Previous record state (for update triggers)

### Function Context

```javascript
// Access trigger information
console.log("Trigger type:", $trigger.type); // "http", "database", or "cron"
console.log("Function name:", $trigger.function);
console.log("Timestamp:", $trigger.timestamp);

// Environment variables
const apiKey = $env.API_KEY;
const dbUrl = $env.DATABASE_URL;

// HTTP requests
if ($trigger.type === 'http') {
    console.log("Method:", $request.method);
    console.log("URL:", $request.url);
    console.log("Headers:", $request.headers);
    console.log("Body:", $request.body);
}

// Database events
if ($trigger.type === 'database') {
    console.log("Record ID:", $record.id);
    console.log("Collection:", $record.collection);
    
    if ($oldRecord) {
        console.log("Previous data:", $oldRecord);
    }
}
```

## Trigger Types

### 1. HTTP Triggers

Execute functions via HTTP requests:

```javascript
// HTTP Response
return {
    status: 200,
    headers: {
        'Content-Type': 'application/json'
    },
    body: {
        message: "Success"
    }
};

// Simple string response
return "Hello World";

// JSON response (auto-detected)
return { message: "Hello World" };
```

**Configuration**:
- Method: GET, POST, PUT, PATCH, DELETE
- Path: `/api/my-endpoint`
- Access: `http://localhost:8090/api/functions/api/my-endpoint`

### 2. Database Triggers

React to record changes:

```javascript
// Create trigger
if ($trigger.type === 'database' && $record) {
    // Send welcome email on user registration
    if ($record.collection === 'users') {
        console.log("New user created:", $record.email);
        
        // Use PocketBase mailer
        // $app.newMailClient().send({
        //     to: $record.email,
        //     subject: "Welcome!",
        //     text: "Welcome to our platform!"
        // });
    }
}

// Update trigger with old/new comparison
if ($oldRecord && $record) {
    if ($oldRecord.status !== $record.status) {
        console.log("Status changed:", $oldRecord.status, "->", $record.status);
    }
}
```

**Configuration**:
- Collection: `users`, `posts`, etc.
- Event: `create`, `update`, `delete`

### 3. Cron Triggers

Schedule periodic execution:

```javascript
// Cron execution
if ($trigger.type === 'cron') {
    console.log("Running scheduled task");
    
    // Cleanup old records
    const oldRecords = $app.findRecordsByFilter(
        'logs', 
        'created < @now(2025, 1, 1, 0, 0, 0)'
    );
    
    oldRecords.forEach(record => {
        $app.delete(record);
    });
    
    console.log("Cleaned up", oldRecords.length, "old records");
}
```

**Configuration**:
- Schedule: `0 2 * * *` (daily at 2 AM)
- Format: Standard cron syntax

## API Access

### PocketBase Database Operations

```javascript
// Find records
const users = $app.findRecordsByFilter('users', 'status = "active"');

// Find single record
const user = $app.findRecordById('users', 'user123');

// Create record
const collection = $app.findCollectionByNameOrId('posts');
const newPost = $app.newRecord(collection);
newPost.set('title', 'New Post');
newPost.set('content', 'Post content');
$app.save(newPost);

// Update record
const post = $app.findRecordById('posts', 'post123');
post.set('title', 'Updated Title');
$app.save(post);

// Delete record
$app.delete(post);
```

### External HTTP Requests

```javascript
// GET request
const response = $http.send({
    url: 'https://api.example.com/data',
    method: 'GET',
    headers: {
        'Authorization': 'Bearer ' + $env.API_TOKEN
    }
});

// POST request
const result = $http.send({
    url: 'https://api.example.com/webhook',
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        event: 'user_created',
        data: $record
    })
});
```

## Testing Functions

### Manual Execution

1. Go to Lambda Functions page
2. Click the "play" button next to your function
3. Provide input JSON: `{"test": true}`
4. Click "Execute Function"
5. View the output and execution time

### Testing HTTP Triggers

```bash
# Test via curl
curl -X GET http://localhost:8090/api/functions/api/hello

# With POST data
curl -X POST http://localhost:8090/api/functions/api/webhook \
  -H "Content-Type: application/json" \
  -d '{"message": "test"}'
```

## Monitoring and Debugging

### Execution Logs

- View real-time execution logs in the admin UI
- Monitor success/failure rates
- Track execution duration
- Debug errors with stack traces

### Error Handling

```javascript
try {
    // Your function logic
    const result = someOperation();
    return { success: true, result };
} catch (error) {
    console.error("Function error:", error);
    return { 
        success: false, 
        error: error.message 
    };
}
```

## Best Practices

### 1. Environment Variables

Store sensitive data in environment variables:

```javascript
// Bad
const apiKey = "sk-1234567890";

// Good
const apiKey = $env.API_KEY;
```

### 2. Error Handling

Always handle errors gracefully:

```javascript
if (!$record) {
    return { error: "No record provided" };
}

try {
    // Your logic here
} catch (error) {
    console.error("Error:", error);
    return { error: "Internal server error" };
}
```

### 3. Timeouts

Keep functions efficient to avoid timeouts:

```javascript
// Use database queries efficiently
const users = $app.findRecordsByFilter(
    'users', 
    'created >= @now(-7days)', // Only recent users
    'created', 
    100 // Limit results
);
```

### 4. HTTP Responses

Return proper HTTP responses:

```javascript
// Success
return {
    status: 200,
    body: { message: "Success" }
};

// Error
return {
    status: 400,
    body: { error: "Invalid input" }
};

// Redirect
return {
    status: 302,
    headers: { Location: "/dashboard" }
};
```

## Security Considerations

1. **Access Control**: Edge functions have superuser access to the database
2. **Input Validation**: Always validate input data
3. **Environment Variables**: Store secrets securely
4. **Error Messages**: Don't expose sensitive information in errors
5. **Rate Limiting**: Consider adding rate limiting for HTTP endpoints

## Examples

### User Welcome Email

```javascript
if ($trigger.type === 'database' && $record.collection === 'users') {
    const mailer = $app.newMailClient();
    
    try {
        mailer.send({
            to: $record.email,
            subject: "Welcome to " + $env.APP_NAME,
            html: `
                <h1>Welcome ${$record.name}!</h1>
                <p>Thanks for joining our platform.</p>
            `
        });
        
        console.log("Welcome email sent to", $record.email);
    } catch (error) {
        console.error("Failed to send welcome email:", error);
    }
}
```

### Webhook Processor

```javascript
if ($trigger.type === 'http' && $request.method === 'POST') {
    const payload = JSON.parse($request.body);
    
    // Validate webhook signature
    const signature = $request.headers['x-signature'];
    const expectedSignature = calculateSignature(payload, $env.WEBHOOK_SECRET);
    
    if (signature !== expectedSignature) {
        return { status: 401, body: { error: "Invalid signature" } };
    }
    
    // Process webhook
    const collection = $app.findCollectionByNameOrId('events');
    const event = $app.newRecord(collection);
    event.set('type', payload.type);
    event.set('data', JSON.stringify(payload));
    $app.save(event);
    
    return { status: 200, body: { message: "Webhook processed" } };
}
```

### Daily Cleanup

```javascript
if ($trigger.type === 'cron') {
    // Clean up old logs (older than 30 days)
    const cutoff = new Date();
    cutoff.setDate(cutoff.getDate() - 30);
    
    const oldLogs = $app.findRecordsByFilter(
        'logs',
        'created < "' + cutoff.toISOString() + '"'
    );
    
    oldLogs.forEach(log => {
        $app.delete(log);
    });
    
    console.log("Cleaned up", oldLogs.length, "old log entries");
    
    return { cleaned: oldLogs.length };
}
```

## Troubleshooting

### Common Issues

1. **Function not executing**
   - Check if function is enabled
   - Verify trigger configuration
   - Check execution logs for errors

2. **Timeout errors**
   - Reduce function complexity
   - Optimize database queries
   - Increase timeout setting

3. **Permission errors**
   - Ensure proper collection access
   - Check if required fields are set

4. **Import errors**
   - Goja doesn't support all Node.js modules
   - Use available PocketBase APIs instead

### Debug Mode

Enable debug logging in your functions:

```javascript
console.log("Debug: Function started");
console.log("Debug: Input data:", JSON.stringify($request));
console.log("Debug: Environment:", Object.keys($env));
```

## Conclusion

PocketBase Lambda Functions provide a powerful way to extend your application with serverless functionality. With the web-based management interface, you can easily create, test, and monitor your functions without deploying separate services.

For more advanced use cases, consider combining multiple trigger types and leveraging PocketBase's built-in features like real-time subscriptions, file storage, and authentication.