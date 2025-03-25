package httpsrv_test

import (
	"context"
	"errors"
	"github.com/vito-go/kaisecurity/internal/httpsrv"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	srv, err := httpsrv.NewServer(":memory:")
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}
func TestStartServer_Shutdown(t *testing.T) {
	srv, _ := httpsrv.NewServer(":memory:")
	go func() {
		err := srv.StartServer(0)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("start error: %v", err)
		}
	}()
	time.Sleep(200 * time.Millisecond) // wait for server to start
	err := srv.ShutDownServer(context.Background())
	if err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
