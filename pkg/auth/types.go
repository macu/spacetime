package auth

import (
	"database/sql"
	"net/http"

	"treetime/pkg/utils/ajax"
)

type AuthOptionalHandler func(
	db *sql.DB,
	auth *ajax.Auth,
	w http.ResponseWriter,
	r *http.Request,
)
