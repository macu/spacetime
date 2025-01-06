package user

import (
	"fmt"
	"spacetime/pkg/utils/db"
	"time"
)

func CheckAdmin(db db.DBConn, userID uint) bool {
	var userRole string
	err := db.QueryRow(`SELECT role FROM user_account WHERE id = $1`, userID).Scan(&userRole)
	if err != nil {
		return false
	}
	return userRole == string(RoleAdmin)
}

func BookmarkSpace(db db.DBConn, userID uint, spaceID uint, bookmark bool) error {
	if bookmark {
		_, err := db.Exec(`INSERT INTO user_space_bookmark
			(user_id, space_id, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, space_id) DO UPDATE SET created_at = $3`,
			userID, spaceID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to bookmark space: %w", err)
		}
	} else {
		_, err := db.Exec(`DELETE FROM user_space_bookmark
			WHERE user_id = $1 AND space_id = $2`,
			userID, spaceID)
		if err != nil {
			return fmt.Errorf("failed to unbookmark space: %w", err)
		}
	}

	return nil
}
