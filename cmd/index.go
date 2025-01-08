package main

import (
	"database/sql"
	_ "embed"
	"html/template"
	"net/http"

	"spacetime/pkg/auth"
	"spacetime/pkg/env"
	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
)

var indexTemplate = template.Must(template.ParseFiles("html/index.html"))

func indexHandler(db *sql.DB, user *ajax.Auth, w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, struct {
		Local              bool
		VersionStamp       string
		PasswordMinLength  uint
		TitleMaxLength     uint
		TagMaxLength       uint
		TextMaxLength      uint
		NakedTextMaxDeltas uint
		MaxLimit           uint
	}{
		env.IsLocal(),
		env.GetCacheControlVersionStamp(),
		auth.PasswordMinLength,
		spacetime.TitleMaxLength,
		spacetime.TagMaxLength,
		spacetime.TextMaxLength,
		spacetime.NakedTextMaxDeltas,
		spacetime.MaxSubspacesPageLimit,
	})
}
