package spacetime

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

const MaxSubspacesPageLimit = 20

func LoadTopSpaces(conn *sql.DB, auth *ajax.Auth,
	parentID *uint, // optional
	afterSpaceID *uint, // pagination
	limit uint,
) ([]*Space, error) {

	var spaces = []*Space{}

	var args = []interface{}{}

	var bookmarkFieldSql string
	if auth != nil {
		bookmarkFieldSql = `EXISTS(SELECT 1 FROM user_space_bookmark
			WHERE user_space_bookmark.user_id = ` + db.Arg(&args, auth.UserID) + `
			AND user_space_bookmark.space_id = space.id) AS user_bookmark`
	} else {
		bookmarkFieldSql = `FALSE AS user_bookmark`
	}

	var parentClauseSql string
	if parentID != nil {
		parentClauseSql = `AND space.parent_id = ` + db.Arg(&args, *parentID)
	} else {
		parentClauseSql = `AND space.parent_id IS NULL`
	}

	rows, err := conn.Query(`SELECT space.id, space.space_type, space.created_at, space.created_by,
		`+bookmarkFieldSql+`
		FROM space
		WHERE space.space_type NOT IN (
			'`+db.Arg(&args, SpaceTypeTitle)+`',
			'`+db.Arg(&args, SpaceTypeTag)+`'
		)
		`+parentClauseSql+`
		ORDER BY space.created_at DESC
		LIMIT `+db.Arg(&args, limit),
		args...,
	)

	defer rows.Close()

	if err == sql.ErrNoRows {
		return spaces, nil
	} else if err != nil {
		return nil, fmt.Errorf("loading top spaces: %w", err)
	}

	for rows.Next() {
		var space = Space{}
		err = rows.Scan(&space.ID, &space.SpaceType, &space.CreatedAt, &space.CreatedBy,
			&space.UserBookmark)
		if err != nil {
			return nil, fmt.Errorf("loading top spaces: %w", err)
		}
		spaces = append(spaces, &space)
	}

	err = loadSpaceDetails(conn, auth, spaces, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	if auth != nil {
		err = loadLastBookmarkedTitleForSpaces(conn, *auth, spaces)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
	}

	err = loadTopTitleForSpaces(conn, spaces)
	if err != nil {
		return nil, fmt.Errorf("loading top titles: %w", err)
	}

	return spaces, nil

}

// --------------------------------------------------
// batch load functions

func loadSpaceDetails(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration,
) error {
	// Load content for multiple spaces

	if hasSpacesOfType(spaces, SpaceTypeTitle) {
		var titleSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeTitle {
				titleSpaces = append(titleSpaces, space)
			}
		}
		loadTitleSpacesContent(conn, auth, titleSpaces)
	}

	if hasSpacesOfType(spaces, SpaceTypeTag) {
		var tagSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeTag {
				tagSpaces = append(tagSpaces, space)
			}
		}
		loadTagSpacesContent(conn, auth, tagSpaces)
	}

	if hasSpacesOfType(spaces, SpaceTypeText) {
		var textSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeText {
				textSpaces = append(textSpaces, space)
			}
		}
		loadTextSpacesContent(conn, auth, textSpaces)
	}

	if hasSpacesOfType(spaces, SpaceTypeNaked) {
		var nakedTextSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeNaked {
				nakedTextSpaces = append(nakedTextSpaces, space)
			}
		}
		loadNakedTextSpacesContent(conn, auth, nakedTextSpaces)
	}

	if hasSpacesOfType(spaces, SpaceTypeCheckin) {
		var checkinSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeCheckin {
				checkinSpaces = append(checkinSpaces, space)
			}
		}
		loadCheckinSpaceDetails(conn, auth, checkinSpaces)
	}

	return nil

}

func loadLastBookmarkedTitleForSpaces(conn *sql.DB, auth ajax.Auth, spaces []*Space) error {
	// Load the last bookmarked title for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	var inClauseSql string

	var args = []interface{}{}

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
	}

	rows, err := conn.Query(``) // TODO

	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading last bookmarked titles: %w", err)
	}

	for rows.Next() {
		var spaceID uint
		var parentID uint
		var text string
		err = rows.Scan(&spaceID, &text)
		if err != nil {
			return fmt.Errorf("loading last bookmarked titles: %w", err)
		}
		for _, space := range spaces {
			if space.ID == parentID {
				var title = &Space{
					ID:        spaceID,
					SpaceType: SpaceTypeTitle,
					Text:      &text,
				}
				var titles = []*Space{title}
				space.BookmarkedTitles = &titles
			}
		}
	}

	return nil

}

