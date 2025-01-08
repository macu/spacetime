package spacetime

import (
	"database/sql"
	"fmt"
	"strings"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func LoadExistingText(conn *sql.DB,
	parentID *uint, text string,
) (*Space, error) {

	// Load text space

	var space = &Space{
		ParentID:  parentID,
		SpaceType: SpaceTypeText,
		Text:      &text,
	}

	var args = []interface{}{}

	err := conn.QueryRow(`SELECT space.id, space.created_at, space.created_by
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE `+db.Eq("space.parent_id", &args, parentID)+`
		AND space.space_type = `+db.Arg(&args, SpaceTypeText)+`
		AND unique_text.text_value = `+db.Arg(&args, text),
		args...,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select text_space: %w", err)
	}

	return space, nil

}

func CreateText(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

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

	// Check if text_space already exists
	existingText, err := LoadExistingText(conn, parentID, text)
	if err != nil {
		return nil, fmt.Errorf("check text_space exists: %w", err)
	}
	if existingText != nil {
		err = LoadSubspaceCount(conn, []*Space{existingText})
		if err != nil {
			return nil, fmt.Errorf("load subspace count: %w", err)
		}
		return existingText, nil
	}

	var space = &Space{
		ParentID:  parentID,
		SpaceType: SpaceTypeText,
		Text:      &text,
	}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		uniqueTextId, err := GetUniqueTextId(conn, text)
		if err != nil {
			return err
		}

		if uniqueTextId == nil {
			uniqueTextId, err = CreateUniqueText(conn, text)
			if err != nil {
				return err
			}
		}

		// Create space
		err = CreateSpace(tx, auth, space, parentID, SpaceTypeText)
		if err != nil {
			return err
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

	})

	if err != nil {
		return nil, fmt.Errorf("create text: %w", err)
	}

	return space, nil

}
