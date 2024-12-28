package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
)

func CreateEmptySpace(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	// Create new space (nameless)

	var space = Space{}

	// If given, check if parent space exists
	if parentID != nil {
		var exists, err = CheckSpaceExists(conn, *parentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("parent space does not exist: %d", *parentID)
		}
	}

	err := conn.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, space_type, created_at, created_by`,
		parentID, SpaceTypeSpace, time.Now(), auth.UserID,
	).Scan(&space.ID, &space.SpaceType,
		&space.CreatedAt, &space.CreatedBy)

	if err != nil {
		return nil, fmt.Errorf("insert space: %w", err)
	}

	return &space, nil

}
