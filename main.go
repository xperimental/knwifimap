package main

import (
	"log"
	"net/http"

	"github.com/spf13/pflag"
)

var addr string
var dbPath string

func main() {
	pflag.StringVarP(&addr, "addr", "a", ":8080", "Address to listen on.")
	pflag.StringVarP(&dbPath, "database", "f", "", "Database file to use.")
	pflag.Parse()

	if len(dbPath) == 0 {
		log.Fatal("Database needs to be set.")
	}

	db, err := newDatabase(dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	http.Handle("/query", QueryHandler(db))
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Printf("Listen on %s...\n", addr)
	http.ListenAndServe(addr, nil)
}
