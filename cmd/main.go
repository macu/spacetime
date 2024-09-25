package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	// Connect to postgres, using DB_HOST, DB_USER, DB_PASS, DB_NAME
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	db.Close() // not currently used

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Write link to /test
		fmt.Fprintf(w, "Welcome to TreeTime!")
	})

	http.ListenAndServe(":8080", nil)

}
