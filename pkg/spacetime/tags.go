package spacetime

import (
	"database/sql"
	"fmt"
	"strings"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func LoadExistingTag(conn *sql.DB,
	parentID uint, tag string,
) (*Space, error) {

	// Load tag space

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeTag,
		Text:      &tag,
	}

	var args = []interface{}{}

	err := conn.QueryRow(`SELECT space.id, space.created_at, space.created_by
		FROM space
		INNER JOIN tag_space ON tag_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = tag_space.unique_text_id
		WHERE space.parent_id = `+db.Arg(&args, parentID)+`
		AND space.space_type = `+db.Arg(&args, SpaceTypeTag)+`
		AND unique_text.text_value = `+db.Arg(&args, tag),
		args...,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select tag_space: %w", err)
	}

	return space, nil

}

func CreateTag(conn *sql.DB, auth ajax.Auth, parentID uint, tag string) (*Space, error) {

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

	// Check if tag already exists
	existingTag, err := LoadExistingTag(conn, parentID, tag)
	if err != nil {
		return nil, fmt.Errorf("check tag exists: %w", err)
	}
	if existingTag != nil {
		err = LoadSubspaceCount(conn, []*Space{existingTag})
		if err != nil {
			return nil, err
		}
		return existingTag, nil
	}

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeTag,
		Text:      &tag,
	}

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		uniqueTextId, err := GetUniqueTextId(tx, tag)
		if err != nil {
			return err
		}

		if uniqueTextId == nil {
			uniqueTextId, err = CreateUniqueText(tx, tag)
			if err != nil {
				return err
			}
		}

		// Create space
		err = CreateSpace(tx, auth, space, &parentID, SpaceTypeTag)
		if err != nil {
			return err
		}

		// Create tag_space
		_, err = tx.Exec(`INSERT INTO tag_space
			(space_id, parent_space_id, unique_text_id)
			VALUES ($1, $2, $3)`,
			space.ID, parentID, *uniqueTextId,
		)
		if err != nil {
			return fmt.Errorf("insert tag_space: %w", err)
		}

		return nil

	})

	if err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}

	return space, nil

}

func LoadTopTags(conn *sql.DB, spaces []*Space,
	offset uint, limit uint,
) error {
	// Load top tags for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	if limit > MaxSubspacesPageLimit {
		limit = MaxSubspacesPageLimit
	}

	for _, space := range spaces {

		tags, err := LoadMoreTags(conn, space.ID, offset, limit)
		if err != nil {
			return err
		}

		space.TopTags = tags

	}

	return nil

}

func LoadMoreTags(conn *sql.DB, parentId uint,
	offset uint, limit uint,
) (*[]*Space, error) {

	rows, err := conn.Query(`SELECT space.id, unique_text.text_value,
		COUNT(subspace.id) AS subspace_total
		FROM space
		INNER JOIN tag_space ON tag_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = tag_space.unique_text_id
		LEFT JOIN space AS subspace ON subspace.parent_id = space.id
		WHERE space.space_type = $1
		AND space.parent_id = $2
		GROUP BY space.id, unique_text.text_value
		ORDER BY subspace_total DESC
		OFFSET $3
		LIMIT $4`,
		SpaceTypeTag, parentId, offset, limit,
	)

	if err != nil {
		return nil, fmt.Errorf("loading top tags: %w", err)
	}

	defer rows.Close()

	var tags = []*Space{}

	for rows.Next() {
		var spaceID uint
		var text string
		var subspacesTotal uint
		err = rows.Scan(&spaceID, &text, &subspacesTotal)
		if err != nil {
			return nil, fmt.Errorf("scanning top tags: %w", err)
		}
		var tag = &Space{
			ID:             spaceID,
			SpaceType:      SpaceTypeTag,
			Text:           &text,
			TotalSubspaces: subspacesTotal,
		}
		tags = append(tags, tag)
	}

	return &tags, nil

}
