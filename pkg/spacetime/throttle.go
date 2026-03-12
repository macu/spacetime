package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CheckCreateSpaceThrottleBlock(db *sql.DB, auth ajax.Auth) (bool, error) {

	// Check if 12 or more spaces were created by the user in the last minute
	// If so, return true (throttle block)
	// Otherwise, return false

	var block bool

	var err = db.QueryRow(`SELECT COUNT(*) >= 12 FROM space
		WHERE created_by = $1
		AND created_at > NOW() - INTERVAL '1 MINUTE'`,
		auth.UserID,
	).Scan(&block)

	if err != nil {
		return true, fmt.Errorf("throttle create space: %w", err)
	}

	return block, nil

}

func CheckCreateCheckinThrottleBlock(db *sql.DB, auth ajax.Auth, parentID uint) (bool, error) {

	// Allow 1 check-in per minute per parent space

	var block bool

	err := db.QueryRow(`SELECT COUNT(*) >= 1 FROM checkin
		WHERE user_id = $1
		AND space_id = $2
		AND created_at > NOW() - INTERVAL '1 MINUTE'`,
		auth.UserID,
		parentID,
	).Scan(&block)

	if err != nil {
		return true, fmt.Errorf("throttle create check-in: %w", err)
	}

	return block, nil

}
