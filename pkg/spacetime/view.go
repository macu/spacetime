package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

const MaxSubspacesPageLimit = 20
const DefaultTitlesLimit = 5
const DefaultTagsLimit = 10

func LoadSpace(conn *sql.DB, auth *ajax.Auth, id uint) (*Space, error) {
	spaces, err := LoadSpaces(conn, auth, []uint{id})
	if err != nil {
		return nil, fmt.Errorf("loading space: %w", err)
	}
	if len(spaces) == 0 {
		return nil, nil
	}
	return spaces[0], nil
}

func LoadSpaces(conn *sql.DB, auth *ajax.Auth, ids []uint) ([]*Space, error) {
	// Load a single space (header details) and its associated content

	var spaces = []*Space{}

	var args = []interface{}{}

	var bookmarkFieldSql string
	if auth != nil {
		bookmarkFieldSql = `EXISTS(SELECT * FROM user_bookmark
			WHERE user_bookmark.space_id=space.id
			AND user_bookmark.user_id = ` + db.Arg(&args, auth.UserID) + `
			) AS user_bookmark`
	} else {
		bookmarkFieldSql = `FALSE AS user_bookmark`
	}

	rows, err := conn.Query(`SELECT space.id, space.parent_id, space.space_type,
		space.created_at, space.created_by,
		unique_text.text_value AS label,
		user_account.handle, user_account.display_name,
		CASE WHEN user_space_config.space_id IS NULL THEN FALSE ELSE TRUE END AS user_pinned,
		`+bookmarkFieldSql+`,
		(SELECT link_space.link_space_id FROM link_space
			WHERE link_space.space_id = space.id
			LIMIT 1) AS link_space_id
		FROM space
		LEFT JOIN branch_space ON space.id = branch_space.space_id
		LEFT JOIN unique_text ON unique_text.id = branch_space.label_text_id
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN user_space_config ON user_space_config.space_id = space.id
			WHERE `+db.In("space.id", &args, ids),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var space = &Space{}
		err = rows.Scan(&space.ID, &space.ParentID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy,
			&space.Label,
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.IsPinned,
			&space.UserBookmark,
			&space.LinkSpaceID,
		)
		if err != nil {
			return nil, fmt.Errorf("loading space details: %w", err)
		}
		spaces = append(spaces, space)
	}
	if len(spaces) == 0 {
		return nil, nil
	}

	err = LoadSpaceContent(conn, auth, spaces, true)
	if err != nil {
		return nil, fmt.Errorf("loading space content: %w", err)
	}

	return spaces, nil

}

