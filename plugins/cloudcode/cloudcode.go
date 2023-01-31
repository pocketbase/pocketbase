package cloudcode

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/cloudcode/api"
	lua "github.com/yuin/gopher-lua"
	"log"
	"os"
)

// MustRegister initializes a Lua environment and panics if anything goes wrong.
// See Register for param information.
func MustRegister(app core.App, cloudCode string, cloudCodeInit string) {
	if err := Register(app, cloudCode, cloudCodeInit); err != nil {
		panic(err)
	}
}

// Register initializes a Lua environment for the application.
// `cloudCode` is the main cloud code file. `cloudCodeInit` (if provided) is run before
// the main cloud code and can be used for things like sandboxing or preloading libraries.
func Register(app core.App, cloudCode string, cloudCodeInit string) error {
	// If the cloud code file isn't available, there's nothing we can do.
	if _, err := os.Stat(cloudCode); os.IsNotExist(err) {
		// TODO: Is there a better way to handle log output?
		log.Println("could not find cloud code")
		return nil
	}

	// If an init script was provided and is inaccessible, exiting early is better (in case the user is relying on the
	// init script for sandboxing or other security measures.
	if cloudCodeInit != "" {
		if _, err := os.Stat(cloudCodeInit); os.IsNotExist(err) {
			// TODO: Is there a better way to handle log output?
			log.Println("could not find cloud code init script")
			return nil
		}
	}

	// Create a new Lua environment and bind all Pocketbase APIs.
	L := lua.NewState()
	api.Bind(&app, L)

	// If we have an init file, run it now.
	if cloudCodeInit != "" {
		err := L.DoFile(cloudCodeInit)
		if err != nil {
			return err
		}
	}

	// Finally, run the user-provided cloud code.
	err := L.DoFile(cloudCode)
	if err != nil {
		return err
	}

	return nil
}
