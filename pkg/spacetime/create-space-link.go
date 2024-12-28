package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateSpaceLink(conn *sql.DB, auth ajax.Auth, parentID, spaceID uint) (*Space, error) {

	// Create new space link
	// If space itself belongs to parent space, create checkin under the space

	// Ensure referenced parent space exists
	var parentExists, err = CheckSpaceExists(conn, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, fmt.Errorf("parent space does not exist: %d", parentID)
	}

	spaceExists, err := CheckSpaceExists(conn, spaceID)
	if err != nil {
		return nil, err
	}
	if !spaceExists {
		return nil, fmt.Errorf("space to check in does not exist: %d", spaceID)
	}

	// Get details about space to check in
	existingSpaceParentID, _, err := GetSpaceMeta(conn, spaceID)
	if err != nil {
		return nil, err
	}

	if existingSpaceParentID != nil && *existingSpaceParentID == parentID {

		// Create direct checkin under existing space
		return CreateCheckin(conn, auth, spaceID)

	}

	var space = Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeLink,
	}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		// Create space link
		err := tx.QueryRow(`INSERT INTO space
			(parent_id, space_type, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, created_by`,
			parentID, SpaceTypeLink, time.Now(), auth.UserID,
		).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

		if err != nil {
			return fmt.Errorf("insert space: %w", err)
		}

		// Create associated data
		_, err = tx.Exec(`INSERT INTO link_space
			(space_id, parent_space_id, link_space_id)
			VALUES ($1, $2, $3)`,
			space.ID, parentID, spaceID,
		)

		if err != nil {
			return fmt.Errorf("insert space_link_space: %w", err)
		}

		var linkSpaceID = &spaceID
		space.LinkSpaceID = &linkSpaceID

		var linkSpace *Space = nil
		space.LinkSpace = &linkSpace // not loaded

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &space, nil

}
