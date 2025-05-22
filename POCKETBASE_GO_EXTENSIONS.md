## Introduction

PocketBase is an open-source Backend-as-a-Service (BaaS) solution that simplifies application development by providing a realtime database, user authentication, file storage, and an auto-generated Admin UI, all within a single portable executable. It's designed for speed and ease of use, allowing developers to get a backend up and running in minutes.

While PocketBase offers a comprehensive suite of features out-of-the-box, there are compelling reasons why you might want to extend its core functionality using the Go programming language. PocketBase itself is built with Go, and by using it as a Go framework (sometimes referred to as "PocketBase as a library"), you can seamlessly integrate custom code to:

*   **Implement Complex Business Logic:** Tailor backend processes to your specific needs, such as custom validation rules, data transformations, or intricate workflows that go beyond standard CRUD operations.
*   **Achieve Higher Performance:** For computationally intensive or performance-critical tasks, Go's efficiency, concurrency features, and direct memory control can provide a significant performance boost compared to alternative extension methods like JavaScript hooks.
*   **Leverage the Go Ecosystem:** Gain access to the rich and mature Go ecosystem, including a vast array of libraries and packages for everything from database interactions and external API integrations to specialized algorithms and data processing tools.
*   **Gain Full Control Over the Backend:** Exercise fine-grained control over the request-response lifecycle, application behavior, and integration with other backend services.
*   **Integrate with External Systems:** Build robust integrations with third-party services, payment gateways, messaging queues, or other enterprise systems.

This guide will focus on the primary and most powerful methods for extending PocketBase using Go:

1.  **Event Hooks:** These allow you to subscribe to various core PocketBase events (e.g., before/after record creation, deletion, user authentication, file uploads). By hooking into these events, you can execute custom Go code to modify default behaviors, enforce additional rules, or trigger side effects.
2.  **Custom API Routes:** Define entirely new API endpoints using familiar Go routing patterns. This enables you to build bespoke functionalities, expose new resources, or create specialized APIs that are not covered by PocketBase's default REST-like API.

By mastering these extension points, you can transform PocketBase from a simple BaaS into a highly customized and powerful backend tailored to your application's unique requirements. Let's explore how to get started!

## Getting Started

This section will guide you through setting up a Go project to extend PocketBase.

### Prerequisites

Before you begin, ensure you have the following installed:

*   **Go:** You'll need Go installed on your system. You can find installation instructions on the official Go website: [https://go.dev/doc/install](https://go.dev/doc/install)
*   **PocketBase (as a Go library):** While you might be familiar with the standalone PocketBase executable, here we'll be using PocketBase as a Go library. For more details on its Go package structure and capabilities, refer to the official PocketBase Go overview: [https://pocketbase.io/docs/go-overview/](https://pocketbase.io/docs/go-overview/)

### Basic Project Setup

1.  **Create a Project Directory:**
    First, create a new directory for your project. For example:
    ```bash
    mkdir my-pb-app
    cd my-pb-app
    ```

2.  **Initialize a Go Module:**
    Inside your project directory, initialize a Go module. This will create a `go.mod` file to manage your project's dependencies.
    ```bash
    go mod init my-pb-app
    ```
    (Replace `my-pb-app` with your actual project name or module path if you plan to host it on a repository like GitHub, e.g., `github.com/yourusername/my-pb-app`).

3.  **Create the `main.go` File:**
    Create a file named `main.go` in your project directory. This will be the entry point for your custom PocketBase application.
    ```bash
    touch main.go
    ```

### Core `main.go` Structure

Open your `main.go` file and add the following basic structure:

```go
package main

import (
	"log"
	// "os" // Recommended for pb_public, etc.

	"github.com/pocketbase/pocketbase"
	// "github.com/pocketbase/pocketbase/apis" // Useful for API related constants and helpers
	"github.com/pocketbase/pocketbase/core" // Required for hook event types like core.ServeEvent
)

func main() {
	log.Println("Initializing PocketBase application...")
	app := pocketbase.New()

	// --- Your custom Go extension code will be added here ---
	// For example, to register a hook or a new route:
	// app.OnServe().Bind(func(e *core.ServeEvent) error { // Note: Bind is preferred over Add
	// 	 log.Println("Server is starting! You could add custom routes here.")
	// 	 // Example: e.Router.GET("/api/hello", func(req *core.RequestEvent) error { ... })
	// 	 return e.Next()
	// })

	log.Println("Starting PocketBase server...")
	if err := app.Start(); err != nil {
		log.Fatalf("Error starting PocketBase: %v", err)
	}
}
```

**Explanation:**

*   `package main`: Declares the package as `main`, which is necessary for an executable program.
*   `import (...)`: Lists the Go packages required.
    *   `log`: For basic logging.
    *   `github.com/pocketbase/pocketbase`: The core PocketBase library.
    *   `github.com/pocketbase/pocketbase/core`: Provides access to core event types like `core.ServeEvent`, which is essential for interacting with hooks.
*   `app := pocketbase.New()`: This line creates a new instance of the PocketBase application. This `app` object is what you'll use to register custom event hooks and API routes.
*   The comment `// --- Your custom Go extension code will be added here ---` indicates the primary location where you'll integrate your custom logic.
*   `app.Start()`: This function starts the PocketBase server, including the Admin UI and API.

### Fetching Dependencies

Once you have your `main.go` file set up with the necessary imports, you need to fetch the PocketBase library and its own dependencies. Navigate to your project directory in the terminal and run:

```bash
go mod tidy
```

This command will find the `github.com/pocketbase/pocketbase` import in your `main.go`, download the package, and update your `go.mod` and `go.sum` files.

### Running the Application

To run your customized PocketBase application, use the following command in your project directory:

```bash
go run main.go serve
```

**Explanation:**

*   `go run main.go`: This command compiles and runs your `main.go` application.
*   `serve`: This is an argument passed to your PocketBase application, instructing it to start the web server. You'll see output similar to the standard PocketBase executable, including the Admin UI and API URLs.

Any Go extensions you've added in your `main.go` file will now be active. You can access the Admin UI (usually at `http://127.0.0.1:8090/_/`) and interact with your API, now powered by your extended Go application.

With this basic setup, you are ready to start adding custom event hooks and API routes, which will be covered in the following sections.

## Event Hooks

Event Hooks are a powerful mechanism in PocketBase that allow you to attach your custom Go code to specific lifecycle events within the application. These events can range from server initialization and record operations (create, update, delete) to user authentication and email sending.

**Purpose of Event Hooks:**

*   **Data Validation:** Enforce complex validation rules beyond what's possible with collection schema settings.
*   **Data Modification:** Change data before it's saved or after it's retrieved.
*   **Custom Notifications:** Trigger notifications (e.g., email, SMS, push notifications) based on specific events.
*   **Side Effects:** Perform actions like logging, calling external APIs, or updating other records in response to an event.
*   **Modifying Request/Response:** Intercept and alter incoming requests or outgoing responses.

**General Pattern:**

Most event hooks follow a similar registration pattern using the `Bind` method (or `Add` for some; `Add` is an alias for `Bind`, but `Bind` is more common in newer examples and documentation). You provide a callback function that receives an event-specific struct (e.g., `*core.RecordCreateRequestEvent`).

```go
// General pattern
app.OnEventName().Bind(func(e *core.EventType) error {
    // Your custom logic here
    // You can access event-specific data and methods via 'e'
    // e.g., e.Record, e.HttpContext, e.Mailer, etc.

    // If you want to stop further processing of this event by other bound hooks
    // or the default PocketBase behavior (if applicable for "Before" hooks),
    // you can return an error.
    // return errors.New("something went wrong, stop processing")
    // return apis.NewBadRequestError("Invalid input.", nil) // For client errors

    // To continue the event chain (allowing other hooks and default behavior to run):
    return e.Next() 
    // OR for some older hook types or simpler cases:
    // return nil 
    // It's generally safer and more explicit to use e.Next() when available
    // to ensure the chain continues as expected.
})
```

