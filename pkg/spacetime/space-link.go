package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func LoadExistingSpaceLink(conn db.DBConn,
	parentID, spaceID uint,
) (*Space, error) {

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeLink,
	}

	err := conn.QueryRow(`SELECT space_id FROM link_space
		WHERE parent_space_id = $1 AND link_space_id = $2`,
		parentID, spaceID,
	).Scan(&space.ID)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select link_space: %w", err)
	}

	return space, nil

}

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

	// Get details about space to check in
	linkedSpace, err := GetSpace(conn, spaceID)
	if err != nil {
		return nil, fmt.Errorf("get space: %w", err)
	}
	if linkedSpace == nil {
		return nil, fmt.Errorf("space to check in does not exist: %d", spaceID)
	}

	// Check if space belongs directly to parent space
	if linkedSpace.ParentID != nil && *linkedSpace.ParentID == parentID {

		// Create direct checkin under existing space
		_, err := CreateCheckin(conn, auth, linkedSpace.ID)
		if err != nil {
			return nil, err
		}

		return linkedSpace, nil

	}

	// Check if this link already exists
	existingSpaceLink, err := LoadExistingSpaceLink(conn, parentID, spaceID)
	if err != nil {
		return nil, err
	}

	if existingSpaceLink == nil {

		var space = Space{
			ParentID:  &parentID,
			SpaceType: SpaceTypeLink,
		}

		err = db.InTransaction(conn, func(tx *sql.Tx) error {

			// Create space link
			err = CreateSpace(tx, auth, &space, &parentID, SpaceTypeLink)
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

	return existingSpaceLink, nil

}
