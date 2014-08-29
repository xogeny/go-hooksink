package testutils

import "fmt"
import "bytes"
import "net/http"
import "net/http/httptest"
import "crypto/hmac"
import "crypto/sha1"

type HandleFunc func(res http.ResponseWriter, req *http.Request);

func Post(f HandleFunc, path string, payload string) (status int, err error) {
	return AuthPost(f, path, payload, "");
/*
	raw := bytes.NewBuffer([]byte(payload));

	req, err := http.NewRequest("POST", path, raw)
	if (err!=nil) { return; }

	rec := httptest.NewRecorder()	

	f(rec, req);

	status = rec.Code;

	return;
*/
}

func AuthPost(f HandleFunc, path string, payload string, secret string) (status int, err error) {
	raw := bytes.NewBuffer([]byte(payload));

	req, err := http.NewRequest("POST", path, raw)
	if (err!=nil) { return; }

	if (secret!="") {
		mac := hmac.New(sha1.New, []byte(secret))
		mac.Reset()
		mac.Write([]byte([]byte(payload)))
		signature := fmt.Sprintf("sha1=%x", mac.Sum(nil))
		req.Header.Add("X-Hub-Signature", signature)
	}

	rec := httptest.NewRecorder()	

	f(rec, req);

	status = rec.Code;

	return;
}
