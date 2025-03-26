package httpsrv

import (
	"fmt"
	"net/http"
)

// responseWriter is a wrapper around a http.ResponseWriter which captures the status code and the number of bytes written.
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (lrw *responseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *responseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += n
	return n, err
}

func formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes) // 直接显示字节
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(bytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2fGB", float64(bytes)/(1024*1024*1024))
	}
}
