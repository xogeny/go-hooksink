package hooksink

import "log"
import "net/http"
import "encoding/json"
import "github.com/go-martini/martini"
import "github.com/rafecolton/vauth"

type AuthFunction func(req *http.Request) bool;

/* Configuration Options */
type Config struct {
	Addr string;
}

/* HookSink struct definition */
type HookSink struct {
	Config Config;
	martini *martini.ClassicMartini;
	auth AuthFunction;
}

/* Add an authentication handler */
func (hs *HookSink) Authenticate(f AuthFunction) {
	hs.auth = f;
}

/* Method to handle a specific path with a given handler */
func (hs *HookSink) Add(path string, handler interface{}) {
	match := false;

	/* Check if this is a PushHandler */
	pusher, ok := handler.(PushHandler)
	if (ok) {
		match = true;
		// TODO: Need to add logic here to make sure this is actually a `push` event
		// TODO: Test what happens if we end up with multiple handlers at a given path
		hs.martini.Get(path, func(res http.ResponseWriter, h HookSink, req *http.Request) {
			decoder := json.NewDecoder(req.Body)
			msg := HubMessage{};
			err := decoder.Decode(&msg)
			if (err!=nil) {
				log.Printf("Invalid GitHub data: %s", err.Error());
				res.WriteHeader(500);
			} else {
				go pusher.Push(msg);
				res.WriteHeader(200);
			}
		});
	}

	/* If the handler didn't match any known interface, throw a fatal error */
	if (!match) {
		log.Fatalf("Handler didn't match any known handler interface: %v", handler);
	}
}

/* Run the underlying Martini server */
func (hs *HookSink) Start() {
	if (hs.Config.Addr!="") {
		hs.martini.RunOnAddr(hs.Config.Addr);
	} else {
		hs.martini.Run();
	}
}

/* This is used for testing */
func (hs HookSink) Handle(res http.ResponseWriter, req *http.Request) {
	hs.martini.ServeHTTP(res, req);
}


/* This creates a HookSink object.  Much of the work here is in setting up
   the underlying Martini server. */
func MakeHookSink() HookSink {
	ret := HookSink{};

	/* Create a martini instance */
	m := martini.Classic()

	/* Add code to check the HMAC of this request */
	m.Use(vauth.GitHub);

	/* Add some middleware to invoke the Authenticate method (if provided) */
	m.Use(func(res http.ResponseWriter, params martini.Params, req *http.Request) {
		if (ret.auth!=nil) {
			ok := ret.auth(req);
			if (!ok) {
				res.WriteHeader(http.StatusUnauthorized);
			}
		}
	});

	/* Add the martini data to the HookSink object */
	ret.martini = m;

	/* And return */
	return ret;
}
