package main

import (
	"database/sql"
	_ "embed"
	"html/template"
	"net/http"

	"treetime/pkg/auth"
	"treetime/pkg/env"
)

var indexTemplate = template.Must(template.ParseFiles("html/index.html"))

func indexHandler(db *sql.DB, userID *uint, w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, struct {
		VersionStamp      string
		PasswordMinLength uint
		Local             bool
	}{
		env.GetCacheControlVersionStamp(),
		auth.PasswordMinLength,
		env.IsLocal(),
	})
}
