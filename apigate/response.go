package apigate

import (
	"net/http"
)

// Middleware function to review the API response
func ResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to capture the response
		captureWriter := &responseCaptureWriter{ResponseWriter: w}

		// Call the next middleware or API handler
		next.ServeHTTP(captureWriter, r)

		// Review the response
		responseBody := captureWriter.Body() // Get the captured response body
		LogResponse(r, captureWriter.Status(), responseBody)
		// // Print the response details
		// fmt.Println("Response Status:", captureWriter.Status())
		// fmt.Println("Response Body:", responseBody)
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
