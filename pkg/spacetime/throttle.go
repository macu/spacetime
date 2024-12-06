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
		WHERE created_by = ? AND created_at > NOW() - INTERVAL 1 MINUTE`,
		auth.UserID,
	).Scan(&block)

	if err != nil {
		return true, fmt.Errorf("checkCreateSpaceThrottleBlock: %w", err)
	}

	return block, nil

}
