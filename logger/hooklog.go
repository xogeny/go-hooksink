// The package is can be used to build a trivial endpoint
package main

import "log"

import hs "github.com/xogeny/go-hooksink"

// Logger is a trivial handler that conforms to the PushHandler interface
type Logger struct{}

// The Push method makes the Logger struct conform to the PushHandler
// interface.  In this case, it simply logs the PushMessage.
func (l Logger) Push(msg hs.PushMessage, params map[string][]string) {
	log.Printf("PUSH: %v", msg)
}

// The main function creates the new HookSink object, adds a path and
// associates that path with an instance of the Logger handler.  After
// all that, it just starts the server.
func main() {
	h := hs.NewHookSink("ssshhhh!")
	h.Add("/build", Logger{})
	h.Start()
}
