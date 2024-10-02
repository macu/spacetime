package auth

import (
	"database/sql"
	"net/http"
)

type AuthOptionalHandler func(
	db *sql.DB,
	userID *uint,
	w http.ResponseWriter,
	r *http.Request,
)