**Crucially**, if your hook is part of a chain (especially "Before" hooks like `OnRecordBeforeCreateRequest`), and you want the default PocketBase action and any subsequent hooks to execute, your function **must** call `return e.Next()`. If you return `nil` directly without `e.Next()`, or return an error, the chain might be halted.

### Practical Examples

Here are some examples of common event hooks. You would typically place this code within the `main` function of your `main.go` file, after `app := pocketbase.New()` and before `app.Start()`.

Make sure to include necessary imports like `github.com/pocketbase/pocketbase/core` and `github.com/pocketbase/pocketbase/apis` if you're using types like `core.ServeEvent` or functions like `apis.NewBadRequestError`.

#### 1. `OnServe`

*   **Purpose:** Executes code when the web server and all PocketBase app components are initialized, right before the server starts listening for requests. This is a good place to register custom API routes or perform setup tasks that require a fully initialized app instance.
*   **Example:** Log a message indicating the server is ready.

```go
// In your main.go, within the main function, after app := pocketbase.New()

app.OnServe().Bind(func(e *core.ServeEvent) error {
	log.Println("PocketBase server instance configured and about to start!")
	
	// Note: Custom API routes are typically added here using e.Router.
	// We'll cover this in the "Custom API Routes" section. For example:
	// e.Router.GET("/api/custom/hello", func(req *core.RequestEvent) error {
	//	 return req.JSON(http.StatusOK, map[string]string{"message": "Hello!"})
	// })

	return e.Next() // Proceed with starting the server
})
```

#### 2. `OnRecordBeforeCreateRequest` (Example: "articles" collection)

*   **Purpose:** Intercepts an API request to create a new record *before* the record is validated and persisted to the database. Ideal for custom validation, data modification, or defaulting values.
*   **Example:** For an "articles" collection, ensure a 'title' field is not empty and default a 'status' field to "draft" if not provided.

```go
// main.go
// Make sure "articles" collection exists in your PocketBase admin UI.
// Fields used: "title" (text), "status" (text or select)

import (
	"github.com/pocketbase/pocketbase/apis"
	// ... other imports like log, pocketbase, core
)

// ... inside func main()
app.OnRecordBeforeCreateRequest("articles").Bind(func(e *core.RecordCreateRequestEvent) error {
	log.Printf("Attempting to create an article. Request Data: %v", e.RequestData)

	// Example: Validate 'title' field from the incoming request data
	// e.Record is an empty models.Record initialized with the request data.
	// For "BeforeCreate", e.Record.Id is not yet set.
	
	title, _ := e.Record.Data()["title"].(string) // Access data using Data() map or GetString, etc.
	if title == "" {
		// Return a structured error to the client
		return apis.NewBadRequestError("The 'title' field cannot be empty for an article.", nil)
	}

	// Example: Default 'status' to "draft" if not provided in the request
	if status, ok := e.Record.Data()["status"].(string); !ok || status == "" {
		e.Record.Set("status", "draft") // Set the value on the record model
		log.Println("Status not provided for new article, defaulting to 'draft'.")
	}
	
	// You can also modify other fields:
	// e.Record.Set("computed_value", "some value based on " + title)

	return e.Next() // Important to continue the event chain for validation and persistence
})
```
*(Note: `e.Record.Id()` is empty in `OnRecordBeforeCreateRequest` as the ID is generated during persistence. `e.RequestData` holds the raw request payload which you can also inspect.)*

#### 3. `OnRecordAfterCreateRequest` (Example: "articles" collection)

*   **Purpose:** Executes *after* a new record has been successfully validated and created via an API request.
*   **Example:** Log the ID and title of the newly created article.

```go
// main.go
// ... inside func main()

app.OnRecordAfterCreateRequest("articles").Bind(func(e *core.RecordCreateRequestEvent) error {
	// e.Record now contains the persisted record, including its ID and other fields.
	log.Printf("Article created successfully! ID: %s, Title: %s, Status: %s", 
		e.Record.Id(), 
		e.Record.GetString("title"),
		e.Record.GetString("status"))
	
	// Example: Send a notification, update an aggregate, trigger a webhook, etc.
	// sendCustomNotification("New article: " + e.Record.GetString("title"))

	// Be mindful of not making this hook too slow, as it may block the response to the client.
	// For longer or potentially failing tasks, consider using background jobs.
	// (See PocketBase documentation on "Jobs and Background Tasks" via `app.NewCron(...)`)

	return e.Next()
})
```

#### 4. `OnRecordBeforeUpdateRequest` (Example: "articles" collection)

*   **Purpose:** Intercepts an API request to update an existing record *before* the changes are validated and persisted. Useful for conditional updates, field protection, or complex validation logic based on existing and new data.
*   **Example:** Prevent updates to a 'legacy_field' or ensure 'status' can only change according to defined rules.

```go
// main.go
// Fields used: "legacy_field" (text), "status" (text or select)

// ... inside func main()
app.OnRecordBeforeUpdateRequest("articles").Bind(func(e *core.RecordUpdateRequestEvent) error {
	log.Printf("Attempting to update article with ID: %s. Request Data: %v", e.Record.Id(), e.RequestData)

	// e.Record contains the record *with changes from the request already applied* but not yet persisted.
	// To get the original record data for comparison, you can use e.Record.OriginalCopy().
	// originalRecord := e.Record.OriginalCopy()
	// originalStatus := originalRecord.GetString("status")

	// Example: Prevent updates to a 'legacy_field' if it's part of the request payload
	// and its value is actually changing.
	if newLegacyValue, payloadHasLegacyField := e.Record.Data()["legacy_field"]; payloadHasLegacyField {
		if newLegacyValue.(string) != e.Record.OriginalCopy().GetString("legacy_field") {
			// One way to prevent change is to revert it to original
			// e.Record.Set("legacy_field", e.Record.OriginalCopy().GetString("legacy_field"))
			// log.Println("Attempt to change 'legacy_field' was reverted.")

			// Or, more directly, return an error:
			return apis.NewBadRequestError("The 'legacy_field' is protected and cannot be changed.", nil)
		}
	}

	// Example: Conditional update for 'status'
	// newStatus := e.Record.GetString("status")
	// if originalStatus == "published" && newStatus == "draft" {
	// 	 return apis.NewBadRequestError("Cannot change status from 'published' directly back to 'draft'.", nil)
	// }
	
	return e.Next() // Continue with validation and persistence
})
```

#### 5. `OnMailerSend` (General example)

*   **Purpose:** Allows modification of an email message (e.g., subject, body, recipients, headers) *before* it's actually sent by PocketBase. This hook applies to all emails sent by the system (verification, password reset, etc.).
*   **Example:** Add a consistent prefix to the subject of all outgoing emails.

```go
// main.go
// ... inside func main()

app.OnMailerSend().Bind(func(e *core.MailerSendEvent) error {
	log.Printf("Sending email with subject: '%s' to %v", e.Message.Subject, e.Message.To)
	
	e.Message.Subject = "[My PB App] " + e.Message.Subject
	// You can also modify e.Message.HTML, e.Message.Text, e.Message.To, e.Message.From, e.Message.Headers etc.
	// For example, to add a CC:
	// e.Message.Cc = append(e.Message.Cc, mail.Address{Address: "archive@example.com"})

	log.Printf("Modified email subject to: '%s'", e.Message.Subject)
	
	// If you wanted to use a completely different email sending service (e.g., SendGrid, Postmark)
	// instead of the one configured in PocketBase settings (default is SMTPSettings),
	// you could replace e.Mailer with your custom mailer implementation that satisfies the mailer.Mailer interface.
	// e.Mailer = &myCustomMailer.Client{}

	return e.Next() // Proceed with sending the email
})
```

