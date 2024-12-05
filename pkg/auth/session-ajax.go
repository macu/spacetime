package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"spacetime/pkg/user"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
)

func AjaxLogin(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if auth != nil {
		// Already authenticated
		return nil, http.StatusForbidden
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return nil, http.StatusBadRequest
	}

	var authHash string
	var userID *uint
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
		logging.LogError(r, auth, fmt.Errorf("loading user: %w", err))
		return false, http.StatusInternalServerError
	}

	if !user.CheckRoleActive(userRole) {
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
		logging.LogError(r, auth, fmt.Errorf("authenticating user: %w", err))
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

func AjaxLoadLogin(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if auth == nil {
		return nil, http.StatusOK
	}

	var user = struct {
		IsAuthenticated bool   `json:"isAuthenticated"`
		ID              uint   `json:"id"`
		Role            string `json:"role"`
		Handle          string `json:"handle"`
		DisplayName     string `json:"displayName"`
		Email           string `json:"email"`
	}{}

	err := db.QueryRow(
		`SELECT id, user_role, handle, display_name, email FROM user_account WHERE id=$1`,
		auth.UserID,
	).Scan(&user.ID, &user.Role, &user.Handle, &user.DisplayName, &user.Email)

	if err != nil {
		logging.LogError(r, auth, fmt.Errorf("loading user: %w", err))
		return nil, http.StatusInternalServerError
	}

	user.IsAuthenticated = true

	return user, http.StatusOK

}

func AjaxLogout(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if r.Method != http.MethodPost {
		return nil, http.StatusBadRequest
	}

	sessionTokenCookie, _ := r.Cookie(sessionTokenCookieName)

	_ = deleteSession(db, sessionTokenCookie.Value) // ignore error

	clearCookie(w)

	return true, http.StatusOK

}
