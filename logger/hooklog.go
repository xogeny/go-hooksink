package main

import "log"

import hs "github.com/xogeny/go-hooksink"

type Logger struct{}

func (l Logger) Push(msg hs.HubMessage, params map[string][]string) {
	log.Printf("PUSH: %v", msg)
}

func main() {
	h := hs.NewHookSink("ssshhhh!")
	h.Add("/build", Logger{})
	h.Start()
}
