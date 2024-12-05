package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"spacetime/pkg/utils/db"
	"spacetime/pkg/utils/random"
)

// Returns a random session ID that includes current Unix time in nanoseconds.
func makeSessionID() string {
	// 20 digits (current time) + 1 (:) + 9 (random) = 30 digit session ID
	// 20 digits gives until around 5138 (over 3117 years from now as of writing)
	// assuming Earth's orbit and day remains stable
	// https://www.epochconverter.com/
	return fmt.Sprintf("%020d:%s", time.Now().UnixNano(), random.RandomToken(9))
}

func authUser(w http.ResponseWriter, db db.DBConn, userID uint) error {

	token := makeSessionID()
	expires := time.Now().Add(sessionTokenCookieExpiry)
	_, err := db.Exec(
		`INSERT INTO user_session (token, user_id, expires) VALUES ($1, $2, $3)`,
		token, userID, expires)
	if err != nil {
		return err
	}

	createCookie(w, expires, token)

	return nil

}

func deleteSession(db *sql.DB, token string) error {

	_, err := db.Exec(
		"DELETE FROM user_session WHERE token=$1",
		token,
	)

	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}

	return nil

}

func createCookie(w http.ResponseWriter, expires time.Time, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionTokenCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,                    // don't expose cookie to JavaScript
		SameSite: http.SameSiteStrictMode, // send in first-party contexts only
	})
}

func clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionTokenCookieName,
		Value:    "",
		Path:     "/", // enable AJAX
		Expires:  time.Unix(0, 0),
		HttpOnly: true,                    // don't expose cookie to JavaScript
		SameSite: http.SameSiteStrictMode, // send in first-party contexts only
	})
}

func DeleteExpiredSessions(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM user_session WHERE expires <= $1", time.Now())
	if err != nil {
		return fmt.Errorf("deleting expired sessions: %w", err)
	}
	return nil
}
