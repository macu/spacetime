package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateSubspace(conn *sql.DB, auth ajax.Auth,
	parentID *uint, label string,
) (*Space, error) {

	var space = &Space{
		ParentID:  parentID,
		SpaceType: SpaceTypeSpace,
	}

	err := db.InTransaction(conn, func(tx *sql.Tx) error {

		uniqueTextID, err := GetOrCreateUniqueTextId(tx, label)
		if err != nil {
			return err
		}

		// Create space
		if err = CreateSpace(tx, auth, space, parentID, SpaceTypeSpace); err != nil {
			return err
		}

		if parentID == nil {
			// 0 for subspace table
			parentID = new(uint)
		}

		if _, err = tx.Exec(`INSERT INTO subspace
			(parent_id, space_id, label_text_id)
			VALUES ($1, $2, $3)`,
			parentID, space.ID, uniqueTextID,
		); err != nil {
			return fmt.Errorf("insert subspace: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("create subspace: %w", err)
	}

	return space, nil

}
