package apigate

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Custom ResponseWriter to capture the response
type ContextKey string

type ResponseCaptureWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *ResponseCaptureWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseCaptureWriter) Write(body []byte) (int, error) {
	rw.body = append(rw.body, body...)
	return rw.ResponseWriter.Write(body)
}

func (rw *ResponseCaptureWriter) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

func (rw *ResponseCaptureWriter) Body() []byte {
	return rw.body
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const RequestIDKey = ContextKey("requestID")
		// Call the next handler in the chain
		requestID := uuid.New().String()
		log.Println("START of End Point Call", r.URL.Path, requestID)

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		requestorDetail := GetRequestorDetail(r)
		requestorDetail.Body = string(body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		LogRequest("", requestorDetail, requestID)

		// log.Println("requestorDetail, requestID", "", requestorDetail, requestID)

		// Move the logging of request after setting the context

		captureWriter := &ResponseCaptureWriter{ResponseWriter: w}

		// Continue with the request handling
		next.ServeHTTP(captureWriter, r)
		LogResponse(r, captureWriter.Status(), captureWriter.Body(), r.Context().Value(RequestIDKey).(string))
		log.Println("END of End Point Call", r.URL.Path, requestID)

		// Log information about the response after the handler is served
		// log.Println("Logging response:", r.Context().Value(RequestIDKey))
		// log.Println("r, captureWriter.Status(), responseBody, r.Context().Value(apigate.RequestIDKey).(string)", r, captureWriter.Status(), string(captureWriter.Body()))
	})
}
