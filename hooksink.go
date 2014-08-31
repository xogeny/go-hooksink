// Package hooksink providers a simple way to create servers that can
// be endpoints for the GitHub webhook API.  It simplifies the process
// of associating code with different paths on the server by handling
// unmarshaling the GitHub payloads into native Go objects and also
// checking the GitHub HMAC signatures prior to invoking the handlers.
package hooksink

import "fmt"
import "log"
import "net/http"
import "io/ioutil"
import "crypto/hmac"
import "crypto/sha1"
import "encoding/json"
import "github.com/go-martini/martini"
import "github.com/martini-contrib/auth"

// Type signature for function to authorize request.  We include the
// entire request to give full access to headers, request path and
// query parameters.
//
// See the Authenticate method on the HookSink object for more
// details.
type AuthFunction func(req *http.Request) bool

// These are configuration options that can be set after the HookSink
// object has been created but before it has been run.
type Config struct {
	// This indicates the address where we should listen for requests
	Addr string
}

// This is the HookSink object.  The only exported field is the Config object.
type HookSink struct {
	Config  Config
	martini *martini.ClassicMartini
	auth    AuthFunction
	secret  string
}

// The Authenticate function allows an authentication function to be
// associated with the endpoint.
//
// We already automatically authenticate the GitHub push (i.e. that it
// came from GitHub using a shared secret).  But this allows another
// potential level of checking if, for example, API keys are part of
// the URL.  However, note that TLS is currently not supported so
// sending of sensitive information is a bad idea until TLS is
// supported for the endpoint.
func (hs *HookSink) Authenticate(f AuthFunction) {
	hs.auth = f
}

// Takes a payload, a secret and a request object and makes sure
// everything is on the up and up.
func checkGitHubSignature(payload []byte, secret string, req *http.Request) bool {
	requestSignature := req.Header.Get("X-Hub-Signature")
	if requestSignature == "" {
		return false
	}

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Reset()
	mac.Write([]byte(payload))
	calculatedSignature := fmt.Sprintf("sha1=%x", mac.Sum(nil))

	return auth.SecureCompare(requestSignature, calculatedSignature)
}

// The Add method allows a given handler to be associated with a
// specified path.
//
// This method uses type assertions to determine what types of
// messages a handler is interested in.  At the moment, only one
// particular kind of handler is supported (see PushHandler
// interface).  If the given handler doesn't conform to any
// expected interface, a fatal error will be logged.
func (hs *HookSink) Add(path string, handler interface{}) {
	match := false

	/* Check if this is a PushHandler */
	pusher, ok := handler.(PushHandler)
	/* If this is a PushHander, setup a Martini handler for it */
	if ok {
		match = true
		// TODO: Need to add logic here to make sure this is actually a `push` event
		//       Use X-GitHub-Event header for this.
		//       see https://developer.github.com/v3/repos/hooks/#webhook-headers
		// TODO: Test what happens if we end up with multiple handlers at a given path
		hs.martini.Post(path, func(res http.ResponseWriter, params martini.Params, req *http.Request) {
			/* Create an empty message */
			foobar := PushMessage{}

			/* Grab the payload from the request body */
			payload, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Printf("Error reading request body: %s", err.Error())
				res.WriteHeader(500)
				return
			}

			/* If this HookSink specified a secret, check the signature */
			if hs.secret != "" {
				if !checkGitHubSignature(payload, hs.secret, req) {
					log.Printf("GitHub signature was not valid")
					res.WriteHeader(401)
					return
				}
			}

			/* Unmarshal the payload into our PushMessage struct */
			err = json.Unmarshal(payload, &foobar)
			if err != nil {
				log.Printf("Error reading JSON data: %s", err.Error())
				res.WriteHeader(500)
				return
			}

			/* Looks like we have a valid PushMessage, let the handler know. */
			go pusher.Push(foobar, map[string][]string(req.URL.Query()))

			/* Reply with "OK" status */
			res.WriteHeader(200)
		})
	}

	/* If the handler didn't match any known interface, throw a fatal error */
	if !match {
		log.Fatalf("Handler didn't match any known handler interface: %v", handler)
	}
}

// The Start method first off the underlying Martini server.
//
// It uses the Config data to determine the address to listen on.  If
// no address is specified, it uses the Martini defaults.
func (hs *HookSink) Start() {
	/* If the user specified a specific address to run on, use that */
	if hs.Config.Addr != "" {
		hs.martini.RunOnAddr(hs.Config.Addr)
	} else {
		hs.martini.Run()
	}
}

// The Handle method is used during testing.  It allows us to avoid
// having to call Start and instead allows us to dispatch a single
// request to the underlying Martini router.  This avoids some of the
// awkwardness and complexity that comes from starting up and shutting
// down the Martini server.
func (hs HookSink) Handle(res http.ResponseWriter, req *http.Request) {
	hs.martini.ServeHTTP(res, req)
}

// This creates a new HookSink object.  Much of the work here is in
// setting up the underlying Martini server.
//
// Note that you should definitely provide a (non-empty) secret here
// (and it should match the secret provided on the GitHub side).  If
// you only provide an empty string (which happens, for example,
// during testing), then no GitHub signature checking can be
// performed.
func NewHookSink(secret string) *HookSink {
	ret := HookSink{
		secret: secret,
	}

	/* Create a martini instance */
	m := martini.Classic()

	/* Add some middleware to invoke the Authenticate method (if provided) */
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if ret.auth != nil {
			ok := ret.auth(req)
			if !ok {
				res.WriteHeader(http.StatusUnauthorized)
			}
		}
	})

	/* Add the martini data to the HookSink object */
	ret.martini = m

	/* And return */
	return &ret
}
