package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"treetime/pkg/env"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/random"
	"treetime/pkg/utils/types"

	"golang.org/x/crypto/bcrypt"
)

func AjaxLoadSignup(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if auth != nil {
		// Already authenticated
		return nil, http.StatusForbidden
	}

	var token = strings.TrimSpace(r.FormValue("token"))

	if token == "" {
		return nil, http.StatusBadRequest
	}

	var email string
	var createdAt time.Time
	err := db.QueryRow(
		`SELECT email, created_at
		FROM user_signup_request
		WHERE token = $1`,
		token,
	).Scan(&email, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return ajax.AjaxErrorPayload{
				ErrorCode: "invalid-token",
			}, http.StatusBadRequest
		}

		logging.LogError(r, nil, fmt.Errorf("error checking signup request: %v", err))
		return nil, http.StatusInternalServerError
	}

	if time.Since(createdAt) > signupTokenExpiry {
		return ajax.AjaxErrorPayload{
			ErrorCode: "token-expired",
		}, http.StatusBadRequest
	}

	return struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Email: email,
		Token: token,
	}, http.StatusOK

}

func AjaxSignup(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if auth != nil {
		// Already authenticated
		return nil, http.StatusForbidden
	}

	var email = strings.ToLower(strings.TrimSpace(r.FormValue("email")))

	if email == "" || !types.ValidateEmailAddress(email) {
		return ajax.AjaxErrorPayload{
			ErrorCode: "invalid-email",
		}, http.StatusBadRequest
	}

	// check if user exists
	var userExists bool
	var err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user_account WHERE email = $1)",
		email,
	).Scan(&userExists)

	if err != nil {
		logging.LogError(r, nil, fmt.Errorf("error checking if user exists: %v", err))
		return nil, http.StatusInternalServerError
	}

	if userExists {
		return ajax.AjaxErrorPayload{
			ErrorCode: "email-exists",
		}, http.StatusBadRequest
	}

	var token = random.RandomToken(signupRequestTokenLength)

	var requestID int64
	err = db.QueryRow(
		`INSERT INTO user_signup_request (email, token, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`,
		email, token, time.Now(),
	).Scan(&requestID)

	if err != nil {
		logging.LogError(r, nil, fmt.Errorf("error inserting signup request: %v", err))
		return nil, http.StatusInternalServerError
	}

	logging.LogNotice(r, struct {
		Event        string
		EmailAddress string
		// IPAddress    string
	}{
		"SignupRequest",
		email,
		// getUserIP(r),
	})

	if env.IsLocal() {
		// Return token for immediate redirect
		return struct {
			Token string `json:"token"`
		}{
			Token: token,
		}, http.StatusOK
	}

	// TODO Send email with verification link

	return true, http.StatusOK

}

func AjaxSignupVerify(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	if auth != nil {
		// Already authenticated
		return nil, http.StatusForbidden
	}

	var token = strings.TrimSpace(r.FormValue("token"))
	var password = r.FormValue("password")
	var handle = strings.TrimSpace(r.FormValue("handle")) // optional
	var displayName = strings.TrimSpace(r.FormValue("displayName"))

	// load user message to send to admin
	// var message = r.FormValue("message")
	// if len(message) > 1000 {
	// 	message = message[:1000]
	// }

	// enforce client-side validation
	if token == "" || displayName == "" ||
		len(handle) > userHandleMaxLength || len(displayName) > userDisplayNameMaxLength ||
		len(strings.TrimSpace(password)) < PasswordMinLength {
		return nil, http.StatusBadRequest
	}

	// load request
	var requestID int64
	var email string
	var createdAt time.Time
	err := db.QueryRow(
		`SELECT id, email, created_at
		FROM user_signup_request
		WHERE token = $1`,
		token,
	).Scan(&requestID, &email, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return ajax.AjaxErrorPayload{
				ErrorCode: "invalid-token",
			}, http.StatusBadRequest
		}

		logging.LogError(r, nil, fmt.Errorf("error checking signup request: %v", err))
		return nil, http.StatusInternalServerError
	}

	if time.Since(createdAt) > signupTokenExpiry {
		return ajax.AjaxErrorPayload{
			ErrorCode: "token-expired",
		}, http.StatusBadRequest
	}

	var handleValue *string
	if handle != "" {
		// verify pattern
		var pattern = regexp.MustCompile(userHandlePattern)
		if !pattern.MatchString(handle) {
			return ajax.AjaxErrorPayload{
				ErrorCode: "invalid-handle",
			}, http.StatusBadRequest
		}

		// verify handle is unused
		var handleExists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM user_account WHERE handle = $1)",
			handle,
		).Scan(&handleExists)

		if err != nil {
			logging.LogError(r, nil, fmt.Errorf("error checking if handle exists: %v", err))
			return nil, http.StatusInternalServerError
		}

		if handleExists {
			return ajax.AjaxErrorPayload{
				ErrorCode: "handle-exists",
			}, http.StatusBadRequest
		}

		handleValue = &handle
	}

	// verify email doesn't yet exist
	var userExists bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user_account WHERE email = $1)",
		email,
	).Scan(&userExists)

	if err != nil {
		logging.LogError(r, nil, fmt.Errorf("error checking if user exists: %v", err))
		return nil, http.StatusInternalServerError
	}

	if userExists {
		return ajax.AjaxErrorPayload{
			ErrorCode: "email-exists",
		}, http.StatusBadRequest
	}

	// hash password
	authHash, err := bcrypt.GenerateFromPassword([]byte(password), passwordBcryptCost)
	if err != nil {
		logging.LogError(r, nil, fmt.Errorf("hashing password: %w", err))
		return nil, http.StatusInternalServerError
	}

	// create user account
	var userID *uint
	err = db.QueryRow(
		`INSERT INTO user_account (email, handle, display_name, auth_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		email, handleValue, displayName, authHash, time.Now(),
	).Scan(&userID)

	if err != nil {
		logging.LogError(r, nil, fmt.Errorf("error creating user account: %v", err))
		return nil, http.StatusInternalServerError
	}

	// delete signup request
	_, err = db.Exec(
		"DELETE FROM user_signup_request WHERE id = $1",
		requestID,
	)

	if err != nil {
		logging.LogError(r, nil, fmt.Errorf("error deleting signup request: %v", err))
		return nil, http.StatusInternalServerError
	}

	logging.LogNotice(r, struct {
		Event        string
		UserID       uint
		EmailAddress string
		Handle       string
		DisplayName  string
	}{
		"SignupVerify",
		*userID,
		email,
		handle,
		displayName,
	})

	// authenticate immediately
	authUser(w, db, *userID)

	// TODO Send email to admin

	return userID, http.StatusOK

}
