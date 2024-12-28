package spacetime

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateTitleCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, title string) (*Space, error) {

	// Load unique_text ID
	// Check for existing title space under parent
	// Create title space if not exists
	// Check-in on title space

	title = strings.TrimSpace(title)

	if !ValidateTitle(title) {
		return nil, fmt.Errorf("invalid title: %s", title)
	}

	// Ensure referenced parent space exists
	var parentExists, err = CheckSpaceExists(conn, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, fmt.Errorf("parent space does not exist: %d", parentID)
	}

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeTitle,
		Text:      &title,
	}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		var uniqueTextId *uint

		// Create function to insert title space
		var runInsertTitleSpace = func() error {

			// Create space
			err = tx.QueryRow(`INSERT INTO space
				(parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at, created_by`,
				parentID, SpaceTypeTitle, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create title_space
			_, err = tx.Exec(`INSERT INTO title_space
				(space_id, parent_space_id, unique_text_id)
				VALUES ($1, $2, $3)`,
				space.ID, parentID, *uniqueTextId,
			)

			if err != nil {
				return fmt.Errorf("insert title_space: %w", err)
			}

			return nil

		}

		// Check for existing unique_text
		err := tx.QueryRow(`SELECT id FROM unique_text WHERE text_value = $1`,
			title,
		).Scan(&uniqueTextId)

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("load unique_text ID: %w", err)
		}

		if uniqueTextId == nil {

			// Create unique_text
			err := tx.QueryRow(`INSERT INTO unique_text (text_value)
				VALUES ($1)
				RETURNING id`,
				title,
			).Scan(&uniqueTextId)

			if err != nil {
				return fmt.Errorf("insert unique_text: %w", err)
			}

			// Create title space now that uniqueTextId is available
			if err = runInsertTitleSpace(); err != nil {
				return fmt.Errorf("insert title space: %w", err)
			}

		} else {

			// Check if title_space already exists
			var existingID *uint
			var existingCreatedAt time.Time
			var existingCreatedBy uint
			err = conn.QueryRow(`SELECT space.id, space.parent_id,
				space.created_at, space.created_by
				FROM space
				INNER JOIN title_space ON title_space.space_id = space.id
				WHERE space.parent_id = $1
				AND space.space_type = $2
				AND title_space.unique_text_id = $3`,
				parentID, SpaceTypeTitle, *uniqueTextId,
			).Scan(&existingID, &existingCreatedAt, &existingCreatedBy)

			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("check title_space exists: %w", err)
			}

			if existingID == nil {

				// Create title subspace
				if err = runInsertTitleSpace(); err != nil {
					return fmt.Errorf("insert title_space: %w", err)
				}

			} else {

				// Return existing space
				space.ID = *existingID
				space.CreatedAt = existingCreatedAt
				space.CreatedBy = existingCreatedBy

				// Check-in under existing title
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
		return nil, fmt.Errorf("create title checkin: %w", err)
	}

	return space, nil

}
