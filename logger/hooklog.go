package main

import "log"
import "net/http"

import hs "github.com/xogeny/go-hooksink"

type Logger struct {}

func (l Logger) Push(msg hs.HubMessage) {
	log.Printf("PUSH: %v", msg);
}

func main() {
	h := hs.MakeHookSink();

	h.Config.Addr = "0.0.0.0:3000";

	h.Authenticate(func(req *http.Request) bool {
		return true;
	});

	logger := Logger{};

	h.Add("/log", logger);
	h.Start();
}
