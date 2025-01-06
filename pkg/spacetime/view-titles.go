package spacetime

import (
	"database/sql"
	"fmt"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
	"time"
)

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
