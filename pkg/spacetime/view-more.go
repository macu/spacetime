package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
)

// call to load user's bookmarked titles ("view all")
func loadLastBookmarkedTitlesByDate(conn *sql.DB, auth ajax.Auth,
	spaceID uint,
	date *time.Time, // optional
	limit uint,
) ([]*Space, error) {

	if limit > MaxSubspacesPageLimit {
		limit = MaxSubspacesPageLimit
	}

	var args = []interface{}{auth.UserID, spaceID, SpaceTypeTitle, limit}

	var dateClause string
	if date != nil {
		dateClause = `AND user_space_bookmark.created_at < $5`
		args = append(args, date)
	}

	rows, err := conn.Query(`SELECT
		space.id, space.space_type, space.created_at, space.created_by, unique_text.text
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		INNER JOIN user_space_bookmark ON user_space_bookmark.space_id = space.id
		WHERE user_space_bookmark.user_id = $1
		AND user_space_bookmark.space_id = $2
		AND space.space_type = $3
		`+dateClause+`
		ORDER BY user_space_bookmark.created_at DESC
		LIMIT $4`,
		args...,
	)

	var spaces = []*Space{}

	if err == sql.ErrNoRows {
		return spaces, nil
	} else if err != nil {
		return nil, fmt.Errorf("loading bookmarked titles: %w", err)
	}

	for rows.Next() {
		var space = Space{}
		err = rows.Scan(&space.ID, &space.SpaceType, &space.CreatedAt, &space.CreatedBy,
			&space.Text)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
		spaces = append(spaces, &space)
	}

	return spaces, nil

}
