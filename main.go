package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/vito-go/kaisecurity/internal/app"
	"github.com/vito-go/mylog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultDBPath := filepath.Join(homeDir, "kai_security.db")
	dbPath := flag.String("db", defaultDBPath, "set db path")
	port := flag.Uint("port", 8080, "set port")
	flag.Parse()
	if *dbPath == "" {
		panic("db path is required")
	}

	appContext, err := app.NewAppContext(*dbPath)
	if err != nil {
		panic(err)
	}
	var exitSignal = make(chan os.Signal, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /scan", appContext.HandleScan)
	// /query
	mux.HandleFunc("POST /query", appContext.HandleQuery)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	mylog.Printf("Server is running on :%d", *port)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-exitSignal
	mylog.Printf("exit signal: %s", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	mylog.Printf("Server is shutting down gracefully")
	err = server.Shutdown(ctx)
	if err != nil {
		mylog.Ctx(ctx).Error(err)
	}
	mylog.Printf("Server is shutdown")
}
