package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"treetime/pkg/ajax"
	"treetime/pkg/auth"
	"treetime/pkg/env"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/secrets"
	"treetime/pkg/utils/types"
)

func main() {

	if types.AtoBool(os.Getenv("MAINTENANCE_MODE")) {
		// Site down for databae upgrades
		maintenanceMode()
		os.Exit(0)
	}

	var port = os.Getenv("PORT")
	var dbHost = os.Getenv("DB_HOST")
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var dbName = os.Getenv("DB_NAME")
	var sslMode = ""
	var err error

	env.SetCacheControlVersionStamp(os.Getenv("VERSION_STAMP"))
	env.SetRecaptchaSiteKey(os.Getenv("RECAPTCHA_SITE_KEY"))
	env.SetRecaptchaSecret(os.Getenv("RECAPTCHA_SECRET"))
	env.SetMailjetApiKey(os.Getenv("MAILJET_API_KEY"))
	env.SetMailjetSecret(os.Getenv("MAILJET_SECRET"))

	if env.IsAppEngine() {

		// Database password loaded from secret manager
		secretName := os.Getenv("DB_PASS_SECRET")
		dbPass, err = secrets.LoadSecret(secretName)
		if err != nil {
			logging.LogErrorFatal(err)
		}

	} else {

		// Disable SSL locally
		sslMode = " sslmode=disable"

	}

	// Connect to postgres, using DB_HOST, DB_USER, DB_PASS, DB_NAME
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s"+sslMode,
		dbHost, dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	if env.IsAppEngine() {
		// Set up AppEngine cron handlers
		// r.HandleFunc("/cron/cleanup", makeCronHandler(db, cleanupHandler))
	} else {
		// Set up static resource routes
		// (These static directories are configured by app.yaml for App Engine)
		r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
		r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
		r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	}

	// set up authenticated routes
	authenticate := auth.MakeAuthenticator(db)

	r.PathPrefix("/ajax/").HandlerFunc(authenticate(ajax.AjaxHandler))

	// All other paths go through index handler
	r.PathPrefix("/").HandlerFunc(authenticate(indexHandler))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Write link to /test
		fmt.Fprintf(w, "Welcome to TreeTime!")
	})

	if err = s.ListenAndServe(); err != http.ErrServerClosed {
		logging.LogErrorFatal(err)
	}

	// Flush pending logs
	logging.CloseLoggingClients()

}
