package apigate

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		log.Println(err)
		requestorDetail := GetRequestorDetail(r)
		requestorDetail.Body = string(body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		LogRequest("", requestorDetail, "")
		next.ServeHTTP(w, r)

	})
}
