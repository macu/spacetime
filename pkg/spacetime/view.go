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
const DefaultTitlesLimit = 5

func LoadSpace(conn *sql.DB, auth *ajax.Auth, id uint, includeSubspaces bool) (*Space, error) {
	// Load a single space

	var space = Space{
		ID: id,
	}

	err := conn.QueryRow(`SELECT
		space.space_type, space.created_at, space.created_by,
		user_account.handle, user_account.display_name,
		EXISTS(SELECT * FROM user_space_bookmark
			WHERE user_space_bookmark.space_id=space.id
			AND user_space_bookmark.user_id = $2) AS user_bookmark,
		COUNT(subspace.id) AS subspace_count
		FROM space
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN space AS subspace ON subspace.parent_id = space.id
		WHERE space.id = $1
		GROUP BY space.space_type, space.created_at, space.created_by, user_account.handle, user_account.display_name, user_bookmark
		LIMIT 1`,
		id, auth.UserID,
	).Scan(&space.SpaceType, &space.CreatedAt, &space.CreatedBy,
		&space.AuthorHandle, &space.AuthorDisplayName,
		&space.UserBookmark,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	err = loadSpaceContent(conn, auth, []*Space{&space}, nil, nil, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	if auth != nil {
		err = loadBookmarkedTitles(conn, *auth, []*Space{&space})
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
	}

	err = loadTopTitles(conn, []*Space{&space}, 0, DefaultTitlesLimit)
	if err != nil {
		return nil, fmt.Errorf("loading top titles: %w", err)
	}

	if includeSubspaces {
		content, err := LoadTopSubspaces(conn, auth,
			&id, 0, MaxSubspacesPageLimit)
		if err != nil {
			return nil, fmt.Errorf("loading subspaces: %w", err)
		}
		space.TopSubspaces = &content
	}

	return &space, nil

}

func LoadTopSubspaces(conn *sql.DB, auth *ajax.Auth,
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
		parentClauseSql = `space.parent_id = ` + db.Arg(&args, *parentID)
	} else {
		parentClauseSql = `space.parent_id IS NULL`
	}

	rows, err := conn.Query(`SELECT space.id,
		space.space_type, space.created_at, space.created_by,
		user_account.handle, user_account.display_name,
		`+bookmarkFieldSql+`,
		COUNT(subspace.id) AS subspace_count
		FROM space
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN spacew AS subspace ON subspace.parent_id = space.id
		WHERE `+parentClauseSql+`
		ORDER BY SELECT COUNT(*) FROM space AS subspace
		WHERE subspace.parent_id = space.id
		GROUP BY space.id, space.space_type, space.created_at, space.created_by, user_account.handle, user_account.display_name, user_bookmark
		ORDER BY subspace_count DESC
		LIMIT `+db.Arg(&args, limit)+`
		OFFSET `+db.Arg(&args, offset),
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
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.UserBookmark,
			&space.TotalSubspaces,
		)
		if err != nil {
			return nil, fmt.Errorf("loading top spaces: %w", err)
		}
		spaces = append(spaces, &space)
	}

	err = loadSpaceContent(conn, auth, spaces, nil, nil, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	if auth != nil {
		err = loadBookmarkedTitles(conn, *auth, spaces)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
	}

	err = loadTopTitles(conn, spaces, 0, DefaultTitlesLimit)
	if err != nil {
		return nil, fmt.Errorf("loading top titles: %w", err)
	}

	return spaces, nil

}

// --------------------------------------------------
// batch load functions

func loadBookmarkedTitles(conn *sql.DB, auth ajax.Auth,
	spaces []*Space,
) error {

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
	}

	rows, err := conn.Query(`SELECT space.id, space.parent_id, unique_text.text
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.space_type = `+db.Arg(&args, SpaceTypeTitle)+`
		AND EXISTS(SELECT 1 FROM user_space_bookmark
			WHERE user_space_bookmark.user_id = `+db.Arg(&args, auth.UserID)+`
			AND user_space_bookmark.space_id = space.id)
		AND space.parent_id IN (`+inClauseSql+`)
		GROUP BY space.id, space.parent_id
		ORDER BY user_space_bookmark.created_at DESC`,
		args...,
	)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("loading bookmarked titles: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var parentID uint
		var text string
		err = rows.Scan(&spaceID, &parentID, &text)
		if err != nil {
			return fmt.Errorf("loading bookmarked titles: %w", err)
		}
		for _, space := range spaces {
			if space.ID == parentID {
				var title = &Space{
					ID:        spaceID,
					SpaceType: SpaceTypeTitle,
					Text:      &text,
				}
				if space.UserTitles == nil {
					space.UserTitles = &[]*Space{title}
				} else {
					*space.UserTitles = append(*space.UserTitles, title)
				}
				break
			}
		}
	}

	return nil

}

