package main

import "log"

import hs "github.com/xogeny/go-hooksink"

type Logger struct {}

func (l Logger) Push(msg hs.HubMessage) {
	log.Printf("PUSH: %v", msg);
}

func main() {
	h := hs.MakeHookSink();
	h.Add("/log", Logger{});
	h.Start();
}
