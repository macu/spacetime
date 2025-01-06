package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateSpace(conn db.DBConn, auth ajax.Auth,
	space *Space, parentID *uint, spaceType string,
) error {

	err := conn.QueryRow(`INSERT INTO space
		(parent_id, space_type, created_at, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, created_by`,
		parentID, spaceType, time.Now(), auth.UserID,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

	if err != nil {
		return fmt.Errorf("insert space: %w", err)
	}

	return nil

}

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

	err := CreateSpace(conn, auth, &space, parentID, SpaceTypeSpace)
	if err != nil {
		return nil, fmt.Errorf("create empty space: %w", err)
	}

	return &space, nil

}

func GetSpace(conn db.DBConn, spaceID uint) (*Space, error) {

	var space = &Space{
		ID: spaceID,
	}

	err := conn.QueryRow(`SELECT parent_id, space_type, created_at, created_by
		FROM space WHERE id = $1`,
		spaceID,
	).Scan(&space.ParentID, &space.SpaceType, &space.CreatedAt, &space.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select space: %w", err)
	}

	return space, nil

}
