package spacetime

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
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

	err := db.InTransaction(conn, func(tx *sql.Tx) error {

		err := tx.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, space_type, created_at, created_by`,
			parentID, SpaceTypeSpace, time.Now(), auth.UserID,
		).Scan(&space.ID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy)

		if err != nil {
			return fmt.Errorf("insert space: %w", err)
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &space, nil

}

func CreateCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, spaceID *uint) (*Space, error) {

	// If spaceID is nil, create user checkin on parent space
	// Else, if spaceID is given, check if checkin space already exists
	// Create new checkin space if not exists

	// Ensure referenced parent space exists
	var parentExists, err = CheckSpaceExists(conn, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, fmt.Errorf("parent space does not exist: %d", parentID)
	}

	var space = Space{}

	// If spaceID is nil, create user checkin on parent space
	if spaceID == nil {

		err := db.InTransaction(conn, func(tx *sql.Tx) error {

			// Create checkin space
			err := tx.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, space_type, created_at, created_by`,
				parentID, SpaceTypeCheckin, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.SpaceType,
				&space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create associated checkin data
			err = tx.QueryRow(`INSERT INTO checkin_space (space_id)
				VALUES ($1)
				RETURNING checkin_space_id`,
				space.ID,
			).Scan(&space.CheckinSpaceID)

			if err != nil {
				return fmt.Errorf("insert checkin_space: %w", err)
			}

			// Increment target space checkin count
			err = incrementCheckinTotal(conn, parentID)
			if err != nil {
				return err
			}

			return nil

		})

		if err != nil {
			return nil, err
		}

		return &space, nil

	}

	// If spaceID is given, check if checkin space already exists

	var existingCheckinSpaceID *uint
	err = conn.QueryRow(`SELECT checkin_space.space_id
		FROM checkin_space
		INNER JOIN space ON space.id = checkin_space.space_id
		WHERE space.parent_id = $1 AND checkin_space.checkin_space_id = $2)`,
		parentID, spaceID,
	).Scan(&existingCheckinSpaceID)

	if err == sql.ErrNoRows {
		// Continue to create checkin
	} else if err != nil {
		return nil, fmt.Errorf("check checkin exists: %w", err)
	}

	if existingCheckinSpaceID != nil {

		// Check-in under existing checkin
		return CreateCheckin(conn, auth, *existingCheckinSpaceID, nil)

	}

	// Create new checkin space if not exists

	spaceExists, err := CheckSpaceExists(conn, *spaceID)
	if err != nil {
		return nil, err
	}
	if !spaceExists {
		return nil, fmt.Errorf("checkin space does not exist: %d", *spaceID)
	}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		// Create checkin space
		err := tx.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, space_type, created_at, created_by`,
			parentID, SpaceTypeCheckin, time.Now(), auth.UserID,
		).Scan(&space.ID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy)

		if err != nil {
			return fmt.Errorf("insert space: %w", err)
		}

		// Create associated checkin data
		err = tx.QueryRow(`INSERT INTO checkin_space (space_id, checkin_space_id)
			VALUES ($1, $2)
			RETURNING checkin_space_id`,
			space.ID, spaceID,
		).Scan(&space.CheckinSpaceID)

		if err != nil {
			return fmt.Errorf("insert checkin space: %w", err)
		}

		// Increment target space checkin count
		err = incrementCheckinTotal(conn, parentID)
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &space, nil

}

func CreateTitleCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, title string) (*Space, error) {

	// Load unique_text ID
	// Check for existing title space under parent
	// Create title space if not exists
	// Check-in on title space

	title = strings.TrimSpace(title)

	// Ensure referenced parent space exists
	var parentExists, err = CheckSpaceExists(conn, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, fmt.Errorf("parent space does not exist: %d", parentID)
	}

	var space = &Space{}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		var uniqueTextId *uint

		err := conn.QueryRow(`SELECT id FROM unique_text WHERE text_value = $1`,
			title,
		).Scan(&uniqueTextId)

		if err == sql.ErrNoRows {
			// Continue to create title
		} else if err != nil {
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

			// Create space
			err = tx.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, space_type, created_at, created_by`,
				parentID, SpaceTypeTitle, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.SpaceType,
				&space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create title_space
			_, err = tx.Exec(`INSERT INTO title_space (space_id, unique_text_id)
				VALUES ($1, $2)`,
				space.ID, *uniqueTextId,
			)

			if err != nil {
				return fmt.Errorf("insert title_space: %w", err)
			}

		} else {

			// Check if title_space already exists
			var existingTitleSpaceID *uint
			err = conn.QueryRow(`SELECT space.id
				FROM space
				INNER JOIN title_space ON title_space.space_id = space.id
				WHERE space.parent_id = $1 AND title_space.unique_text_id = $2`,
				parentID, *uniqueTextId,
			).Scan(&existingTitleSpaceID)

			if err == sql.ErrNoRows {
				// Continue to create title
			} else if err != nil {
				return fmt.Errorf("check title exists: %w", err)
			}

			if existingTitleSpaceID == nil {

				// Create title subspace
				err = tx.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
					VALUES ($1, $2, $3, $4)
					RETURNING id, space_type, created_at, created_by`,
					parentID, SpaceTypeTitle, time.Now(), auth.UserID,
				).Scan(&space.ID, &space.SpaceType,
					&space.CreatedAt, &space.CreatedBy)

				if err != nil {
					return fmt.Errorf("insert space: %w", err)
				}

				// Create title_space
				_, err = tx.Exec(`INSERT INTO title_space (space_id, unique_text_id)
					VALUES ($1, $2)`,
					space.ID, *uniqueTextId,
				)

				if err != nil {
					return fmt.Errorf("insert title_space: %w", err)
				}

			} else {

				// Check-in under existing title
				space, err = CreateCheckin(conn, auth, *existingTitleSpaceID, nil)

				if err != nil {
					return err
				}

			}

		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return space, nil

}

func CreateTagCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing tag space under parent
	// Create tag space if not exists
	// Check-in on tag space

	return nil, nil

}

func CreateTextCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing text space under parent
	// Create text space if not exists
	// Check-in on text space

	return nil, nil

}

func CreateNakedText(conn *sql.DB, auth ajax.Auth, parentID uint, finalText, replayData string) (*Space, error) {

	// Create naked text space with given replay data

	return nil, nil

}

func CreateStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	// Create an open stream of consciousness space
	// (will hold a series of naked texts created by author)

	return nil, nil

}

func CloseStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, id uint) error {

	// Mark stream of consciousness as "closed" by author

	return nil

}

func CreateJSONAttribute(conn *sql.DB, auth ajax.Auth, parentID uint, url, path string) (*Space, error) {

	// Check if space exists
	// Create if not exists

	return nil, nil

}

func incrementCheckinTotal(conn db.DBConn, spaceID uint) error {

	// Increment overall_checkin_total for given space

	_, err := conn.Exec(`UPDATE space
		SET overall_checkin_total = overall_checkin_total + 1
		WHERE id = $1`,
		spaceID,
	)

	if err != nil {
		return fmt.Errorf("increment checkin total: %w", err)
	}

	return nil

}
