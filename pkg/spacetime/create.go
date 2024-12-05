package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateEmptySpace(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	var space = Space{}

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

	// Check if checkin already exists

	if spaceID != nil {

		var existingCheckinID *uint
		err := conn.QueryRow(`SELECT checkin_space.space_id
			FROM checkin_space
			INNER JOIN space ON space.id = checkin_space.space_id
			WHERE space.parent_id = $1 AND checkin_space.checkin_space_id = $2)`,
			parentID, spaceID,
		).Scan(&existingCheckinID)

		if err == sql.ErrNoRows {
			// Continue to create checkin
		} else if err != nil {
			return nil, fmt.Errorf("check checkin exists: %w", err)
		}

		if existingCheckinID != nil {

			// Check-in under existing checkin
			return CreateCheckin(conn, auth, *existingCheckinID, nil)

		}

	}

	// Create checkin space if not exists

	var parentExists, err = CheckSpaceExists(conn, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, fmt.Errorf("parent space does not exist: %d", parentID)
	}

	if spaceID != nil {
		var spaceExists, err = CheckSpaceExists(conn, *spaceID)
		if err != nil {
			return nil, err
		}
		if !spaceExists {
			return nil, fmt.Errorf("checkin space does not exist: %d", *spaceID)
		}
	}

	var space = Space{}

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
		_, err = tx.Exec(`UPDATE space
			SET overall_checkin_total = overall_checkin_total + 1
			WHERE id = $1`,
			parentID,
		)

		if err != nil {
			return fmt.Errorf("increment checkin count: %w", err)
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &space, nil

}

func CreateTitleCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, title string) (*Space, error) {

	// Load unique_text ID
	// Check for existing title space under parent
	// Create title space if not exists
	// Check-in on title space

	return nil, nil

}

func CreateTagCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing tag space under parent
	// Create tag space if not exists
	// Check-in on tag space

	return nil, nil

}

func CreateTextCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing text space under parent
	// Create text space if not exists
	// Check-in on text space

	return nil, nil

}

func CreateNakedText(conn *sql.DB, auth ajax.Auth, parentID *uint, finalText, replayData string) (*Space, error) {

	// Create naked text space with given replay data

	return nil, nil

}

func CreateStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	// Create an open stream of consciousness space

	return nil, nil

}

func CloseStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, id *uint) error {

	// Mark stream of consciousness as "closed" by user who created it

	return nil

}

func CreateJSONAttribute(conn *sql.DB, auth ajax.Auth, parentID *uint, url, path string) (*Space, error) {

	// Check if space exists
	// Create if not exists

	return nil, nil

}
