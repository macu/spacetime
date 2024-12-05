package spacetime

import (
	"database/sql"
	"fmt"
)

func CheckSpaceExists(conn *sql.DB, spaceID uint) (bool, error) {

	var exists bool

	var err = conn.QueryRow(`SELECT EXISTS (
		SELECT 1
		FROM space
		WHERE id = $1
	)`, spaceID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check space exists: %w", err)
	}

	return exists, nil

}
