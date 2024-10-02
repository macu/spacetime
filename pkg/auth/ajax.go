package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"treetime/pkg/user"
	"treetime/pkg/utils/logging"

	"golang.org/x/crypto/bcrypt"
)

func AjaxLogin(db *sql.DB, userID *uint, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if userID != nil {
		// Already authenticated
		return nil, http.StatusForbidden
	}

	if r.Method != http.MethodPost {
		return nil, http.StatusBadRequest
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return nil, http.StatusBadRequest
	}

	var authHash string
	var userRole string

	err := db.QueryRow(
		`SELECT id, auth_hash, user_role FROM user_account WHERE email=$1`,
		email,
	).Scan(&userID, &authHash, &userRole)

	if err != nil {
		if err == sql.ErrNoRows {
			logging.LogNotice(r, struct {
				Event string
				Email string
				// IPAddress string
			}{
				"InvalidLogin",
				email,
				// getUserIP(r),
			})
			return false, http.StatusForbidden
		}
		logging.LogError(r, userID, fmt.Errorf("loading user: %w", err))
		return false, http.StatusInternalServerError
	}

	if user.CheckBanned(userRole) {
		// Don't allow authentication
		logging.LogNotice(r, struct {
			Event string
			Email string
			// IPAddress string
		}{
			"InvalidLogin",
			email,
			// getUserIP(r),
		})
		return false, http.StatusForbidden
	}

	err = bcrypt.CompareHashAndPassword([]byte(authHash), []byte(password))
	if err != nil {
		logging.LogNotice(r, struct {
			Event string
			Email string
			// IPAddress string
		}{
			"InvalidLogin",
			email,
			// getUserIP(r),
		})
		return false, http.StatusForbidden
	}

	err = authUser(w, db, *userID)
	if err != nil {
		logging.LogError(r, userID, fmt.Errorf("authenticating user: %w", err))
		return false, http.StatusInternalServerError
	}

	logging.LogDefault(r, struct {
		Event  string
		UserID uint
		// IPAddress string
	}{
		"UserLogin",
		*userID,
		// getUserIP(r),
	})

	return true, http.StatusOK

}

func AjaxLogout(db *sql.DB, userID uint, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if r.Method != http.MethodPost {
		return nil, http.StatusBadRequest
	}

	sessionTokenCookie, _ := r.Cookie(sessionTokenCookieName)

	_, err := db.Exec(
		"DELETE FROM user_session WHERE token=$1",
		sessionTokenCookie.Value,
	)
	if err != nil {
		logging.LogError(r, &userID, fmt.Errorf("deleting session: %w", err))
		return false, http.StatusInternalServerError
	}

	clearCookie(w)

	return true, http.StatusOK

}