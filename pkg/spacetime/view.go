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
const DefaultTagsLimit = 10

func LoadSpace(conn *sql.DB, auth *ajax.Auth, id uint) (*Space, error) {
	// Load a single space (header details) and its associated content

	var space = Space{
		ID: id,
	}

	var args = []interface{}{}

	var bookmarkFieldSql string
	if auth != nil {
		bookmarkFieldSql = `EXISTS(SELECT * FROM user_space_bookmark
			WHERE user_space_bookmark.space_id=space.id
			AND user_space_bookmark.user_id = ` + db.Arg(&args, auth.UserID) + `
			) AS user_bookmark`
	} else {
		bookmarkFieldSql = `FALSE AS user_bookmark`
	}

	err := conn.QueryRow(`SELECT space.parent_id, space.space_type,
		space.created_at, space.created_by,
		user_account.handle, user_account.display_name,
		`+bookmarkFieldSql+`
		FROM space
		LEFT JOIN user_account ON user_account.id = space.created_by
		WHERE space.id = `+db.Arg(&args, id)+`
		LIMIT 1`,
		args...,
	).Scan(&space.ParentID, &space.SpaceType,
		&space.CreatedAt, &space.CreatedBy,
		&space.AuthorHandle, &space.AuthorDisplayName,
		&space.UserBookmark,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	err = LoadSpaceContent(conn, auth, []*Space{&space}, true)
	if err != nil {
		return nil, fmt.Errorf("loading space content: %w", err)
	}

	return &space, nil

}

func LoadParentPath(conn *sql.DB, auth *ajax.Auth, id uint) ([]*Space, error) {

	// recursively load space details following parent_id of space with given id
	// until reaching the root space

	// TODO Use recusive query

	var spaces = []*Space{}

	for {
		space, err := LoadSpace(conn, auth, id)
		if err != nil {
			return nil, fmt.Errorf("loading parent path: %w", err)
		}
		if space == nil {
			break
		}

		spaces = append([]*Space{space}, spaces...)

		if space.ParentID == nil {
			break
		}

		id = *space.ParentID
	}

	if hasSpacesOfType(spaces, SpaceTypeTitle) {
		err := loadTitleSpacesContent(conn,
			extractSpacesByType(spaces, SpaceTypeTitle))
		if err != nil {
			return nil, err
		}
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

func LoadTopSubspaces(conn *sql.DB, auth *ajax.Auth,
	parentID *uint, // optional
	offset uint, limit uint, // pagination
	date *time.Time, interval *time.Duration, // review
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
		LEFT JOIN space AS subspace ON subspace.parent_id = space.id
		WHERE `+parentClauseSql+`
		GROUP BY space.id, space.space_type, space.created_at, space.created_by, user_account.handle, user_account.display_name, user_bookmark
		ORDER BY subspace_count DESC, space.created_at DESC
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
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.UserBookmark,
			&space.TotalSubspaces,
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

func LoadSubspaceCount(conn *sql.DB, spaces []*Space) error {
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
		COUNT(subspace.id) AS subspace_count
		FROM space
		LEFT JOIN space AS subspace ON subspace.parent_id = space.id
		WHERE space.id IN (`+inClauseSql+`)
		GROUP BY space.id`,
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

func LoadSpaceContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	loadLinkedSpaces bool, // prevent recursion
) error {
	// Load content for multiple spaces

	if hasSpacesOfType(spaces, SpaceTypeTitle) {
		err := loadTitleSpacesContent(conn,
			extractSpacesByType(spaces, SpaceTypeTitle))
		if err != nil {
			return err
		}
	}

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

	if hasSpacesOfType(spaces, SpaceTypeNaked) {
		err := loadNakedTextSpacesContent(conn,
			extractSpacesByType(spaces, SpaceTypeNaked))
		if err != nil {
			return err
		}
	}

	if loadLinkedSpaces && hasSpacesOfType(spaces, SpaceTypeLink) {
		err := loadLinkSpaceDetails(conn, auth,
			extractSpacesByType(spaces, SpaceTypeLink))
		if err != nil {
			return err
		}
	}

	if hasSpacesOfType(spaces, SpaceTypeStream) {
		err := loadStreamSpaceDetails(conn,
			extractSpacesByType(spaces, SpaceTypeStream))
		if err != nil {
			return err
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

func loadTitleSpacesContent(conn *sql.DB, spaces []*Space) error {
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
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text_value
		FROM space
		INNER JOIN title_space ON title_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = title_space.unique_text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
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
		INNER JOIN unique_text ON unique_text.id = tag_space.unique_text_id
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
		space.id, unique_text.text_value
		FROM space
		INNER JOIN text_space ON text_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = text_space.unique_text_id
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

func loadNakedTextSpacesContent(conn *sql.DB, spaces []*Space) error {
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
	}

	rows, err := conn.Query(`SELECT
		space.id, unique_text.text_value, naked_text_space.replay_data
		FROM space
		INNER JOIN naked_text_space ON naked_text_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = naked_text_space.final_unique_text_id
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
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

func loadLinkSpaceDetails(conn *sql.DB, auth *ajax.Auth, spaces []*Space) error {
	// Load link content for multiple spaces

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

		// Set link space to nil for direct check-ins
		var nullId *uint = nil
		var nullSpace *Space = nil
		space.LinkSpaceID = &nullId
		space.LinkSpace = &nullSpace
	}

	var bookmarkFieldSql string
	if auth != nil {
		bookmarkFieldSql = `EXISTS(SELECT * FROM user_space_bookmark
			WHERE linked_space.id IS NOT NULL
			AND user_space_bookmark.space_id = linked_space.id
			AND user_space_bookmark.user_id = ` + db.Arg(&args, auth.UserID) + `
		) AS user_bookmark`
	} else {
		bookmarkFieldSql = `FALSE AS user_bookmark`
	}

	rows, err := conn.Query(`SELECT space.id,
		link_space.link_space_id,
		linked_space.space_type, linked_space.created_at, linked_space.created_by,
		user_account.handle, user_account.display_name,
		`+bookmarkFieldSql+`
		FROM space
		INNER JOIN link_space ON link_space.space_id = space.id
		LEFT JOIN space AS linked_space ON linked_space.id = link_space.link_space_id
		LEFT JOIN user_account ON user_account.id = linked_space.created_by
		WHERE space.id IN (`+inClauseSql+`)`,
		args...,
	)

	if err != nil {
		return fmt.Errorf("loading link space details: %w", err)
	}

	defer rows.Close()

	var linkSpaces = []*Space{}

	for rows.Next() {
		var spaceID uint
		var linkSpace = &Space{}
		err = rows.Scan(&spaceID,
			&linkSpace.ID, &linkSpace.SpaceType,
			&linkSpace.CreatedAt, &linkSpace.CreatedBy,
			&linkSpace.AuthorHandle, &linkSpace.AuthorDisplayName,
			&linkSpace.UserBookmark,
		)
		if err != nil {
			return fmt.Errorf("loading link space details: %w", err)
		}
		for _, space := range spaces {
			if space.ID == spaceID {
				if linkSpace.ID != 0 {
					var id = &linkSpace.ID
					space.LinkSpaceID = &id
					space.LinkSpace = &linkSpace
					linkSpaces = append(linkSpaces, linkSpace)
				}
			}
		}
	}

	err = LoadSpaceContent(conn, auth, linkSpaces, false)
	if err != nil {
		return fmt.Errorf("loading link space details: %w", err)
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
