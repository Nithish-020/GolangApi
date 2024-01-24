package apigate

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func ResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const RequestIDKey = ContextKey("requestID")
		// Call the next handler in the chain
		requestID := uuid.New().String()

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

		log.Println("requestorDetail, requestID", "", requestorDetail, requestID)

		// Move the logging of request after setting the context

		captureWriter := &ResponseCaptureWriter{ResponseWriter: w}

		// Continue with the request handling
		next.ServeHTTP(captureWriter, r)
		LogResponse(r, captureWriter.Status(), captureWriter.Body(), r.Context().Value(RequestIDKey).(string))

		// Log information about the response after the handler is served
		log.Println("Logging response:", r.Context().Value(RequestIDKey))
		log.Println("r, captureWriter.Status(), responseBody, r.Context().Value(apigate.RequestIDKey).(string)", r, captureWriter.Status(), string(captureWriter.Body()))
	})
}

// Middleware function to review the API response
func ResponseMiddleware0(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the next handler in the chain
		requestID := uuid.New().String()
		log.Println("Start of End Point Call", r.URL.Path, requestID)
		// Create a custom ResponseWriter to capture the response
		captureWriter := &responseCaptureWriter{ResponseWriter: w}

		// Call the next middleware or API handler
		next.ServeHTTP(captureWriter, r)

		// Review the response
		responseBody := captureWriter.Body() // Get the captured response body
		LogResponse(r, captureWriter.Status(), responseBody, "")
		log.Println("End of End Point Call", r.URL.Path, requestID)
	})
}

// Middleware function to review the API response
func ResponseMiddleware2(next http.Handler, getWriteTimeout func(*http.Request) time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to capture the response
		captureWriter := &responseCaptureWriter{ResponseWriter: w}
		// Call the next middleware or API handler
		// 		next.ServeHTTP(captureWriter, r)
		// Call the timeoutMiddleware with the custom write timeout
		timeoutMiddleware(next, getWriteTimeout).ServeHTTP(captureWriter, r)

		responseBody := captureWriter.Body() // Get the captured response body
		LogResponse(r, captureWriter.Status(), responseBody, "")
	})
}

// Custom ResponseWriter to capture the response
type responseCaptureWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *responseCaptureWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseCaptureWriter) Write(body []byte) (int, error) {
	rw.body = append(rw.body, body...)
	return rw.ResponseWriter.Write(body)
}

func (rw *responseCaptureWriter) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

func (rw *responseCaptureWriter) Body() []byte {
	return rw.body
}

// function to specify timeouts for different endpoints,
//
//	enabling flexibility such as setting a specific timeout or no timeout (-1) for particular endpoints.
func CustomWriteTimeout(r *http.Request) time.Duration {
	// Example: Set no timeout for a particular endpoint set -1
	if r.URL.Path == "/sgbEndDateSch" {
		return 15 * time.Second
	}
	// Default write timeout
	return 15 * time.Second
}

func timeoutMiddleware(next http.Handler, getWriteTimeout func(*http.Request) time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the write timeout duration from the custom function
		writeTimeout := getWriteTimeout(r)

		// Check if a specific timeout is set for this endpoint
		if writeTimeout < 0 {
			// No timeout for this endpoint
			next.ServeHTTP(w, r)
			return
		}
		// Create a new context with the specified write timeout
		ctx, cancel := context.WithTimeout(r.Context(), writeTimeout)
		defer cancel()
		r = r.WithContext(ctx)

		// Create a channel to signal when the request handling is finished
		done := make(chan struct{})

		// Start a goroutine to execute the actual HTTP request handling
		go func() {
			//  Based on Done Return Server will close
			defer close(done) // Signal that handling is finished
			next.ServeHTTP(w, r)
		}()

		// Use a select statement to wait for either the request handling to complete or the timeout to occur
		select {
		case <-done:
			// Request completed within the specified timeout
		case <-ctx.Done():
			// request timed out
			w.WriteHeader(http.StatusGatewayTimeout)
			// handle timeout, like logging or sending a specific response
		}
	})
}
