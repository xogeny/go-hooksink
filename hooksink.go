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
	secret string
}

/* Add an authentication handler */
func (hs *HookSink) Authenticate(f AuthFunction) {
	hs.auth = f;
}

func checkGitHubSignature(payload []byte, secret string, req *http.Request) bool {
	requestSignature := req.Header.Get("X-Hub-Signature")
	if (requestSignature=="") {
		return false;
	}

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Reset()
	mac.Write([]byte(payload))
	calculatedSignature := fmt.Sprintf("sha1=%x", mac.Sum(nil))

	return auth.SecureCompare(requestSignature, calculatedSignature);
}

/* Method to handle a specific path with a given handler */
func (hs *HookSink) Add(path string, handler interface{}) {
	match := false;

	/* Check if this is a PushHandler */
	pusher, ok := handler.(PushHandler)
	if (ok) {
		match = true;
		// TODO: Need to add logic here to make sure this is actually a `push` event
        //       Use X-GitHub-Event header for this.
		//       see https://developer.github.com/v3/repos/hooks/#webhook-headers
		// TODO: Test what happens if we end up with multiple handlers at a given path
		hs.martini.Post(path, func(res http.ResponseWriter, req *http.Request) {
			var foobar HubMessage;
			
			payload, err := ioutil.ReadAll(req.Body);

			if (hs.secret!="") {
				if (!checkGitHubSignature(payload, hs.secret, req)) {
					log.Printf("GitHub signature was not valid");
					res.WriteHeader(401);
					return;
				}
			}

			if (err!=nil) {
				log.Printf("Error reading request body: %s", err.Error());
				res.WriteHeader(500);
				return;
			}

			err = json.Unmarshal(payload, &foobar);
			if (err!=nil) {
				log.Printf("Error reading JSON data: %s", err.Error());
				res.WriteHeader(500);
				return;
			}

			go pusher.Push(foobar);
			res.WriteHeader(200);
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
func NewHookSink(secret string) *HookSink {
	ret := HookSink{
		secret: secret,
	};

	/* Create a martini instance */
	m := martini.Classic()

	/* Add some middleware to invoke the Authenticate method (if provided) */
	m.Use(func(res http.ResponseWriter, req *http.Request) {
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
	return &ret;
}
