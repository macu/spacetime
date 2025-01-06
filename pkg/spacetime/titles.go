package spacetime

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func LoadExistingTitle(conn *sql.DB,
	parentID uint, title string,
) (*Space, error) {

	// Load title space

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeTitle,
		Text:      &title,
	}

	var args = []interface{}{}

	err := conn.QueryRow(`SELECT space.id, space.created_at, space.created_by
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.parent_id = `+db.Arg(&args, parentID)+`
		AND space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
		AND unique_text.text_value = `+db.Arg(&args, title),
		args...,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select title_space: %w", err)
	}

	return space, nil

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
			err = CreateSpace(tx, auth, space, &parentID, SpaceTypeTitle)
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
			existingTitle, err := LoadExistingTitle(conn, parentID, title)
			if err != nil {
				return fmt.Errorf("check title_space exists: %w", err)
			}

			if existingTitle == nil {

				// Create title subspace
				if err = runInsertTitleSpace(); err != nil {
					return fmt.Errorf("insert title_space: %w", err)
				}

			} else {

				space = existingTitle

				// Check-in under existing title
				_, err = CreateCheckin(conn, auth, space.ID)

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

func LoadOriginalTitles(conn *sql.DB, spaces []*Space) error {
	// Load earliest associated titles by same creator

	if len(spaces) == 0 {
		return nil
	}

	var allTitles = []*Space{}

	for _, space := range spaces {

		// Load one by one

		var args = []interface{}{}
		var spaceID uint
		var createdAt time.Time
		var createdBy uint
		var text string

		err := conn.QueryRow(`SELECT space.id,
		space.created_at, space.created_by,
		unique_text.text_value
		FROM space AS space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
		AND space.parent_id = `+db.Arg(&args, space.ID)+`
		AND space.created_by = `+db.Arg(&args, space.CreatedBy)+`
		ORDER BY space.created_at ASC
		LIMIT 1`,
			args...,
		).Scan(&spaceID, &createdAt, &createdBy, &text)

		if err == sql.ErrNoRows {

			var nullTitle *Space = nil
			space.OriginalTitle = &nullTitle

			continue

		} else if err != nil {
			return fmt.Errorf("loading original titles: %w", err)
		}

		var titleSpace = &Space{
			ID:        spaceID,
			ParentID:  &space.ID,
			SpaceType: SpaceTypeTitle,
			Text:      &text,
			CreatedAt: createdAt,
			CreatedBy: createdBy,
		}

		space.OriginalTitle = &titleSpace

		allTitles = append(allTitles, titleSpace)

	}

	err := LoadSubspaceCount(conn, allTitles)
	if err != nil {
		return err
	}

	return nil

}

func LoadLastUserTitles(conn *sql.DB, auth ajax.Auth,
	spaces []*Space,
) error {
	// Load user titles by last checkin

	if len(spaces) == 0 {
		return nil
	}

	var allTitles = []*Space{}

	for _, space := range spaces {

		// Load one by one

		var args = []interface{}{}
		var spaceID uint
		var createdAt time.Time
		var createdBy uint
		var text string
		var lastCheckin time.Time

		err := conn.QueryRow(`SELECT space.id,
			space.created_at, space.created_by,
			unique_text.text_value,
			MAX(checkin_space.created_at) AS last_checkin
			FROM space AS space
			INNER JOIN title_space ON title_space.space_id = space.id
			INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
			INNER JOIN space AS checkin_space ON checkin_space.parent_id = space.id
			WHERE space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
			AND space.parent_id = `+db.Arg(&args, space.ID)+`
			AND checkin_space.space_type = `+db.Arg(&args, SpaceTypeCheckin)+`
			AND checkin_space.created_by = `+db.Arg(&args, auth.UserID)+`
			GROUP BY space.id, space.created_at, space.created_by,
				unique_text.text_value
			ORDER BY last_checkin DESC
			LIMIT 1`,
			args...,
		).Scan(&spaceID, &createdAt, &createdBy, &text, &lastCheckin)

		if err == sql.ErrNoRows {

			var nullTitle *Space = nil
			space.UserTitle = &nullTitle

			continue

		} else if err != nil {
			return fmt.Errorf("loading user titles by last checkin: %w", err)
		}

		var titleSpace = &Space{
			ID:        spaceID,
			ParentID:  &space.ID,
			SpaceType: SpaceTypeTitle,
			Text:      &text,
			CreatedAt: createdAt,
			CreatedBy: createdBy,
		}

		space.UserTitle = &titleSpace

		allTitles = append(allTitles, titleSpace)

	}

	err := LoadSubspaceCount(conn, allTitles)
	if err != nil {
		return err
	}

	return nil

}

func LoadTopTitles(conn *sql.DB, spaces []*Space) error {

	if len(spaces) == 0 {
		return nil
	}

	for _, space := range spaces {

		var args = []interface{}{}
		var spaceID uint
		var createdAt time.Time
		var createdBy uint
		var text string
		var totalSubspaces uint

		err := conn.QueryRow(`SELECT space.id,
			space.created_at, space.created_by, unique_text.text_value,
			COUNT(subspace.id) AS subspaces_total
			FROM space
			INNER JOIN title_space ON title_space.space_id = space.id
			INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
			LEFT JOIN space AS subspace ON subspace.parent_id = space.id
			WHERE space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
				AND space.parent_id = `+db.Arg(&args, space.ID)+`
			GROUP BY space.id, space.created_at, space.created_by,
				unique_text.text_value
			ORDER BY subspaces_total DESC
			LIMIT 1`,
			args...,
		).Scan(&spaceID, &createdAt, &createdBy, &text, &totalSubspaces)

		if err == sql.ErrNoRows {

			var nullTitle *Space = nil
			space.TopTitle = &nullTitle

			continue

		} else if err != nil {
			return fmt.Errorf("loading top title: %w", err)
		}

		var title = &Space{
			ID:             spaceID,
			ParentID:       &space.ID,
			SpaceType:      SpaceTypeTitle,
			Text:           &text,
			CreatedAt:      createdAt,
			CreatedBy:      createdBy,
			TotalSubspaces: totalSubspaces,
		}

		space.TopTitle = &title

	}

	return nil

}

func LoadMoreTitles(conn *sql.DB,
	parentId uint, offset uint, limit uint,
) (*[]Space, error) {

	return nil, nil

}
