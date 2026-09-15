package routine

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

// FireAndForget executes f() in a new go routine and auto recovers if panic.
//
// **Note:** Use this only if you are not interested in the result of f()
// and don't want to block the parent go routine.
func FireAndForget(f func(), wg ...*sync.WaitGroup) {
	if len(wg) > 0 && wg[0] != nil {
		wg[0].Add(1)
	}

	go func() {
		if len(wg) > 0 && wg[0] != nil {
			defer wg[0].Done()
		}

		defer func() {
			if rec := recover(); rec != nil {
				log.Println("[FireAndForget] RECOVERED FROM PANIC:", rec)

				stack := make([]byte, 2<<10) // 2 KB
				length := runtime.Stack(stack, false)
				log.Println(string(stack[:length]))
			}
		}()

		f()
	}()
}

// SafeWrap wraps the provided function with auto panic recover handling
// and returns any eventual panic as regular error.
func SafeWrap(f func() error) func() error {
	return func() (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := make([]byte, 2<<10) // 2 KB
				length := runtime.Stack(stack, false)
				err = fmt.Errorf("[SafeWrap] recovered from panic: %v\n%s", rec, string(stack[:length]))
			}
		}()

		return f()
	}
}