func loadTopTitles(conn *sql.DB, spaces []*Space,
	offset, limit uint,
) error {

	if len(spaces) == 0 {
		return nil
	}

	if limit > MaxSubspacesPageLimit {
		limit = MaxSubspacesPageLimit
	}

	for _, space := range spaces {

		rows, err := conn.Query(`SELECT space.id, space.parent_id, unique_text.text
			FROM space
			INNER JOIN title_space ON title_space.space_id = space.id
			INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
			WHERE space.space_type = $1
			AND space.parent_id = $2
			GROUP BY space.id, space.parent_id
			ORDER BY space.overall_checkin_total DESC
			OFFSET $3
			LIMIT $4`,
			SpaceTypeTitle, space.ID, offset, limit,
		)

		if err == sql.ErrNoRows {
			space.TopTitles = &[]*Space{}
			continue
		} else if err != nil {
			return fmt.Errorf("loading top titles: %w", err)
		}

		defer rows.Close()

		var titles = []*Space{}

		for rows.Next() {
			var spaceID uint
			var parentID uint
			var text string
			err = rows.Scan(&spaceID, &parentID, &text)
			if err != nil {
				return fmt.Errorf("loading top titles: %w", err)
			}
			var title = &Space{
				ID:        spaceID,
				SpaceType: SpaceTypeTitle,
				Text:      &text,
			}
			titles = append(titles, title)
		}

		space.TopTitles = &titles

	}

	return nil

}

func loadSpaceContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval *time.Duration,
	loadCheckinSpace bool,
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

	if hasSpacesOfType(spaces, SpaceTypeCheckin) && loadCheckinSpace {
		var checkinSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeCheckin {
				checkinSpaces = append(checkinSpaces, space)
			}
		}
		loadCheckinSpaceDetails(conn, auth, checkinSpaces)
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
		var replayData NakedText
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

	rows, err := conn.Query(`SELECT space.id,
		check_space.checkin_space_id,
		checked_space.space_type, checked_space.created_at, checked_space.created_by
		FROM space
		INNER JOIN checkin_space check_space ON check_space.space_id = space.id
		INNER JOIN space AS checked_space ON checked_space.id = check_space.checkin_space_id
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
		var checkedSpaceType string
		var checkedSpaceCreatedAt time.Time
		var checkedSpaceCreatedBy uint
		err = rows.Scan(&spaceID, &checkinSpaceID,
			&checkedSpaceType, &checkedSpaceCreatedAt, &checkedSpaceCreatedBy)
		if err != nil {
			return fmt.Errorf("loading checkin space details: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				var space = &Space{
					ID:        *checkinSpaceID,
					SpaceType: checkedSpaceType,
					CreatedAt: checkedSpaceCreatedAt,
					CreatedBy: checkedSpaceCreatedBy,
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

	err = loadSpaceContent(conn, auth, checkinSpaces, nil, nil, false)
	if err != nil {
		return fmt.Errorf("loading checkin space details: %w", err)
	}

	return nil

}

func loadStreamSpaceDetails(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
) error {
	// Load stream content for multiple spaces

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

	rows, err := conn.Query(`SELECT space.id,
		stream_space.stream_closed_at
		FROM space
		INNER JOIN stream_space ON stream_space.space_id = space.id
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
		var streamClosedAt *time.Time
		err = rows.Scan(&spaceID, &streamClosedAt)
		if err != nil {
			return fmt.Errorf("loading stream space details: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.StreamClosedAt = &streamClosedAt
			}
		}
	}

	return nil

}
