package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateText(conn *sql.DB, auth ajax.Auth, parentID uint, text string) (*Space, error) {

	if !ValidateText(text) {
		return nil, fmt.Errorf("invalid text")
	}

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeText,
		Text:      &text,
	}

	err := db.InTransaction(conn, func(tx *sql.Tx) error {

		uniqueTextId, err := GetOrCreateUniqueTextId(tx, text)
		if err != nil {
			return err
		} else if uniqueTextId == nil {
			return fmt.Errorf("unique text id is nil")
		}

		// Create space
		if err = CreateSpace(tx, auth, space, &parentID, SpaceTypeText); err != nil {
			return err
		}

		// Create text_space
		if _, err = tx.Exec(`INSERT INTO text_space
			(space_id, parent_id, text_id)
			VALUES ($1, $2, $3)`,
			space.ID, parentID, *uniqueTextId,
		); err != nil {
			return fmt.Errorf("insert text_space: %w", err)
		}

		return nil

	})

	if err != nil {
		return nil, fmt.Errorf("create text: %w", err)
	}

	return space, nil

}
