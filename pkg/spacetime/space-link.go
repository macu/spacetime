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

	var spaceIDPtr = &spaceID

	var space = &Space{
		ParentID:    &parentID,
		SpaceType:   SpaceTypeLink,
		LinkSpaceID: &spaceIDPtr,
	}

	err := conn.QueryRow(`SELECT space_id FROM link_space
		WHERE parent_id = $1 AND link_space_id = $2`,
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

	// Get details about space to check in
	linkedSpace, err := GetSpace(conn, spaceID)
	if err != nil {
		return nil, fmt.Errorf("get space: %w", err)
	}
	if linkedSpace == nil {
		return nil, fmt.Errorf("space to check in does not exist: %d", spaceID)
	}

	if linkedSpace.ParentID != nil && *linkedSpace.ParentID == parentID {
		return nil, fmt.Errorf("space %d already belongs to parent space %d", spaceID, parentID)
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
				(space_id, parent_id, link_space_id)
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
