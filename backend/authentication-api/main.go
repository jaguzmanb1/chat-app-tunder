package main

import (
	"authentication-api/data"
	"authentication-api/handlers"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"

	gohandlers "github.com/gorilla/handlers"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

var _ = godotenv.Load(".env")
var (
	//ConnectionString cadena de conexión a la base de datos
	ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("user"),
		os.Getenv("pass"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("db_name"))
)

func main() {
	// Database object
	db, err := sql.Open("mysql", ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	l := hclog.New(&hclog.LoggerOptions{
		Name:  "authentication-api",
		Level: hclog.LevelFromString("DEBUG")})

	// Creating logger for each service
	dl := l.Named("Data")
	al := l.Named("Auth")

	// Validator object
	v := data.NewValidation()

	// Se crea servicio de usuario
	us := data.New(db, dl)

	// Creating mux server to save handlers
	sm := mux.NewRouter()

	// Creating auth handler
	ah := handlers.New(al, us, v)

	// Subrouter to hanlde post requests
	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/createuser", ah.Signup)
	postR.Use(ah.MiddlewareValidateUser)

	postSignR := sm.Methods(http.MethodPost).Subrouter()
	postSignR.HandleFunc("/signin", ah.Signin)
	postSignR.Use(ah.MiddlewareValidateUserSignin)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// New custom server
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
			l.Error("[main] Error starting server %s\n", err)
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
