package testutils

import "log"
import "bytes"
import "net/http"
import "net/http/httptest"

type HandleFunc func(res http.ResponseWriter, req *http.Request);

func Post(f HandleFunc, path string, payload string) (status int, err error) {
	raw := bytes.NewBuffer([]byte(payload));
	log.Printf("POST Body: %v", raw);

	req, err := http.NewRequest("POST", path, raw)
	if (err!=nil) { return; }

	rec := httptest.NewRecorder()	

	f(rec, req);

	status = rec.Code;

	return;
}
