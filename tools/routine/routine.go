package routine

import (
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
			if err := recover(); err != nil {
				log.Println("RECOVERED FROM PANIC (safe to ignore):", err)

				stack := make([]byte, 2<<10) // 2 KB
				length := runtime.Stack(stack, false)
				log.Println(string(stack[:length]))
			}
		}()

		f()
	}()
}
