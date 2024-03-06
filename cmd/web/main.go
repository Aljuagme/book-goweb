package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Rule of thumb: only Panic() or Fatal() in main.

	// Load environment variables
	godotenv.Load()

	// Define a new command line flag with the name 'addr', and some explanation.
	// Now you can call go run ./cmd/web -addr=":8888", or by defautl APP_PORT
	// You can use now: go run./cmd/Web -help as well
	addr := flag.String("addr", os.Getenv("APP_PORT"), "HTTP network address")
	// you need to parse vefore using it
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("snippet/create", snippetCreate)

	// To use our loggers, we create a new struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// prints to console with date and hour
	infoLog.Printf("Starting server on :%s", *addr)

	// err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
