package main

import (
	"log"
	// It's good practice to eventually add "os" if you interact with OS, like for pb_public
	// "os" 

	"github.com/pocketbase/pocketbase"
	// "github.com/pocketbase/pocketbase/apis" // Not strictly needed for the minimal example but often used
	"github.com/pocketbase/pocketbase/core" // Needed for core.ServeEvent if you add hooks later
)

func main() {
	// Log a message to confirm app is starting
	log.Println("Starting PocketBase app...")

	app := pocketbase.New()

	// --- Extension code will go here later ---
	// Example:
	// app.OnServe().Add(func(e *core.ServeEvent) error {
	// log.Println("PocketBase server is about to start!")
	// return nil
	// })

	// Start the PocketBase application.
	// The `serve` command is implicitly handled by app.Start() when no other command is specified.
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start PocketBase: %v", err)
	}
}
