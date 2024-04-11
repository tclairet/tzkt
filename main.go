package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	if err := os.MkdirAll("db", os.ModePerm); err != nil {
		panic(err)
	}

	store, err := NewBadger("db")
	if err != nil {
		panic(err)
	}
	port := 3333
	if envPort := os.Getenv("TZKT_PORT"); envPort != "" {
		port, err = strconv.Atoi(envPort)
		if err != nil {
			panic(fmt.Errorf("TZKT_PORT '%s' is not a valid number", envPort))
		}
	}
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: API{store: store}.Routes()}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	p := newPoller(&tzkt{}, 5*time.Second, 0, store)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()

		_ = p.Stop()
	}()

	_ = p.Start()

	logger.Info("server started", "port", port)
	// Run the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
