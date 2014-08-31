package testutils

import "fmt"
import "bytes"
import "net/http"
import "net/http/httptest"
import "crypto/hmac"
import "crypto/sha1"

type HandleFunc func(res http.ResponseWriter, req *http.Request)

/* This function POSTs the given payload to the server (directly, without having
   to start up an actual web server) */
func Post(f HandleFunc, path string, payload string) (status int, err error) {
	/* Just call AuthPost with no secret. */
	return AuthPost(f, path, payload, "")
}

/* This function POSTs the given payload to the server (directly,
   without having to start up an actual web server) and provides a
   secret for HMAC checking */
func AuthPost(f HandleFunc, path string, payload string, secret string) (status int, err error) {
	/* Create a buffer for the payload */
	raw := bytes.NewBuffer([]byte(payload))

	/* Create a new request */
	req, err := http.NewRequest("POST", path, raw)
	if err != nil {
		return
	}

	/* If a secret was given, add an X-Hub-Signature header */
	if secret != "" {
		mac := hmac.New(sha1.New, []byte(secret))
		mac.Reset()
		mac.Write([]byte([]byte(payload)))
		signature := fmt.Sprintf("sha1=%x", mac.Sum(nil))
		req.Header.Add("X-Hub-Signature", signature)
	}

	/* Get ready to record the response */
	rec := httptest.NewRecorder()

	/* Submit the request */
	f(rec, req)

	/* Record the status code */
	status = rec.Code

	return
}