func LoadParentPath(conn *sql.DB, auth *ajax.Auth, id uint) ([]*Space, error) {

	// recursively load space details following parent_id of space with given id
	// until reaching the root space

	// return array has root space first

	var spaces = []*Space{}

	var args = []interface{}{}

	rows, err := conn.Query(`WITH RECURSIVE parent_spaces AS (
		SELECT space.id, space.parent_id, space.space_type,
			space.created_at, space.created_by,
			unique_text.text_value AS label,
			user_account.handle, user_account.display_name,
			CASE WHEN user_space_config.space_id IS NULL THEN FALSE ELSE TRUE END AS user_pinned
		FROM space
		LEFT JOIN branch_space ON space.id = branch_space.space_id
		LEFT JOIN unique_text ON unique_text.id = branch_space.label_text_id
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN user_space_config ON user_space_config.space_id = space.id
		WHERE space.id = `+db.Arg(&args, id)+`
		UNION ALL
		SELECT space.id, space.parent_id, space.space_type,
			space.created_at, space.created_by,
			unique_text.text_value AS label,
			user_account.handle, user_account.display_name,
			CASE WHEN user_space_config.space_id IS NULL THEN FALSE ELSE TRUE END AS user_pinned
		FROM space
		INNER JOIN parent_spaces ON parent_spaces.parent_id = space.id
		LEFT JOIN branch_space ON space.id = branch_space.space_id
		LEFT JOIN unique_text ON unique_text.id = branch_space.label_text_id
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN user_space_config ON user_space_config.space_id = space.id
	)
	SELECT * FROM parent_spaces`, args...)
	if err != nil {
		return nil, fmt.Errorf("loading parent path: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var space = &Space{}
		err = rows.Scan(&space.ID, &space.ParentID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy,
			&space.Label,
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.IsPinned,
		)
		if err != nil {
			return nil, fmt.Errorf("loading parent path: %w", err)
		}
		spaces = append(spaces, space)
	}

	// reverse the order of spaces to have root space first
	for i, j := 0, len(spaces)-1; i < j; i, j = i+1, j-1 {
		spaces[i], spaces[j] = spaces[j], spaces[i]
	}

	if hasSpacesOfType(spaces, SpaceTypeTag) {
		err := loadTagSpacesContent(conn,
			extractSpacesByType(spaces, SpaceTypeTag))
		if err != nil {
			return nil, err
		}
	}

	return spaces, nil

}

func LoadSubspaces(conn *sql.DB, auth *ajax.Auth,
	parentID *uint, // optional
	offset uint, limit uint, // pagination
	filter *SpaceFilter, // optional filter by date
	typesFilter *TypesFilter, // optional filter by space type
) ([]*Space, error) {

	var spaces = []*Space{}

	var args = []interface{}{}

	var bookmarkFieldSql string
	if auth != nil {
		bookmarkFieldSql = `EXISTS(SELECT 1 FROM user_bookmark
			WHERE user_bookmark.user_id = ` + db.Arg(&args, auth.UserID) + `
			AND user_bookmark.space_id = space.id) AS user_bookmark`
	} else {
		bookmarkFieldSql = `FALSE AS user_bookmark`
	}

	var parentClauseSql string
	if parentID != nil {
		parentClauseSql = `space.parent_id = ` + db.Arg(&args, *parentID)
	} else {
		parentClauseSql = `space.parent_id IS NULL`
	}

	var typesClauseSql string
	if typesFilter != nil {
		if len(typesFilter.Types) > 0 {
			if typesFilter.Exclude {
				typesClauseSql = `AND space.space_type NOT IN (`
			} else {
				typesClauseSql = `AND space.space_type IN (`
			}
			for i, spaceType := range typesFilter.Types {
				if i > 0 {
					typesClauseSql += `, `
				}
				typesClauseSql += db.Arg(&args, spaceType)
			}
			typesClauseSql += `)`
		}
	}

	var filterModeClauseSql string
	var orderByClauseSql string
	if filter != nil {

		switch filter.Mode {

		case SpaceFilterModeTopSubspaces:
			orderByClauseSql = `checkin_count DESC, space.created_at DESC`

		case SpaceFilterModeMostRecent:
			orderByClauseSql = `space.created_at DESC`

		case SpaceFilterModePinned:
			filterModeClauseSql = `AND EXISTS (SELECT 1 FROM user_space_config
				WHERE user_space_config.space_id = space.id)`
			orderByClauseSql = `user_space_config.order_number ASC`

		default:
			return nil, fmt.Errorf("invalid filter mode: %s", filter.Mode)
		}

		if filter.PinnedFirst && filter.Mode != SpaceFilterModePinned {
			orderByClauseSql = `CASE WHEN user_space_config.space_id IS NULL THEN 1 ELSE 0 END, ` +
				`user_space_config.order_number ASC, ` + orderByClauseSql
		}

		if filter.Date != nil {
			// up to given date
			filterModeClauseSql = `AND space.created_at <= ` + db.Arg(&args, *filter.Date)
		}
		if filter.Window != nil {
			// within given window
			var windowDuration string
			switch *filter.Window {
			case "day":
				windowDuration = "1 DAY"
			case "week":
				windowDuration = "1 WEEK"
			case "month":
				windowDuration = "1 MONTH"
			case "year":
				windowDuration = "1 YEAR"
			default:
				return nil, fmt.Errorf("invalid filter window: %s", *filter.Window)
			}
			filterModeClauseSql += ` AND space.created_at >= NOW() - INTERVAL '` + windowDuration + `'`
		}

	} else {
		orderByClauseSql = `space.created_at DESC`
	}

	rows, err := conn.Query(`SELECT space.id,
		space.space_type, space.created_at, space.created_by,
		unique_text.text_value AS label,
		user_account.handle, user_account.display_name,
		`+bookmarkFieldSql+`,
		EXISTS(SELECT 1 FROM user_space_config
			WHERE user_space_config.space_id = space.id
		) AS is_pinned,
		(SELECT COUNT(*) FROM checkin
			WHERE checkin.space_id = space.id) AS checkin_count,
		(SELECT link_space.link_space_id FROM link_space
			WHERE link_space.space_id = space.id
			LIMIT 1) AS link_space_id
		FROM space
		LEFT JOIN branch_space ON branch_space.space_id = space.id
		LEFT JOIN unique_text ON unique_text.id = branch_space.label_text_id
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN user_space_config ON user_space_config.space_id = space.id
		WHERE `+parentClauseSql+` `+typesClauseSql+` `+filterModeClauseSql+`
		ORDER BY `+orderByClauseSql+`
		LIMIT `+db.Arg(&args, limit)+`
		OFFSET `+db.Arg(&args, offset),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf("loading top spaces: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var space = Space{
			ParentID: parentID,
		}
		err = rows.Scan(&space.ID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy,
			&space.Label,
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.UserBookmark,
			&space.IsPinned,
			&space.CheckinCount,
			&space.LinkSpaceID,
		)
		if err != nil {
			return nil, fmt.Errorf("loading top spaces: %w", err)
		}
		spaces = append(spaces, &space)
	}

	err = LoadSpaceContent(conn, auth, spaces, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	return spaces, nil

}

// --------------------------------------------------
// batch load functions

func LoadCheckinCount(conn *sql.DB, spaces []*Space) error {

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

	rows, err := conn.Query(`SELECT space.id,
		(SELECT COUNT(*) FROM checkin
			WHERE checkin.space_id = space.id) AS checkin_count
		FROM space
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
		return fmt.Errorf("loading check-in count: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var checkinCount uint
		err = rows.Scan(&spaceID, &checkinCount)
		if err != nil {
			return fmt.Errorf("loading check-in count: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.CheckinCount = checkinCount
				break
			}
		}
	}

	return nil

}

/*
func LoadSubspaceCount(conn *sql.DB, spaces []*Space, filter *SpaceFilter) error {
	// Load subspace count for multiple spaces

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

	rows, err := conn.Query(`SELECT space.id,
		(SELECT COUNT(*) FROM space AS subspace
			WHERE subspace.parent_id = space.id) AS subspace_count
		FROM space
		LEFT JOIN space AS subspace ON subspace.parent_id = space.id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
		return fmt.Errorf("loading subspace count: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var subspaceCount uint
		err = rows.Scan(&spaceID, &subspaceCount)
		if err != nil {
			return fmt.Errorf("loading subspace count: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.TotalSubspaces = subspaceCount
				break
			}
		}
	}

	return nil

}
*/

func LoadSpaceContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	loadLinkedSpaces bool, // prevent recursion
) error {
	// Load content for multiple spaces

	if hasSpacesOfType(spaces, SpaceTypeTag) {
		err := loadTagSpacesContent(conn,
			extractSpacesByType(spaces, SpaceTypeTag))
		if err != nil {
			return err
		}
	}

	if hasSpacesOfType(spaces, SpaceTypeText) {
		err := loadTextSpacesContent(conn,
			extractSpacesByType(spaces, SpaceTypeText))
		if err != nil {
			return err
		}
	}

	if loadLinkedSpaces && hasSpacesOfType(spaces, SpaceTypeLink) {
		linkSpaces := extractSpacesByType(spaces, SpaceTypeLink)
		var linkedSpaceIds []uint
		for _, space := range linkSpaces {
			linkedSpaceIds = append(linkedSpaceIds, **space.LinkSpaceID)
		}
		linkedSpaces, err := LoadSpaces(conn, auth, linkedSpaceIds)
		if err != nil {
			return err
		}
		for _, space := range linkSpaces {
			space.LinkSpace = nil
			for _, linkedSpace := range linkedSpaces {
				if **space.LinkSpaceID == linkedSpace.ID {
					space.LinkSpace = &linkedSpace

					parentPath, err := LoadParentPath(conn, auth, linkedSpace.ID)
					if err != nil {
						return err
					}
					linkedSpace.ParentPath = &parentPath

					break
				}
			}
		}
	}

	// if hasSpacesOfType(spaces, SpaceTypeStream) {
	// 	err := loadStreamSpaceDetails(conn,
	// 		extractSpacesByType(spaces, SpaceTypeStream))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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

func extractSpacesByType(spaces []*Space, spaceType string) []*Space {
	// Extract spaces of a certain type from a list of spaces

	var extractedSpaces = []*Space{}

	for _, space := range spaces {
		if space.SpaceType == spaceType {
			extractedSpaces = append(extractedSpaces, space)
		}
	}

	return extractedSpaces

}

func loadTagSpacesContent(conn *sql.DB, spaces []*Space) error {
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
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text_value
		FROM space
		INNER JOIN tag_space ON tag_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = tag_space.text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
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

func loadTextSpacesContent(conn *sql.DB, spaces []*Space) error {
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
	}

	rows, err := conn.Query(`SELECT
		space.id, text_unique_text.text_value, title_unique_text.text_value,
		CASE WHEN text_space.recording IS NOT NULL THEN TRUE ELSE FALSE END AS has_recording
		FROM space
		INNER JOIN text_space ON text_space.space_id = space.id
		INNER JOIN unique_text AS text_unique_text ON text_unique_text.id = text_space.text_id
		LEFT JOIN unique_text AS title_unique_text ON title_unique_text.id = text_space.title_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
		return fmt.Errorf("loading text spaces content: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var spaceID uint
		var text string
		var title *string
		var hasRecording bool
		err = rows.Scan(&spaceID, &text, &title, &hasRecording)
		if err != nil {
			return fmt.Errorf("loading text spaces content: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				space.Text = &text
				space.Title = &title
				space.HasRecording = &hasRecording
			}
		}
	}

	return nil

}

func loadStreamSpaceDetails(conn *sql.DB, spaces []*Space) error {
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

	if err != nil {
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
