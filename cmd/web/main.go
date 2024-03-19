package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Aljuagme/book-goweb/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// To use it in other files (same package)
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.BookDB
	templateCache map[string]*template.Template
}

func main() {
	godotenv.Load()

	// Define a new command line flag with the name 'addr', and some explanation.
	// Now you can call go run ./cmd/web -addr=":8888", or by defautl APP_PORT
	// You can use now: go run./cmd/Web -help as well
	addr := flag.String("addr", os.Getenv("APP_PORT"), "HTTP network address")
	seed := flag.Bool("seed", false, "Seed the DB")

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_SERVICE"), os.Getenv("MYSQL_DATABASE"))

	fmt.Println(connStr)

	// you need to parse vefore using it
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := models.OpenDB(connStr)
	if err != nil {
		errorLog.Fatal(err)
	}

	if *seed {
		if err := db.SeedDB(); err != nil {
			errorLog.Fatal(err)
		}
	}

	defer db.DB.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      db,
		templateCache: templateCache,
	}

	// To use our loggers, we create a new struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes
		Handler: app.routes(),
	}

	// prints to console with date and hour
	infoLog.Printf("Starting server on :%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
