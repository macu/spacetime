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

func LoadSpace(conn *sql.DB, auth *ajax.Auth, id uint) (*Space, error) {
	// Load a single space

	var space = Space{
		ID: id,
	}

	err := conn.QueryRow(`SELECT
		space.space_type, space.created_at, space.created_by
		FROM space
		WHERE space.id = $1`,
		id,
	).Scan(&space.SpaceType, &space.CreatedAt, &space.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	err = loadSpaceContent(conn, auth, []*Space{&space}, nil, 0, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	return &space, nil

}

func LoadSubspacesByCheckinTotal(conn *sql.DB, auth *ajax.Auth,
	parentID *uint, // optional
	offset uint, limit uint, // pagination
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

	rows, err := conn.Query(`SELECT space.id,
		space.space_type, space.created_at, space.created_by,
		`+bookmarkFieldSql+`, space.overall_checkin_total
		FROM space
		WHERE space.space_type NOT IN (
			'`+db.Arg(&args, SpaceTypeTitle)+`',
			'`+db.Arg(&args, SpaceTypeTag)+`'
		)
		`+parentClauseSql+`
		ORDER BY space.overall_checkin_total DESC
		OFFSET `+db.Arg(&args, offset)+`
		LIMIT `+db.Arg(&args, limit),
		args...,
	)

	if err == sql.ErrNoRows {
		return spaces, nil
	} else if err != nil {
		return nil, fmt.Errorf("loading top spaces: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var space = Space{}
		err = rows.Scan(&space.ID, &space.SpaceType, &space.CreatedAt, &space.CreatedBy,
			&space.UserBookmark, &space.CheckinTotal)
		if err != nil {
			return nil, fmt.Errorf("loading top spaces: %w", err)
		}
		spaces = append(spaces, &space)
	}

	err = loadSpaceContent(conn, auth, spaces, nil, 0, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	if auth != nil {
		err = loadLastUserTitleForSpaces(conn, *auth, spaces)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
	}

	err = loadTopTitleForSpaces(conn, spaces, 0, MaxSubspacesPageLimit)
	if err != nil {
		return nil, fmt.Errorf("loading top titles: %w", err)
	}

	return spaces, nil

}

// --------------------------------------------------
// batch load functions

func loadSpaceContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration,
	loadSubspaces bool,
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
		loadCheckinSpaceDetails(conn, auth, checkinSpaces, loadSubspaces)
	}

	if hasSpacesOfType(spaces, SpaceTypeStream) {
		var streamSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeStream {
				streamSpaces = append(streamSpaces, space)
			}
		}
		loadStreamSpaceDetails(conn, auth, streamSpaces)
	}

	return nil

}

func loadLastUserTitleForSpaces(conn *sql.DB, auth ajax.Auth, spaces []*Space) error {
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

	rows, err := conn.Query(`SELECT space.id, spare.parent_id, unique_text.text
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
		AND EXISTS(SELECT 1 FROM user_space_bookmark
			WHERE user_space_bookmark.user_id = `+db.Arg(&args, auth.UserID)+`
			AND user_space_bookmark.space_id = space.id)
		AND space.parent_id IN (`+inClauseSql+`)
		GROUP BY space.id, space.parent_id
		ORDER BY user_space_bookmark.created_at DESC
		LIMIT 1`,
		args...,
	)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading last bookmarked titles: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var parentID uint
		var text string
		err = rows.Scan(&spaceID, &parentID, &text)
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
				space.LastUserTitle = &title
				break
			}
		}
	}

	return nil

}

func loadTopTitleForSpaces(conn *sql.DB, spaces []*Space, offset, limit uint) error {
	// Load the top title for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	if limit > MaxSubspacesPageLimit {
		limit = MaxSubspacesPageLimit
	}

	var inClauseSql string

	var args = []interface{}{}

	for i, space := range spaces {
		if i > 0 {
			inClauseSql += `, `
		}
		inClauseSql += db.Arg(&args, space.ID)
	}

	rows, err := conn.Query(`SELECT space.id, spare.parent_id, unique_text.text
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
		AND space.parent_id IN (`+inClauseSql+`)
		GROUP BY space.id
		ORDER BY space.overall_checkin_total DESC
		OFFSET `+db.Arg(&args, offset)+`
		LIMIT `+db.Arg(&args, limit),
		args...,
	)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading top titles: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var parentID uint
		var text string
		err = rows.Scan(&spaceID, &parentID, &text)
		if err != nil {
			return fmt.Errorf("loading top titles: %w", err)
		}
		for _, space := range spaces {
			if space.ID == parentID {
				var title = &Space{
					ID:        spaceID,
					SpaceType: SpaceTypeTitle,
					Text:      &text,
				}
				var titles = []*Space{title}
				space.TopTitles = &titles
			}
		}
	}

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

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading title spaces content: %w", err)
	}

	defer rows.Close()

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

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading tag spaces content: %w", err)
	}

	defer rows.Close()

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

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading text spaces content: %w", err)
	}

	defer rows.Close()

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

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading title spaces content: %w", err)
	}

	defer rows.Close()

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

func loadCheckinSpaceDetails(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	loadSubSpaces bool,
) error {
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

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading checkin space details: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var checkinSpaceID *uint
		err = rows.Scan(&spaceID, &checkinSpaceID)
		if err != nil {
			return fmt.Errorf("loading checkin space details: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				var space = &Space{
					ID: *checkinSpaceID,
				}
				space.CheckinSpace = &space
			}
		}
	}

	var checkinSpaces []*Space
	for _, space := range spaces {
		if space.CheckinSpace != nil {
			checkinSpaces = append(checkinSpaces, *space.CheckinSpace)
		}
	}

	err = loadSpaceContent(conn, auth, checkinSpaces, nil, 0, false)
	if err != nil {
		return fmt.Errorf("loading checkin space details: %w", err)
	}

	return nil

}

func loadStreamSpaceDetails(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
) error {
	// Load stream-of-conscoiusness content for multiple spaces

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
		space.id, stream_of_conscioussness_space.closed_at
		FROM space
		INNER JOIN stream_of_conscioussness_space ON stream_of_conscioussness_space.space_id = space.id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading stream space details: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var closedAt *time.Time
		err = rows.Scan(&spaceID, &closedAt)
		if err != nil {
			return fmt.Errorf("loading checkin space details: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.StreamClosedAt = &closedAt
			}
		}
	}

	// load stream texts
	// for _, space := range spaces {
	// 	var streamTexts []*Space
	// 	var limitTextCreatedAtQuery string
	// 	if space.StreamClosedAt != nil {
	// 		// Up to time closed
	// 		limitTextCreatedAtQuery = ``
	// 	} else {
	// 		// Up to now
	// 		limitTextCreatedAtQuery = ``
	// 	}
	// }

	return nil

}
