package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CheckCreateSpaceThrottleBlock(db *sql.DB, auth ajax.Auth) (bool, error) {

	// Check if 60 or more spaces were created by the user in the last minute
	// If so, return true
	// Otherwise, return false

	var block bool

	var err = db.QueryRow(`SELECT COUNT(*) >= 60 FROM space
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

	// Check if 60 or more spaces were created by the user under the given parent in the last hour
	// If so, return true
	// Otherwise, return false

	block, err := CheckCreateSpaceThrottleBlock(db, auth)

	if err != nil {
		return true, err
	}

	if block {
		return true, nil
	}

	err = db.QueryRow(`SELECT COUNT(*) >= 1 FROM space
		WHERE created_by = $1
		AND parent_id = $2
		AND space_type = $3
		AND created_at > NOW() - INTERVAL '1 MINUTE'`,
		auth.UserID,
		parentID,
		SpaceTypeCheckin,
	).Scan(&block)

	if err != nil {
		return true, fmt.Errorf("throttle create check-in: %w", err)
	}

	return block, nil

}