### Other Available Hooks

PocketBase offers a wide range of other event hooks that you can leverage, including (but not limited to):

*   **Authentication Hooks:** `OnUserBeforeAuthenticateRequest`, `OnUserAfterAuthenticateRequest`, `OnUserBeforeVerifyRequest`, `OnUserAfterVerifyRequest`, etc.
*   **File Hooks:** `OnRecordBeforeFileDownloadRequest`, `OnFileBeforeUploadRequest`, `OnFileAfterUploadRequest`.
*   **Settings Hooks:** `OnSettingsBeforeUpdateRequest`, `OnSettingsAfterUpdateRequest`.
*   **Realtime Hooks:** `OnRealtimeConnectRequest`, `OnRealtimeDisconnectRequest`, `OnRealtimeSubscribeRequest`.
*   **CRUD Hooks for specific actions:** `OnRecordBeforeDeleteRequest`, `OnRecordAfterDeleteRequest`, `OnRecordBeforeViewRequest`, `OnRecordAfterViewRequest`, etc.

For a comprehensive list of all available Go event hooks, their specific event argument types, and more detailed explanations, please refer to the official PocketBase documentation:

**[https://pocketbase.io/docs/go-event-hooks/](https://pocketbase.io/docs/go-event-hooks/)**

By understanding and utilizing these hooks, you can significantly extend and customize your PocketBase application's behavior to fit your exact requirements. The next section will cover how to define completely new API functionality using Custom API Routes.

## Custom API Routes

While Event Hooks allow you to modify existing PocketBase behavior, Custom API Routes empower you to define entirely new API endpoints. This is essential when:

*   Your desired functionality doesn't fit the standard Create, Read, Update, Delete (CRUD) model of PocketBase collections.
*   You need to implement complex business processes that involve multiple steps or external service integrations.
*   You want to expose specific data or operations in a controlled manner, separate from the default collection API.

**Registration:**

Custom routes are typically registered within the `OnServe` event hook. This hook provides access to `e.Router`, which is an instance of `echo.Router` (PocketBase uses the [Echo web framework](https://echo.labstack.com/) internally).

The handler function for a custom route has the signature: `func(req *core.RequestEvent) error`.
The `core.RequestEvent` struct provides access to the HTTP request, response writer, PocketBase application instance (`req.App`), and authenticated user/admin model (`req.AuthRecord` / `req.AuthAdmin`).

```go
// In your main.go
// Ensure you have these imports for the examples below:
import (
	"log"
	"net/http"
	// "errors" // If you need to return generic errors

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	// "github.com/pocketbase/pocketbase/models" // For database interaction examples
)

// ... inside func main()
app.OnServe().Bind(func(e *core.ServeEvent) error {
	// Register your custom routes here using e.Router
	// e.Router.GET("/path", handlerFunc, optionalMiddlewares...)
	// e.Router.POST("/path", handlerFunc, optionalMiddlewares...)
	// e.Router.PUT("/path", handlerFunc, optionalMiddlewares...)
	// e.Router.DELETE("/path", handlerFunc, optionalMiddlewares...)

	// --- EXAMPLES WILL GO HERE ---

	return e.Next()
})
```

### Practical Examples

All route definitions below are placed inside the `app.OnServe().Bind(func(e *core.ServeEvent) error { ... })` block in your `main.go`.

#### 1. Simple GET Request

*   **Endpoint:** `/api/custom/hello`
*   **Action:** Returns a JSON message `{"message": "Hello from custom route!"}`.

```go
// In main.go, within app.OnServe().Bind(func(e *core.ServeEvent) error { ... })
e.Router.GET("/api/custom/hello", func(req *core.RequestEvent) error {
	return req.JSON(http.StatusOK, map[string]string{"message": "Hello from custom route!"})
})
```

#### 2. GET Request with Path Parameters

*   **Endpoint:** `/api/custom/greet/{name}`
*   **Action:** Reads the `name` path parameter and returns `{"greeting": "Hello, [name]!"}`.

```go
// In main.go, within app.OnServe().Bind(func(e *core.ServeEvent) error { ... })
e.Router.GET("/api/custom/greet/{name}", func(req *core.RequestEvent) error {
	name := req.Request.PathValue("name") // Echo's way to get path parameters
	return req.JSON(http.StatusOK, map[string]string{"greeting": "Hello, " + name + "!"})
})
```

#### 3. POST Request with JSON Body

*   **Endpoint:** `/api/custom/submit_data`
*   **Action:** Reads a JSON body (e.g., `{"item": "test", "value": 123}`), logs it, and returns a success message.

```go
// Define a struct for the request body (can be defined at package level in main.go or another file)
type SubmitDataPayload struct {
	Item  string `json:"item"`
	Value int    `json:"value"`
}

// In main.go, within app.OnServe().Bind(func(e *core.ServeEvent) error { ... })
e.Router.POST("/api/custom/submit_data", func(req *core.RequestEvent) error {
	data := &SubmitDataPayload{} // Must be a pointer for BindBody
	if err := req.BindBody(data); err != nil {
		// apis.NewBadRequestError provides a structured error response
		return apis.NewBadRequestError("Failed to parse request body.", err)
	}

	log.Printf("Received data via POST: Item = %s, Value = %d", data.Item, data.Value)

	// Example interaction with PocketBase app instance:
	// (Ensure "github.com/pocketbase/pocketbase/models" is imported for this part)
	// appInstance := req.App // Get the PocketBase app instance from the request event
	// itemsCollection, err := appInstance.Dao().FindCollectionByNameOrId("items_collection")
	// if err != nil {
	//     log.Printf("Error finding collection: %v", err) // Log the detailed error
	//     return apis.NewInternalServerError("Could not find the 'items_collection'.", err)
	// }
	// record := models.NewRecord(itemsCollection)
	// record.Set("name", data.Item)
	// record.Set("quantity", data.Value)
	// if err := appInstance.Dao().SaveRecord(record); err != nil {
	//     log.Printf("Error saving record: %v", err) // Log the detailed error
	//     return apis.NewInternalServerError("Failed to save the submitted item.", err)
	// }
	// log.Printf("Successfully saved record with ID: %s", record.Id())
	
	return req.JSON(http.StatusOK, map[string]string{"status": "success", "received_item": data.Item})
})
```

#### 4. Route with Middleware (Built-in and Custom)

*   **Endpoint:** `/api/custom/protected_resource`
*   **Action:** A simple resource that requires authentication and logs a message via custom middleware.

```go
// MyCustomMiddleware is an example of a custom middleware.
// Define this at the package level in main.go or in its own .go file.
func MyCustomMiddleware(next core.RouteHandlerFunc) core.RouteHandlerFunc {
	return func(req *core.RequestEvent) error {
		userId := "anonymous"
		// Safely check if AuthRecord is populated before accessing its ID
		if req.AuthRecord != nil {
			userId = req.AuthRecord.Id()
		}
		log.Println("MyCustomMiddleware: Request received for", req.Request.URL.Path, "from user:", userId)
		
		err := next(req) // Call the next handler in the chain (could be another middleware or the main handler)
		
		log.Println("MyCustomMiddleware: Request processing finished for", req.Request.URL.Path)
		return err
	}
}

// In main.go, within app.OnServe().Bind(func(e *core.ServeEvent) error { ... })
e.Router.GET(
	"/api/custom/protected_resource", 
	func(req *core.RequestEvent) error { // This is the main handler
		// req.AuthRecord is guaranteed to be non-nil here 
		// if apis.RequireAuth() middleware is placed before this handler and succeeds.
		authRecord := req.AuthRecord 
		return req.JSON(http.StatusOK, map[string]string{
			"message":   "You are accessing a protected resource!",
			"userId":    authRecord.Id(),
			"userEmail": authRecord.Email(), // Example of accessing auth record data
		})
	},
	// Middlewares are executed in the order they are provided:
	apis.RequireAuth(), // Built-in PocketBase middleware to ensure the request is authenticated.
	                    // If not authenticated, it returns an error and stops the chain.
	MyCustomMiddleware,   // Your custom middleware (defined above).
)
```

### Key Request/Response Handling in Custom Routes

The `req *core.RequestEvent` object provides several methods and properties for handling requests and responses:

*   **Reading Path Parameters:**
    `name := req.Request.PathValue("name")` (for routes like `/api/custom/items/{name}`)
*   **Reading Query Parameters:**
    `searchTerm := req.Request.URL.Query().Get("search")` (for URLs like `/api/custom/items?search=keyword`)
    `limit := req.Request.URL.Query().Get("limit")`
*   **Reading Request Body:**
    Define a struct matching your expected JSON payload.
    `payload := &YourStructType{}`
    `if err := req.BindBody(payload); err != nil { return apis.NewBadRequestError("Invalid body", err) }`
    *(Ensure `payload` is a pointer. `BindBody` handles JSON by default if `Content-Type` is `application/json`.)*
*   **Sending JSON Response:**
    `return req.JSON(http.StatusOK, data)`
    *(Where `data` can be a `map[string]any`, a struct, or any other serializable type.)*
*   **Sending String Response:**
    `return req.String(http.StatusOK, "Your plain text message")`
*   **Sending HTML Response:**
    `return req.HTML(http.StatusOK, "<h1>Hello World</h1>")`
*   **Returning Errors:** PocketBase provides structured error types in the `apis` package:
    *   `return apis.NewBadRequestError("Missing required field.", data)`
    *   `return apis.NewNotFoundError("Resource not found.", nil)`
    *   `return apis.NewForbiddenError("You are not allowed to perform this action.", nil)`
    *   `return apis.NewUnauthorizedError("Authentication failed.", nil)` (though `apis.RequireAuth` often handles this)
    *   `return apis.NewInternalServerError("An unexpected error occurred.", err)`
    These helpers ensure consistent JSON error responses that the client-side SDKs can understand.
*   **Accessing Authenticated User/Admin:**
    If a route is protected by `apis.RequireAuth()` (or `apis.RequireAdminAuth()`), and authentication is successful:
    `authRecord := req.AuthRecord` will give you the authenticated user's record (`*models.Record`).
    `authAdmin := req.AuthAdmin` will give you the authenticated admin's model (`*models.Admin`).
    If no authentication middleware is used, or if it's optional and no valid token is provided, these will be `nil`.
*   **Accessing PocketBase App Instance:**
    `appInstance := req.App`
    This gives you the `*pocketbase.PocketBase` instance, allowing you to interact with the database and other core functionalities:
    `record, err := appInstance.Dao().FindRecordById("my_collection", "RECORD_ID")`
    `appInstance.Settings().Smtp.Enabled = false` (though changing settings on the fly should be done with care)

### Route Grouping

For better organization, especially when dealing with multiple related routes, you can use route groups. A group allows you to define a common path prefix and apply middlewares to all routes within that group.

```go
// In main.go, within app.OnServe().Bind(func(e *core.ServeEvent) error { ... })

// Create a group for /api/v1/myfeature
myFeatureGroup := e.Router.Group(
    "/api/v1/myfeature",
    // Optional: middlewares applied to ALL routes in this group
    // MyCommonGroupMiddleware, 
)

// This route will be accessible at /api/v1/myfeature/items
myFeatureGroup.GET("/items", func(req *core.RequestEvent) error {
	return req.JSON(http.StatusOK, map[string]string{"feature": "myfeature", "resource": "items"})
})

// This route will be accessible at /api/v1/myfeature/details/{id}
// and also includes the group's middleware (if any) plus its own specific middleware.
myFeatureGroup.GET("/details/{id}", 
    func(req *core.RequestEvent) error {
        id := req.Request.PathValue("id")
        return req.JSON(http.StatusOK, map[string]string{"feature": "myfeature", "detail_id": id})
    }, 
    // Specific middleware for this route only (e.g., require auth for this specific sub-route)
    // apis.RequireAuth(), 
)

// You can also apply middleware to the group after its creation
// myFeatureGroup.Use(AnotherGroupMiddleware)
```

### Further Reading

The routing in PocketBase (via Echo) is very flexible. For more advanced routing patterns, middleware usage, and detailed API for `echo.Router` and `echo.Context` (which `core.RequestEvent` wraps), refer to the official PocketBase and Echo documentation:

*   **PocketBase Go Routing:** [https://pocketbase.io/docs/go-routing/](https://pocketbase.io/docs/go-routing/)
*   **Echo Framework Guide:** [https://echo.labstack.com/guide/](https://echo.labstack.com/guide/)

With custom routes, you have the complete freedom to design your API exactly how your application needs it, complementing PocketBase's built-in features. The next section will discuss how to build and deploy your Go-extended PocketBase application.

## Interacting with the Database

When extending PocketBase with Go, you'll often need to interact with the database. This is typically done within Event Hooks (using `e.App` to get the application instance) or Custom Route handlers (using `req.App`).

The `core.App` interface (which `pocketbase.PocketBase` implements) provides methods for database interaction. While `app.DB()` gives you direct access to the underlying `dbx.DB` instance for raw SQL queries or using the `dbx` query builder, PocketBase's higher-level **Record operations** (like `app.Dao().FindRecordById()`, `app.Dao().SaveRecord()`, etc.) are generally recommended. (Note: `app.Dao()` is the Data Access Object that provides these methods).

**Why prefer Record operations via `app.Dao()`?**

*   **Ease of Use:** They are simpler for common CRUD tasks.
*   **Model Integration:** They work directly with PocketBase's `models.Record` type, providing type safety and convenience for accessing record fields.
*   **Automatic Hook Triggering:** These operations automatically trigger relevant record event hooks (e.g., `OnRecordBeforeCreateRequest`, `OnRecordAfterUpdateOperation`), ensuring your custom logic and PocketBase's internal logic are consistently applied.
*   **Collection Rule Enforcement:** They usually respect collection API rules and permissions, unless specifically bypassed (which is rare and requires careful consideration).

Use `app.DB()` for more complex scenarios not covered by standard record operations, such as complex aggregations not supported by `app.Dao().CountRecords`, bulk updates outside of transactions, or when you need fine-grained control with raw SQL.

### Using Record Operations via `app.Dao()`

Here are examples of common record operations. Assume `app` is an instance of `core.App` (e.g., `app := e.App` in a hook, or `app := req.App` in a route handler). The DAO (Data Access Object) is available via `app.Dao()`.

Remember to include necessary imports:
```go
import (
	"log"
	"database/sql" // For sql.ErrNoRows
	"errors"       // For errors.Is

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/apis" // For structured error responses
	"github.com/pocketbase/dbx"          // For dbx.Params and dbx.NewExp
	// "net/http" // For http.StatusCreated etc. if in a route handler
)
```

#### Fetching Records

##### `app.Dao().FindRecordById("collectionNameOrId", "recordId")`
Retrieves a single record by its ID from the specified collection.

```go
// Make sure 'app' is in scope, e.g., app := req.App or e.App
record, err := app.Dao().FindRecordById("articles", "some_article_id")
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        log.Println("Article not found:", err)
        // If in a route handler, you might return a 404:
        // return apis.NewNotFoundError("Article not found.", err)
    } else {
        log.Println("Error fetching article:", err)
        // If in a route handler, you might return a 500:
        // return apis.NewInternalServerError("Error fetching article.", err)
    }
    // return err // Or handle as appropriate for your context
}
if record != nil { 
    log.Println("Found article:", record.GetString("title"))
}
```

##### `app.Dao().FindFirstRecordByFilter("collectionNameOrId", "filter string", dbx.Params{"key": value})`
Retrieves the first record matching the filter criteria. The filter string uses placeholders like `{:key}` which are replaced by values from `dbx.Params`.

```go
userRecord, err := app.Dao().FindFirstRecordByFilter(
    "users",
    "email = {:email} && status = {:status}", // Filter placeholder
    dbx.Params{"email": "test@example.com", "status": "active"},
)
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        log.Println("Active user with that email not found.")
    } else {
        log.Println("Error finding user:", err)
    }
    // return err
}
if userRecord != nil {
    log.Println("Found user:", userRecord.Id(), "Email:", userRecord.Email())
}
```

##### `app.Dao().FindRecordsByFilter("collectionNameOrId", "filter string", sort, limit, offset, dbx.Params{"key": value})`
Retrieves multiple records matching the filter.
*   `filter`: Similar to `FindFirstRecordByFilter`. Use an empty string `""` to match all records (subject to limit/offset).
*   `sort`: A string defining sort order, e.g., `"-created"` (descending by created date), `"name ASC"`. Use `""` for default/no specific sort.
*   `limit`, `offset`: For pagination.
The result of this function needs to be scanned into a slice using `.All(&yourSlice)`.

```go
var posts []*models.Record // Slice to hold the records
err := app.Dao().FindRecordsByFilter(
    "posts",
    "status = {:status} && category = {:category}", // Filter
    "-created",            // Sort (newest first)
    10,                    // Limit
    0,                     // Offset (start from the beginning)
    dbx.Params{"status": "published", "category": "tech"},
).All(&posts) // Use .All() to scan into the slice

if err != nil {
    log.Println("Error fetching posts:", err)
    // return err
}
log.Printf("Found %d published tech posts.", len(posts))
for _, post := range posts {
    log.Println("Published tech post:", post.GetString("title"))
}
```

##### `app.Dao().CountRecords("collectionNameOrId", filterExpression)`
Counts records matching a filter. The filter is a `dbx.Expression`.

```go
count, err := app.Dao().CountRecords("articles", dbx.NewExp("status = {:status}", dbx.Params{"status": "published"}))
if err != nil {
    log.Println("Error counting articles:", err)
    // return err
}
log.Printf("There are %d published articles.", count)
```

#### Creating Records

To create a new record, you first get the collection, create a new record instance, set its data, and then save it.

```go
// Example: Create a new "comments" record.
// Assumes "comments" collection has a "text" field (text) and "user" field (relation to "users").

collection, err := app.Dao().FindCollectionByNameOrId("comments")
if err != nil {
    log.Println("Error finding collection 'comments':", err)
    // if in a route: return apis.NewNotFoundError("Collection 'comments' not found.", err)
    // return err // Or handle appropriately
}

record := models.NewRecord(collection) // Create a new record for this collection
record.Set("text", "A new comment added via Go functions.")
record.Set("some_other_field", 123.45)

// Example: Set a relation field if you have the related record's ID
// record.Set("user", "USER_RECORD_ID_HERE") // Assuming "user" is a relation field

// If in a route handler and the user is authenticated:
// if req != nil && req.AuthRecord != nil {
//     record.Set("user", req.AuthRecord.Id())
// }


if err := app.Dao().SaveRecord(record); err != nil {
    log.Println("Error saving new comment:", err)
    // if in a route: return apis.NewInternalServerError("Failed to save comment.", err)
    // return err
}

log.Println("Successfully created new comment with ID:", record.Id())
// if in a route, you might return the created record:
// return req.JSON(http.StatusCreated, record)
```

#### Updating Records

To update an existing record, first fetch it, modify its data using `Set()`, and then save it.

```go
recordToUpdate, err := app.Dao().FindRecordById("users", "some_user_id_to_update")
if err != nil {
    log.Println("Error finding user to update:", err)
    // Handle error (e.g., return apis.NewNotFoundError if in a route)
    // return err
}

if recordToUpdate != nil {
    recordToUpdate.Set("bio", "This user's bio has been updated programmatically by Go code.")
    recordToUpdate.Set("isVerified", true) // Example: setting a boolean field

    if err := app.Dao().SaveRecord(recordToUpdate); err != nil {
        log.Println("Error updating user bio:", err)
        // if in a route: return apis.NewInternalServerError("Failed to update user bio.", err)
        // return err
    }
    log.Println("Successfully updated bio for user:", recordToUpdate.Id())
}
```

#### Deleting Records

To delete a record, first fetch it, then pass the record model to `app.Dao().DeleteRecord()`.

```go
recordToDelete, err := app.Dao().FindRecordById("comments", "comment_id_to_delete")
if err != nil {
    log.Println("Error finding comment to delete:", err)
    // Handle error (e.g., return apis.NewNotFoundError if in a route)
    // return err
}

if recordToDelete != nil {
    if err := app.Dao().DeleteRecord(recordToDelete); err != nil {
        log.Println("Error deleting comment:", err)
        // if in a route: return apis.NewInternalServerError("Failed to delete comment.", err)
        // return err
    }
    log.Println("Successfully deleted comment:", recordToDelete.Id())
}
```

#### Custom Record Queries with `app.Dao().RecordQuery()`

For more complex queries not covered by `FindRecordByFilter`, use `app.Dao().RecordQuery()`. This returns a `dbx.SelectQuery` that you can customize with various `WHERE`, `JOIN`, `ORDER BY`, `LIMIT`, etc., conditions.

```go
var articles []*models.Record
err := app.Dao().RecordQuery("articles"). // Specify the collection name or ID
    AndWhere(dbx.HashExp{"status": "active", "type": "featured"}). // Simple key-value conditions
    AndWhere(dbx.NewExp("likes > {:min_likes}", dbx.Params{"min_likes": 100})). // Custom expression
    // Example of joining and filtering on related data (if 'author_details' is a relation)
    // OuterJoin("users", dbx.NewExp("articles.author = users.id")). 
    // AndWhere(dbx.NewExp("users.status = {:userStatus}", dbx.Params{"userStatus": "active"})).
    OrderBy("created DESC"). // Newest first
    Limit(5).                // Get top 5
    All(&articles)           // Execute query and scan into the slice
            
if err != nil {
    log.Println("Error executing custom record query for articles:", err)
    // return err
}

log.Printf("Found %d featured articles with >100 likes:", len(articles))
for _, article := range articles {
    log.Println("Fetched article via RecordQuery:", article.GetString("title"), "Likes:", article.GetInt("likes"))
}
```

### Using `app.DB()` (For Advanced/Raw SQL)

For scenarios where Record operations are insufficient, or if you need to execute raw SQL, you can use `app.DB()`. This returns a `*dbx.DB` instance. For example, this can be useful for complex aggregations, bulk updates not easily handled by record operations, or when specific SQL constructs are needed.

```go
// Example: Get a count of published posts using a raw SQL query
var publishedPostsCount int
err := app.DB().NewQuery("SELECT COUNT(*) FROM posts WHERE status = {:status}").
    Bind(dbx.Params{"status": "published"}).
    Scalar(&publishedPostsCount) // Scalar is for retrieving a single value from the first row

if err != nil {
    log.Println("Error executing raw DB query for post count:", err)
    // return err
}
log.Printf("There are %d published posts (queried via app.DB()).", publishedPostsCount)
```

**Security Note on Raw SQL:** When constructing SQL queries, especially if any part of the query string might come from user input, **always** use parameterized queries (like `dbx.Params` with placeholders `{:key}`) or the `dbx` expression methods. This is crucial to prevent SQL injection vulnerabilities. Avoid string concatenation to build SQL queries with external input.

### Transactions

To perform multiple database operations as a single atomic unit, use `app.RunInTransaction()`. If any operation within the transaction callback returns an error, the entire transaction is rolled back. If `nil` is returned, the transaction is committed.

**Important:** Inside the transaction callback, you must use the `txApp core.App` instance provided to the function for all your database calls (e.g., `txApp.Dao().SaveRecord(...)`). Do not use the global `app` instance for operations you want to be part of the transaction.

```go
// Example: Creating an 'order' and related 'order_items' atomically.

err := app.RunInTransaction(func(txApp core.App) error {
    // txApp is a transactional app instance. Use its Dao for all DB ops here.

    // 1. Create the main order record
    orderCollection, err := txApp.Dao().FindCollectionByNameOrId("orders")
    if err != nil { 
        log.Println("Transaction: Failed to find 'orders' collection.", err)
        return err // Causes rollback
    }
    
    newOrder := models.NewRecord(orderCollection)
    newOrder.Set("customer_name", "Jane Doe Transactional")
    newOrder.Set("status", "pending_transaction")
    newOrder.Set("total_amount", 150.99)

    if err := txApp.Dao().SaveRecord(newOrder); err != nil {
        log.Println("Transaction: Failed to save new order.", err)
        return err // Causes rollback
    }
    log.Printf("Transaction: Order %s created (pending commit).", newOrder.Id())

    // 2. Create related order items
    orderItemsCollection, err := txApp.Dao().FindCollectionByNameOrId("order_items")
    if err != nil { 
        log.Println("Transaction: Failed to find 'order_items' collection.", err)
        return err // Causes rollback
    }

    itemsData := []map[string]any{
        {"product_id": "prod_123_tx", "quantity": 2, "price": 50.00},
        {"product_id": "prod_456_tx", "quantity": 1, "price": 50.99},
    }

    for _, itemMap := range itemsData {
        orderItem := models.NewRecord(orderItemsCollection)
        orderItem.Set("order_relation", newOrder.Id()) // Link to the new order
        orderItem.Set("product_name", itemMap["product_id"]) // Assuming 'product_name' field
        orderItem.Set("quantity", itemMap["quantity"])
        orderItem.Set("unit_price", itemMap["price"])
        if err := txApp.Dao().SaveRecord(orderItem); err != nil {
            log.Printf("Transaction: Failed to save order item for product %v. Error: %v", itemMap["product_id"], err)
            return err // Causes rollback
        }
    }
    
    log.Printf("Transaction: All items for order %s created successfully (pending commit).", newOrder.Id())
    return nil // Returning nil commits the transaction
})

if err != nil {
    log.Println("Order creation transaction failed and was rolled back:", err)
    // If in a route handler, you might return an error to the client:
    // return apis.NewInternalServerError("Failed to process order due to a transaction error.", err)
} else {
    log.Println("Order transaction completed and committed successfully.")
}
```

### Necessary Imports Reminder

When working with the database, ensure you have the relevant imports. Common ones include:

*   `log` (for logging messages)
*   `database/sql` (for `sql.ErrNoRows` when checking if a record was found)
*   `errors` (for `errors.Is` to robustly check error types)
*   `github.com/pocketbase/pocketbase/core` (for `core.App` and event types)
*   `github.com/pocketbase/pocketbase/models` (for `models.Record` and collection types)
*   `github.com/pocketbase/dbx` (for `dbx.Params`, `dbx.NewExp`, `dbx.HashExp` etc.)
*   `github.com/pocketbase/pocketbase/apis` (for returning structured HTTP errors like `apis.NewNotFoundError`)
*   `net/http` (if you are in a route handler and need to return HTTP status codes like `http.StatusCreated`)

### Further Reading

For more in-depth information on database interactions with Go in PocketBase:

*   **Go Database:** [https://pocketbase.io/docs/go-database/](https://pocketbase.io/docs/go-database/) (Covers `dbx` query builder and raw SQL)
*   **Go Record Operations:** [https://pocketbase.io/docs/go-record-operations/](https://pocketbase.io/docs/go-record-operations/) (Detailed look at `app.Dao().FindRecordBy...`, `app.Dao().SaveRecord`, etc.)
*   **Go Collection Operations:** [https://pocketbase.io/docs/go-collection-operations/](https://pocketbase.io/docs/go-collection-operations/) (For managing collections themselves, e.g., `app.Dao().FindCollectionByNameOrId`)

Mastering these database interaction methods is key to building powerful and dynamic extensions for your PocketBase application. The next section will cover building and deploying your Go-extended application.

## Structuring Your Go Extensions

As your PocketBase Go extensions grow in complexity, organizing your code becomes crucial. While simple extensions can reside entirely within the `main.go` file, moving logic into separate Go packages offers significant benefits:

*   **Improved Maintainability:** Smaller, focused packages are easier to understand, modify, and debug.
*   **Better Readability:** A well-organized project structure makes it clear where different functionalities are located.
*   **Easier Testing:** Individual packages and functions can be unit-tested in isolation.

### Example Project Layout

Here's a suggested directory structure. This is just a recommendation and can be adapted to your project's specific needs:

```
your_project_name/
|-- go.mod
|-- go.sum
|-- main.go         // PocketBase app initialization and registration hub
|-- hooks/          // Go package for event hook handlers
|   |-- record_hooks.go
|   |-- mailer_hooks.go // Example for other types of hooks
|-- routes/         // Go package for custom API route handlers
|   |-- user_routes.go  // Example for user-related custom routes
|   |-- product_routes.go // Example for product-related custom routes
|   |-- middleware.go   // Optional: custom middleware used by your routes
|-- corelogic/      // Optional: for business logic shared across hooks/routes
|   |-- validation_service.go
|   |-- notification_service.go
|-- models/         // Optional: for custom Go structs/types used in your extensions
|   |-- custom_types.go 
```

**Explanation of Suggested Packages:**

*   **`hooks/`**: Contains Go files that define and register your event hook handlers. You might further subdivide this (e.g., `hooks/record/` and `hooks/user/`) if you have many hooks.
    *   `record_hooks.go`: For hooks related to record operations (e.g., `OnRecordBeforeCreateRequest`).
    *   `mailer_hooks.go`: For hooks related to mailer events (e.g., `OnMailerSend`).
*   **`routes/`**: Holds Go files that define and register your custom API routes and any associated custom middleware.
    *   `user_routes.go`: For custom API endpoints related to user management or profiles.
    *   `product_routes.go`: For custom API endpoints related to products or catalog.
    *   `middleware.go`: For custom middleware functions that might be applied to your routes or route groups.
*   **`corelogic/` (Optional)**: A place for business logic that is shared between different hooks or routes. This helps keep your hook and route handlers cleaner and focused on the event/request itself.
    *   `validation_service.go`: Functions for complex custom validation logic.
    *   `notification_service.go`: Functions for sending various types of notifications.
*   **`models/` (Optional)**: If your extensions use custom Go structs (e.g., for request/response payloads, or complex data structures), you can define them here.

### Registering Extensions from Packages in `main.go`

Your `main.go` file acts as the central hub for initializing the PocketBase app and registering all your extensions from the different packages.

**Important:** In the import paths below (e.g., `"your_module_name/hooks"`), you **must** replace `your_module_name` with the actual module name of your project. This is the name you defined in your `go.mod` file (e.g., `module github.com/yourusername/your_project_name`).

#### Example `main.go`

```go
package main

import (
	"log"
	// "os" // Uncomment if needed, e.g., for static file serving

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	
	// ========================================================================
	// IMPORTANT: Replace 'your_module_name' with your actual Go module name
	// as defined in your project's go.mod file.
	// For example, if your go.mod says "module myapp", use "myapp/hooks".
	// If it says "module github.com/user/myapp", use "github.com/user/myapp/hooks".
	// ========================================================================
	"your_module_name/hooks"  // Assuming your module is 'your_module_name'
	"your_module_name/routes" // Assuming your module is 'your_module_name'
	// "your_module_name/corelogic" // Uncomment if you create and use this package
)

func main() {
	app := pocketbase.New()

	log.Println("Initializing and registering custom Go extensions...")

	// Register event hooks from the hooks package
	// These functions would take 'app core.App' as an argument.
	hooks.RegisterRecordEventHooks(app)
	// Example: if you had mailer hooks in a separate file or function:
	// hooks.RegisterMailerEventHooks(app) 

	// Register custom API routes from the routes package.
	// This is typically done within the OnServe hook to get access to e.Router.
	app.OnServe().Bind(func(e *core.ServeEvent) error {
		// Pass the router (e.Router) and the app instance (app) 
		// to your route registration functions.
		// The 'app' instance can be useful in route handlers for DB operations.
		routes.RegisterUserApiRoutes(e.Router, app) 
		// routes.RegisterProductApiRoutes(e.Router, app) // Example for other route groups

		// You could also register global middleware from your routes/middleware.go
		// if it's designed to be applied to the root router.
		// For example: e.Router.Use(routes.GlobalRequestLoggerMiddleware(app))

		log.Println("Custom routes and server-related extensions registered.")
		return e.Next()
	})

	// --- Any other application setup code before starting ---
	// For example, if you have a corelogic package that needs initialization:
	// corelogic.InitializeServices(app)


	log.Println("Starting PocketBase server with extensions...")
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start PocketBase: %v", err)
	}
}
```

#### Example `hooks/record_hooks.go`

```go
package hooks // Belongs to the 'hooks' package

import (
	"log"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/apis" // For apis.NewBadRequestError
	// "your_module_name/corelogic" // Example: if you need shared business logic
)

// RegisterRecordEventHooks sets up all record-related event hooks.
// It takes core.App which allows it to register hooks on the app instance.
func RegisterRecordEventHooks(app core.App) {
	app.OnRecordBeforeCreateRequest("articles").Bind(func(e *core.RecordCreateRequestEvent) error {
		log.Printf("[Hooks] Before creating 'articles' record. Current data: %v", e.Record.Data())
		
		// Example: Ensure title is not empty
		if e.Record.GetString("title") == "" {
			return apis.NewBadRequestError("Article title cannot be empty (checked by hook).", nil)
		}
		
		// Example: Default a status if not provided
		if e.Record.GetString("status") == "" {
			e.Record.Set("status", "pending_review_from_hook")
			log.Println("[Hooks] Defaulted status to 'pending_review_from_hook' for new article.")
		}

		// Example of using a hypothetical corelogic service:
		// if corelogic.IsSpam(app, e.Record.GetString("content")) {
		//     return apis.NewBadRequestError("Content appears to be spam.", nil)
		// }

		return e.Next()
	})

	// You can add more hook registrations here for different collections or events:
	// app.OnRecordAfterUpdateRequest("products").Bind(func(e *core.RecordUpdateRequestEvent) error {
	// 	log.Printf("[Hooks] After updating 'products' record: %s", e.Record.Id())
	// 	// Perform some action, like cache invalidation or notification
	// 	return e.Next()
	// })

	log.Println("Successfully registered record event hooks from 'hooks' package.")
}
```

#### Example `routes/user_routes.go`

```go
package routes // Belongs to the 'routes' package

import (
	"log"
	"net/http"
	"github.com/pocketbase/pocketbase/core"
	// "github.com/pocketbase/pocketbase/apis" // For error helpers or auth middleware like apis.RequireAuth()
	// "your_module_name/corelogic" // Example: if you need shared business logic
)

// RegisterUserApiRoutes sets up custom API routes related to users.
// It needs the router instance (core.Router) from ServeEvent 
// and can use the app instance (core.App) for DB operations or other app functionalities.
func RegisterUserApiRoutes(router core.Router, app core.App) {
	
	// Example route: GET /api/custom/users/{id}/profile_strength
	// This route might calculate a "strength" score for a user's profile.
	router.GET("/api/custom/users/{id}/profile_strength", func(req *core.RequestEvent) error {
		userID := req.Request.PathValue("id") // Extract 'id' from URL path
		log.Printf("[Routes] Calculating profile strength for user: %s", userID)
		
		// --- Example: Fetch user record and perform some logic ---
		// user, err := app.Dao().FindRecordById("users", userID)
		// if err != nil {
		// 	 log.Printf("[Routes] User %s not found: %v", userID, err)
		// 	 return apis.NewNotFoundError("User profile not found.", err)
		// }
		//
		// strength := 0
		// if user.GetString("bio") != "" { strength += 30 }
		// if user.Has("avatar") { strength += 20 } // Check if avatar field exists and is set
		// if user.GetBool("isVerified") { strength += 50 }
		//
		// // For more complex logic, consider moving it to a 'corelogic' package:
		// // strength := corelogic.CalculateProfileStrength(app, user)

		// Placeholder response
		return req.JSON(http.StatusOK, map[string]any{
			"userId": userID,
			"profileStrengthScore": 75, // Replace with actual calculated strength
			"message": "Profile strength calculated (example).",
		})
	} /*, apis.RequireAuth() */) // Optionally, add middleware like apis.RequireAuth() here

	// Add more user-specific routes here...
	// Example: POST /api/custom/users/{id}/send_message
	// router.POST("/api/custom/users/{id}/send_message", func(req *core.RequestEvent) error {
	//     // Handler logic for sending a message
	//     // ...
	//     return req.JSON(http.StatusOK, map[string]string{"status": "message_sent"})
	// }, apis.RequireAuth()) // Protect this route

	log.Println("Successfully registered custom user API routes from 'routes' package.")
}
```

### Tips for Good Go Project Structure

*   **Keep `main.go` Lean:** Its primary role should be initialization and wiring up components from other packages. Avoid putting actual business logic directly in `main.go`.
*   **Group by Feature or Domain:** As your application grows, consider organizing packages by feature (e.g., `users`, `posts`, `orders`) rather than just by type (`hooks`, `routes`). Each feature package could then contain its own hooks, routes, and core logic.
*   **Shared Logic in `corelogic` (or `services`, `internal`):** Place reusable business logic (e.g., complex validation rules, notification sending, data transformation algorithms, interactions with external services) in a separate package (or packages). This avoids duplication in your hook and route handlers and makes the logic easier to test and maintain.
*   **Clear Naming:** Use descriptive and consistent names for packages, files, functions, and variables. This greatly improves the readability and understandability of your code.
*   **Testing:** Structuring code into smaller, well-defined packages makes it significantly easier to write unit tests for individual pieces of logic, ensuring your extensions are robust and reliable.

By adopting a structured approach early on, you'll find your Go extensions for PocketBase are more manageable, scalable, and enjoyable to work on as your project evolves. The next section will cover building and deploying your Go-extended PocketBase application.

## Building and Deployment

Once you've developed your custom Go extensions, the final steps are to build your application into a single executable and deploy it.

### Building Your Custom PocketBase Executable

The process of building your Go-extended PocketBase application is straightforward using standard Go tooling.

1.  **Navigate to your project's root directory** (where your `main.go` and `go.mod` files are located) in your terminal.
2.  **Run the `go build` command:**
    ```bash
    go build -o my_custom_pb main.go
    ```
    *   `go build`: This is the standard Go command to compile packages and dependencies.
    *   `-o my_custom_pb`: This flag specifies the output file name for your executable. You can choose any name you prefer (e.g., `my_app_server`, `pocketbase_extended`). If you omit `-o my_custom_pb` and `main.go` (and your current directory contains `main.go`), Go will typically create an executable named after your module or the directory.
    *   `main.go`: This tells the Go compiler that `main.go` is the entry point of your application.

    This command compiles your `main.go` file, along with all the Go packages you've imported (including your custom `hooks/`, `routes/`, etc., packages and the PocketBase framework itself), into a single, self-contained executable file (e.g., `my_custom_pb`). This executable includes everything needed to run your application; there are no external Go dependencies to manage on the server separately.

3.  **Cross-Compilation (Optional):**
    If you are developing on one operating system (e.g., macOS or Windows) but need to deploy to another (e.g., Linux), Go's cross-compilation capabilities are very useful. You can specify the target operating system and architecture using environment variables:
    ```bash
    # Example: Build for Linux (amd64 architecture) from another OS
    GOOS=linux GOARCH=amd64 go build -o my_custom_pb_linux main.go
    ```
    *   `GOOS`: Target operating system (e.g., `linux`, `windows`, `darwin` for macOS).
    *   `GOARCH`: Target architecture (e.g., `amd64`, `arm64`).

    This will produce an executable (e.g., `my_custom_pb_linux`) specifically for the target environment.

### Deployment

Deploying your Go-extended PocketBase application is remarkably simple due to its single-executable nature.

1.  **The Executable is Your Server:**
    The custom executable file generated by the `go build` command (e.g., `my_custom_pb`) *is* your new, complete PocketBase server, now with all your Go extensions built-in.

2.  **Replace Standard PocketBase Executable:**
    This custom executable effectively replaces the standard pre-built PocketBase executable that you might otherwise download from the PocketBase website. You will not use the standard executable if you are deploying a Go-extended version.

3.  **Running Your Custom Executable:**
    Deployment involves:
    *   Copying your custom executable (e.g., `my_custom_pb`) to your server.
    *   Making it executable (e.g., `chmod +x my_custom_pb`).
    *   Running it with the `serve` command, just like the standard PocketBase:
        ```bash
        ./my_custom_pb serve
        ```
        You can also specify other PocketBase command-line flags as needed, for example:
        ```bash
        ./my_custom_pb serve --http="0.0.0.0:80" --automigrate=false
        ```

4.  **Data, Public, and Migrations Directories:**
    Your custom PocketBase executable works with the standard PocketBase data directories in the same way as the original:
    *   `pb_data/`: This directory will be created (if it doesn't exist) when your application starts and will store your application's SQLite database, settings, and uploaded files (unless S3 or another external file storage is configured).
    *   `pb_public/`: If you have static public files you want to serve (e.g., a custom landing page), place them in a `pb_public` directory next to your executable. Your Go extensions can also serve static files if needed, but `pb_public` is the PocketBase default.
    *   `pb_migrations/`: If you are using PocketBase's built-in migration system (managed via the Admin UI or `migrate` commands), your custom executable will use this directory to find and apply database migrations.

    Typically, when deploying, you would upload your custom executable and, if you have existing data, your `pb_data` directory (and `pb_migrations`, `pb_public` if used) to your server.

### Considerations for Production

*   **Process Management:** Since your custom executable is a long-running server process, you'll want to manage it using a process manager like `systemd` (common on Linux), `supervisor`, or run it within a Docker container. This ensures it restarts if it crashes and can be managed as a service.
*   **Reverse Proxy:** It's common practice to run PocketBase (and your custom version) behind a reverse proxy like Nginx or Caddy. This allows you to handle SSL/TLS termination, custom domain names, rate limiting, caching, and serve multiple applications from the same server.
*   **Security and Backups:** Standard server security practices apply. Ensure your server is hardened, and regularly back up your `pb_data` directory.
*   **PocketBase Production Guide:** For more comprehensive advice on deploying PocketBase instances (which largely still applies to your custom Go-extended version), refer to the official PocketBase documentation:
    [**Going to Production - PocketBase Docs**](https://pocketbase.io/docs/going-to-production/)

Building and deploying your Go-extended PocketBase application leverages Go's strengths in creating portable, efficient executables, making the process relatively straightforward. The next section will provide a summary and point to further resources.

## Conclusion and Further Resources

### Conclusion

Extending PocketBase with Go unlocks a significant level of power and flexibility, transforming it from a simple backend-in-a-box into a highly customizable platform. Throughout this guide, we've explored how to set up a Go development environment for PocketBase and leverage its core extension mechanisms:

*   **Event Hooks:** Allowing you to tap into PocketBase's lifecycle to inject custom logic, modify behavior, and integrate side effects.
*   **Custom API Routes:** Enabling you to build entirely new API functionalities tailored to your specific needs.
*   **Database Interactions:** Providing the tools to query and manipulate your data directly from your Go code, either through high-level Record operations or direct database access.
*   **Project Structure:** Guiding you on how to organize your growing codebase for better maintainability and scalability.

By using Go, you can implement sophisticated business logic, achieve higher performance for critical tasks, and integrate deeply with external systems, all while benefiting from Go's strong typing, rich standard library, and vibrant ecosystem.

We encourage you to experiment with the concepts and examples presented. Start small, incrementally add features, and explore the vast possibilities that Go extensions open up. Whether you're building a complex SaaS product, a specialized internal tool, or a high-performance API, extending PocketBase with Go provides a robust foundation to meet your unique requirements.

### Further Resources

To continue your journey and delve deeper into specific topics, the following resources are highly recommended:

*   **Official PocketBase Go Documentation:** This is your primary reference for all aspects of extending PocketBase with Go. It includes detailed API references, more examples, and advanced usage patterns.
    *   [PocketBase Go Overview](https://pocketbase.io/docs/go-overview/)
    *   [Go Event Hooks](https://pocketbase.io/docs/go-event-hooks/)
    *   [Go Routing](https://pocketbase.io/docs/go-routing/)
    *   [Go Database & Record Operations](https://pocketbase.io/docs/go-database/)
    *   [Go Collection Operations](https://pocketbase.io/docs/go-collection-operations/)

*   **PocketBase GitHub Repository:** The official source code for PocketBase. Useful for understanding internals, finding examples, and tracking development.
    *   [https://github.com/pocketbase/pocketbase](https://github.com/pocketbase/pocketbase)

*   **PocketBase Discussions:** A great place to ask questions, share your projects, and connect with the PocketBase community.
    *   [https://github.com/pocketbase/pocketbase/discussions](https://github.com/pocketbase/pocketbase/discussions)

*   **Go Language Documentation:** If you are new to Go or want to brush up on its fundamentals, the official Go documentation is an excellent resource.
    *   [https://go.dev/doc/](https://go.dev/doc/)

Happy coding, and we look forward to seeing what you build with your Go-extended PocketBase applications!

[end of POCKETBASE_GO_EXTENSIONS.md]
