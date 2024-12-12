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

	err := conn.QueryRow(`INSERT INTO space (parent_id, space_type, created_at, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, space_type, created_at, created_by`,
		parentID, SpaceTypeSpace, time.Now(), auth.UserID,
	).Scan(&space.ID, &space.SpaceType,
		&space.CreatedAt, &space.CreatedBy)

	if err != nil {
		return nil, fmt.Errorf("insert space: %w", err)
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

		// Create checkin space
		err := conn.QueryRow(`INSERT INTO space
			(parent_id, space_type, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, space_type, created_at, created_by`,
			parentID, SpaceTypeCheckin, time.Now(), auth.UserID,
		).Scan(&space.ID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy)

		if err != nil {
			return nil, fmt.Errorf("insert space: %w", err)
		}

		var checkinSpaceID *uint = nil
		space.CheckinSpaceID = &checkinSpaceID

		var checkinSpace *Space = nil
		space.CheckinSpace = &checkinSpace

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
		return nil, fmt.Errorf("space to check in does not exist: %d", *spaceID)
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
		_, err = tx.Exec(`INSERT INTO checkin_space (space_id, checkin_space_id)
			VALUES ($1, $2)`,
			space.ID, spaceID,
		)

		if err != nil {
			return fmt.Errorf("insert checkin_space: %w", err)
		}

		space.CheckinSpaceID = &spaceID

		var checkinSpace = &Space{
			ID: *spaceID,
		}
		space.CheckinSpace = &checkinSpace

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

	var space = &Space{}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		var uniqueTextId *uint

		var runInsertTitleSpace = func() error {

			// Create space
			err = tx.QueryRow(`INSERT INTO space
				(parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, space_type, created_at, created_by`,
				parentID, SpaceTypeTitle, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.SpaceType,
				&space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create title_space
			_, err = tx.Exec(`INSERT INTO title_space
				(space_id, unique_text_id)
				VALUES ($1, $2)`,
				space.ID, *uniqueTextId,
			)

			if err != nil {
				return fmt.Errorf("insert title_space: %w", err)
			}

			return nil

		}

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
				RETURNING id, text_value`,
				title,
			).Scan(&uniqueTextId, space.Text)

			if err != nil {
				return fmt.Errorf("insert unique_text: %w", err)
			}

			if err = runInsertTitleSpace(); err != nil {
				return fmt.Errorf("insert title space: %w", err)
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

				return fmt.Errorf("check title_space exists: %w", err)

			}

			if existingTitleSpaceID == nil {

				// Create title subspace
				if err = runInsertTitleSpace(); err != nil {
					return fmt.Errorf("insert title_space: %w", err)
				}

			} else {

				// Check-in under existing title
				space, err = CreateCheckin(conn, auth, *existingTitleSpaceID, nil)

				if err != nil {
					return fmt.Errorf("create checkin: %w", err)
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

func CreateTagCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, tag string) (*Space, error) {

	// Load unique_text ID
	// Check for existing tag space under parent
	// Create tag space if not exists
	// Check-in on tag space

	tag = strings.TrimSpace(tag)

	if !ValidateTag(tag) {
		return nil, fmt.Errorf("invalid tag: %s", tag)
	}

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

		var runInsertTagSpace = func() error {

			// Create space
			err = tx.QueryRow(`INSERT INTO space
				(parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, space_type, created_at, created_by`,
				parentID, SpaceTypeTag, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.SpaceType,
				&space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create title_space
			_, err = tx.Exec(`INSERT INTO tag_space
				(space_id, unique_text_id)
				VALUES ($1, $2)`,
				space.ID, *uniqueTextId,
			)

			if err != nil {
				return fmt.Errorf("insert tag_space: %w", err)
			}

			return nil

		}

		err := conn.QueryRow(`SELECT id FROM unique_text WHERE text_value = $1`,
			tag,
		).Scan(&uniqueTextId)

		if err == sql.ErrNoRows {
			// Continue to create tag
		} else if err != nil {
			return fmt.Errorf("load unique_text ID: %w", err)
		}

		if uniqueTextId == nil {

			// Create unique_text
			err := tx.QueryRow(`INSERT INTO unique_text (text_value)
				VALUES ($1)
				RETURNING id, text_value`,
				tag,
			).Scan(&uniqueTextId, space.Text)

			if err != nil {
				return fmt.Errorf("insert unique_text: %w", err)
			}

			if err = runInsertTagSpace(); err != nil {
				return fmt.Errorf("insert tag space: %w", err)
			}

		} else {

			// Check if tag_space already exists
			var existingTagSpaceID *uint
			err = conn.QueryRow(`SELECT space.id
				FROM space
				INNER JOIN tag_space ON tag_space.space_id = space.id
				WHERE space.parent_id = $1 AND tag_space.unique_text_id = $2`,
				parentID, *uniqueTextId,
			).Scan(&existingTagSpaceID)

			if err == sql.ErrNoRows {

				// Continue to create tag

			} else if err != nil {

				return fmt.Errorf("check tag_space exists: %w", err)

			}

			if existingTagSpaceID == nil {

				// Create tag subspace
				if err = runInsertTagSpace(); err != nil {
					return fmt.Errorf("insert tag_space: %w", err)
				}

			} else {

				// Check-in under existing tag
				space, err = CreateCheckin(conn, auth, *existingTagSpaceID, nil)

				if err != nil {
					return fmt.Errorf("create checkin: %w", err)
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

func CreateTextCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing text space under parent
	// Create text space if not exists
	// Check-in on text space

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

	var space = &Space{}

	err := db.InTransaction(conn, func(tx *sql.Tx) error {

		var uniqueTextId *uint

		var runInsertTextSpace = func() error {

			// Create space
			err := tx.QueryRow(`INSERT INTO space
				(parent_id, space_type, created_at, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING id, space_type, created_at, created_by`,
				parentID, SpaceTypeText, time.Now(), auth.UserID,
			).Scan(&space.ID, &space.SpaceType,
				&space.CreatedAt, &space.CreatedBy)

			if err != nil {
				return fmt.Errorf("insert space: %w", err)
			}

			// Create title_space
			_, err = tx.Exec(`INSERT INTO text_space
				(space_id, unique_text_id)
				VALUES ($1, $2)`,
				space.ID, *uniqueTextId,
			)

			if err != nil {
				return fmt.Errorf("insert text_space: %w", err)
			}

			return nil

		}

		err := conn.QueryRow(`SELECT id FROM unique_text WHERE text_value = $1`,
			text,
		).Scan(&uniqueTextId)

		if err == sql.ErrNoRows {
			// Continue to create text
		} else if err != nil {
			return fmt.Errorf("load unique_text ID: %w", err)
		}

		if uniqueTextId == nil {

			// Create unique_text
			err := tx.QueryRow(`INSERT INTO unique_text (text_value)
				VALUES ($1)
				RETURNING id, text_value`,
				text,
			).Scan(&uniqueTextId, space.Text)

			if err != nil {
				return fmt.Errorf("insert unique_text: %w", err)
			}

			if err = runInsertTextSpace(); err != nil {
				return fmt.Errorf("insert text space: %w", err)
			}

		} else {

			// Check if text_space already exists
			var existingTextSpaceID *uint
			err = conn.QueryRow(`SELECT space.id
				FROM space
				INNER JOIN text_space ON text_space.space_id = space.id
				WHERE space.parent_id = $1 AND text_space.unique_text_id = $2`,
				parentID, *uniqueTextId,
			).Scan(&existingTextSpaceID)

			if err == sql.ErrNoRows {

				// Continue to create text

			} else if err != nil {

				return fmt.Errorf("check text_space exists: %w", err)

			}

			if existingTextSpaceID == nil {

				// Create text subspace
				if err = runInsertTextSpace(); err != nil {
					return fmt.Errorf("insert text_space: %w", err)
				}

			} else {

				// Check-in under existing text
				space, err = CreateCheckin(conn, auth, *existingTextSpaceID, nil)

				if err != nil {
					return fmt.Errorf("create checkin: %w", err)
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

func CreateNakedText(conn *sql.DB, auth ajax.Auth, parentID *uint, finalText, replayData string) (*Space, error) {

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