func loadTopTitleForSpaces(conn *sql.DB, spaces []*Space) error {
	// Load the top title for multiple spaces

	return nil

}

func hasSpacesOfType(spaces []*Space, spaceType string) bool {
	// Check if a space of a certain type exists in a list of spaces

	for _, space := range spaces {
		if space.SpaceType == spaceType {
			return true
		}
	}

	return false

}

func loadTitleSpacesContent(conn *sql.DB, auth *ajax.Auth, spaces []*Space) error {
	// Load title content for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	var args = []interface{}{}

	var inClauseSql string

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
		args = append(args, space.ID)
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading title spaces content: %w", err)
	}

	for rows.Next() {
		var spaceID uint
		var text string
		err = rows.Scan(&spaceID, &text)
		if err != nil {
			return fmt.Errorf("loading title spaces content: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.Text = &text
			}
		}
	}

	return nil

}

func loadTagSpacesContent(conn *sql.DB, auth *ajax.Auth, spaces []*Space) error {
	// Load tag content for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	var args = []interface{}{}

	var inClauseSql string

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
		args = append(args, space.ID)
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text
		FROM space
		INNER JOIN tag_space ON tag_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = tag_space.unique_text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading tag spaces content: %w", err)
	}

	for rows.Next() {
		var spaceID uint
		var text string
		err = rows.Scan(&spaceID, &text)
		if err != nil {
			return fmt.Errorf("loading tag spaces content: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.Text = &text
			}
		}
	}

	return nil

}

func loadTextSpacesContent(conn *sql.DB, auth *ajax.Auth, spaces []*Space) error {
	// Load text content for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	var args = []interface{}{}

	var inClauseSql string

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
		args = append(args, space.ID)
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text
		FROM space
		INNER JOIN text_space ON text_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = text_space.unique_text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading text spaces content: %w", err)
	}

	for rows.Next() {
		var spaceID uint
		var text string
		err = rows.Scan(&spaceID, &text)
		if err != nil {
			return fmt.Errorf("loading text spaces content: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.Text = &text
			}
		}
	}

	return nil

}

func loadNakedTextSpacesContent(conn *sql.DB, auth *ajax.Auth, spaces []*Space) error {
	// Load naked text

	if len(spaces) == 0 {
		return nil
	}

	var args = []interface{}{}

	var inClauseSql string

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
		args = append(args, space.ID)
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text, naked_text_space.replay_data
		FROM space
		INNER JOIN naked_text_space ON naked_text_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = naked_text_space.final_unique_text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading title spaces content: %w", err)
	}

	for rows.Next() {
		var spaceID uint
		var text string
		var replayJSON string
		var replayData map[string]interface{}
		err = rows.Scan(&spaceID, &text, &replayJSON)
		if err != nil {
			return fmt.Errorf("loading title spaces content: %w", err)
		}
		err = json.Unmarshal([]byte(replayJSON), &replayData)
		if err != nil {
			return fmt.Errorf("unmarshalling naked text replay: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.Text = &text
				space.ReplayData = &replayData
			}
		}
	}

	return nil

}

func loadCheckinSpaceDetails(conn *sql.DB, auth *ajax.Auth, spaces []*Space) error {
	// Load checkin content for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	var args = []interface{}{}

	var inClauseSql string

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
		args = append(args, space.ID)
	}

	rows, err := conn.Query(`SELECT
		space.id, check_space.checkin_space_id
		FROM space
		INNER JOIN checkin_space check_space ON check_space.space_id = space.id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading checkin space details: %w", err)
	}

	for rows.Next() {
		var spaceID uint
		var checkinSpaceID *uint
		err = rows.Scan(&spaceID, &checkinSpaceID)
		if err != nil {
			return fmt.Errorf("loading checkin space details: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.CheckinSpaceID = checkinSpaceID
			}
		}
	}

	var checkinSpaces []*Space
	for _, space := range spaces {
		if space.CheckinSpaceID != nil {
			var checkinSpace = Space{
				ID: *space.CheckinSpaceID,
			}
			checkinSpaces = append(checkinSpaces, &checkinSpace)
		}
	}

	err = loadSpaceDetails(conn, auth, checkinSpaces, nil, 0)
	if err != nil {
		return fmt.Errorf("loading checkin space details: %w", err)
	}

	for _, space := range spaces {
		if space.CheckinSpaceID != nil {
			for _, checkinSpace := range checkinSpaces {
				if checkinSpace.ID == *space.CheckinSpaceID {
					space.CheckinSpace = &checkinSpace
				}
			}
		}
	}

	return nil

}