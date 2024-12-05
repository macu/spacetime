package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CheckCreateSpaceThrottleBlock(db *sql.DB, auth ajax.Auth) (bool, error) {

	// Check if a space was created by the user within the window
	// If so, return false
	// Otherwise, return true

	var exists bool

	var err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM spaces
		WHERE created_by = ? AND created_at > NOW() - INTERVAL 1 SECOND)`,
		auth.UserID,
	).Scan(&exists)

	if err != nil {
		return true, fmt.Errorf("checkThrottleAllowed: %w", err)
	}

	return exists, nil

}
