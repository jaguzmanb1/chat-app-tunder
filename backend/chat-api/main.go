package main

import (
	"context"
	"go-chat/auth"
	"go-chat/data"
	"go-chat/handlers"
	"go-chat/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var _ = godotenv.Load(".env")

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "rest-api",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Creating logers for each service
	authLogger := l.Named("Auth")
	chatLogger := l.Named("Chat")
	serverLogger := l.Named("Server")

	// Token validator handler
	auth := auth.New(authLogger)

	// Creates message broadcast channel
	m := make(chan data.Message)

	// Chat server to recieve messages
	cs := *server.New(serverLogger, m)

	// Chat http handler
	cha := handlers.New(chatLogger, cs, m)

	// Router creationg
	sm := mux.NewRouter()

	// Injecting authentication validation to mux router
	sm.Use(auth.MiddlewareTokenValidationSocket)

	// Register handlers
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/ws", cha.HandleConnections)

	go cs.HandleMessages()

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         os.Getenv("bindAddress"),                         // configure the bind address
		Handler:      ch(sm),                                           // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Debug("[main] Starting server on", "port", os.Getenv("bindAddress"))

		err := s.ListenAndServeTLS("cert/server.crt", "cert/server.key")
		if err != nil {
			l.Error("[main] Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
