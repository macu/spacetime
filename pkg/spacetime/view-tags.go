package spacetime

import (
	"database/sql"
	"fmt"
)

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

func LoadMoreTags(conn *sql.DB,
	parentId uint, offset uint, limit uint,
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
