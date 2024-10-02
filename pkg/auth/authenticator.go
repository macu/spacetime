package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"treetime/pkg/user"
	"treetime/pkg/utils/logging"
)

// Returns a function that wraps a handler in an authentication intercept that loads
// the authenticated user ID and occasionally updates the expiry of the session cookie.
// The wrapped handler is not called and 401 is returned if no user is authenticated.
func MakeAuthenticator(db *sql.DB) func(handler AuthOptionalHandler) func(http.ResponseWriter, *http.Request) {

	selectUserStmt, err := db.Prepare(
		`SELECT user_session.user_id, user_session.expires, user_account.user_role
		FROM user_session
		INNER JOIN user_account ON user_session.user_id=user_account.id
		WHERE user_session.token=$1 AND user_session.expires>$2`,
	)
	if err != nil {
		panic(err)
	}

	// Return factory function for wrapping handlers that require authentication
	return func(handler AuthOptionalHandler) func(http.ResponseWriter, *http.Request) {

		// Return standard http.Handler which calls the authenticated handler passing db and userID
		return func(w http.ResponseWriter, r *http.Request) {

			var userID *uint
			var userRole string

			// Read auth cookie
			sessionTokenCookie, err := r.Cookie(sessionTokenCookieName)

			if err == nil {

				// Look up session and read authenticated userID
				now := time.Now()
				var expires time.Time

				err = selectUserStmt.QueryRow(
					sessionTokenCookie.Value, now,
				).Scan(&userID, &expires, &userRole)

				if err == sql.ErrNoRows {
					userID = nil
				} else if err != nil {
					logging.LogError(r, nil, fmt.Errorf("loading user from session token: %w", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if userID != nil {
					if user.CheckBanned(userRole) {

						// Don't allow authentication
						userID = nil

						// Delete cookie
						clearCookie(w)

					} else if expires.Before(now.Add(sessionTokenCookieRenewIfExpiresIn)) {
						// Refresh session and cookie if old

						// Update session expires time
						expires := now.Add(sessionTokenCookieExpiry)
						_, err = db.Exec(
							`UPDATE user_session SET expires=$1 WHERE token=$2`,
							expires, sessionTokenCookie.Value)
						if err != nil {
							logging.LogError(r, userID, fmt.Errorf("updating session expiry: %w", err))
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						// Update cookie expires time
						createCookie(w, expires, sessionTokenCookie.Value)

					}
				}

			}

			// Invoke route with authenticated user info
			handler(db, userID, w, r)
		}
	}
}
