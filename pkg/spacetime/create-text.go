package spacetime

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateTextCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing tag space under parent
	// Create tag space if not exists
	// Check-in on tag space

	text = strings.TrimSpace(text)

	if !ValidateText(text) {
		return nil, fmt.Errorf("invalid text: %s", text)
	}

	if parentID != nil {
		// Ensure referenced parent space exists
		var parentExists, err = CheckSpaceExists(conn, *parentID)
		if err != nil {
			return nil, err
		}
		if !parentExists {
			return nil, fmt.Errorf("parent space does not exist: %d", parentID)
		}
	}

	var space = &Space{
		ParentID:  parentID,
		SpaceType: SpaceTypeText,
		Text:      &text,
	}

	err := db.InTransaction(conn, func(tx *sql.Tx) error {

		uniqueTextId, err := GetUniqueTextId(conn, text)
		if err != nil {
			return err
		}

		// Create function to insert text space
		var runInsertTextSpace = func() error {

			// Create space
			err := tx.QueryRow(`INSERT INTO space
				(parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at, created_by`,
				parentID, SpaceTypeText, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create text_space
			_, err = tx.Exec(`INSERT INTO text_space
				(space_id, parent_space_id, unique_text_id)
				VALUES ($1, $2, $3)`,
				space.ID, parentID, *uniqueTextId,
			)

			if err != nil {
				return fmt.Errorf("insert text_space: %w", err)
			}

			return nil

		}

		if uniqueTextId == nil {

			uniqueTextId, err = CreateUniqueText(conn, text)
			if err != nil {
				return err
			}

			// Create text space now that uniqueTextId is available
			if err = runInsertTextSpace(); err != nil {
				return fmt.Errorf("insert text space: %w", err)
			}

		} else {

			// Check if text_space already exists
			var existingID *uint
			var existingCreatedAt time.Time
			var existingCreatedBy uint
			err = conn.QueryRow(`SELECT space.id,
				space.created_at, space.created_by
				FROM space
				INNER JOIN text_space ON text_space.space_id = space.id
				WHERE space.parent_id = $1
				AND space.space_type = $2
				AND text_space.unique_text_id = $3`,
				parentID, SpaceTypeText, *uniqueTextId,
			).Scan(&existingID, &existingCreatedAt, &existingCreatedBy)

			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("check text_space exists: %w", err)
			}

			if existingID == nil {

				// Create text subspace
				if err = runInsertTextSpace(); err != nil {
					return fmt.Errorf("insert text_space: %w", err)
				}

			} else {

				// Return existing space details
				space.ID = *existingID
				space.CreatedAt = existingCreatedAt
				space.CreatedBy = existingCreatedBy

				// Check-in under existing text
				_, err = CreateCheckin(conn, auth, *existingID)
				if err != nil {
					return fmt.Errorf("create checkin: %w", err)
				}

				err = LoadSubspaceCount(conn, []*Space{space})
				if err != nil {
					return fmt.Errorf("load subspace count: %w", err)
				}

			}

		}

		return nil

	})

	if err != nil {
		return nil, fmt.Errorf("create text checkin: %w", err)
	}

	return space, nil

}
