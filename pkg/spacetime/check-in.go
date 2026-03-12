package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CreateCheckin(conn *sql.DB, auth ajax.Auth, parentID uint) error {

	_, err := conn.Exec(`INSERT INTO checkin (space_id, user_id, created_at)
	VALUES ($1, $2, NOW())`,
		parentID,
		auth.UserID,
	)

	if err != nil {
		return fmt.Errorf("create checkin: %w", err)
	}

	return nil

}
