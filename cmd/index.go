package main

import (
	"database/sql"
	_ "embed"
	"html/template"
	"net/http"

	"treetime/pkg/auth"
	"treetime/pkg/env"
	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
)

var indexTemplate = template.Must(template.ParseFiles("html/index.html"))

func indexHandler(db *sql.DB, user *ajax.Auth, w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, struct {
		Local                        bool
		VersionStamp                 string
		PasswordMinLength            uint
		TreeMaxDepth                 uint
		CategoryTitleMaxLength       uint
		CategoryDescriptionMaxLength uint
		TagTitleMaxLength            uint
		TypeTitleMaxLength           uint
		TypeDescriptionMaxLength     uint
		PostTitleMaxLength           uint
		PostBlockMaxLength           uint
		PostBlockCount               uint
		PostURLMaxLength             uint
		CommentMaxLength             uint
	}{
		env.IsLocal(),
		env.GetCacheControlVersionStamp(),
		auth.PasswordMinLength,
		treetime.TreeMaxDepth,
		treetime.CategoryTitleMaxLength,
		treetime.CategoryDescriptionMaxLength,
		treetime.TagTitleMaxLength,
		treetime.TypeTitleMaxLength,
		treetime.TypeDescriptionMaxLength,
		treetime.PostTitleMaxLength,
		treetime.PostBlockMaxLength,
		treetime.PostBlockCount,
		treetime.PostURLMaxLength,
		treetime.CommentMaxLength,
	})
}
