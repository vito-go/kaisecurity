package httpsrv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingResponseWriter(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		lrw := responseWriter{ResponseWriter: w, statusCode: 0}
		lrw.WriteHeader(http.StatusCreated)
		n, _ := lrw.Write([]byte("hello world"))
		if lrw.statusCode != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, lrw.statusCode)
		}
		if lrw.bytesWritten != n {
			t.Errorf("expected bytes written %d, got %d", n, lrw.bytesWritten)
		}
	}

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input    int
		expected string
	}{
		{512, "512B"},
		{2048, "2.00KB"},
		{1048576, "1.00MB"},
		{1073741824, "1.00GB"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%d bytes", c.input), func(t *testing.T) {
			result := formatBytes(c.input)
			if result != c.expected {
				t.Errorf("expected %s, got %s", c.expected, result)
			}
		})
	}
}
