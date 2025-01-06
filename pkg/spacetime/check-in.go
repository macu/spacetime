package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CreateCheckin(conn *sql.DB, auth ajax.Auth, parentID uint) (*Space, error) {

	// Create new user checkin

	// Ensure referenced parent space exists
	var parentExists, err = CheckSpaceExists(conn, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, fmt.Errorf("parent space does not exist: %d", parentID)
	}

	var space = Space{}

	err = CreateSpace(conn, auth, &space, &parentID, SpaceTypeCheckin)
	if err != nil {
		return nil, fmt.Errorf("create checkin: %w", err)
	}

	return &space, nil

}
