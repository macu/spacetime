package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
)

const SubspacesPageLimit = 20

func LoadSpaceTopCheckinsAllTime(conn *sql.DB, auth *ajax.Auth,
	spaceID uint) (*Space, error) {
	// Load a space and its subspaces ordered by all time checkin count

	var space = Space{
		ID: spaceID,
	}

	var args = []interface{}{spaceID}

	var bookmarkQuery string

	if auth != nil {
		bookmarkQuery = `EXISTS(SELECT 1 FROM user_space_bookmark
			WHERE user_space_bookmark.user_id = $2
			AND user_space_bookmark.space_id = space.id) AS user_bookmark`
		args = append(args, auth.UserID)
	} else {
		bookmarkQuery = `FALSE AS user_bookmark`
	}

	// Load space
	err := conn.QueryRow(`SELECT space.space_type, space.created_at, space.created_by,
		`+bookmarkQuery+`
		FROM space WHERE space.id = $1`,
		args...,
	).Scan(&space.SpaceType, &space.CreatedAt, &space.CreatedBy, &space.UserBookmark)
	if err != nil {
		return nil, fmt.Errorf("loading node class: %w", err)
	}

	// Load space content
	loadSpaceDetails(conn, auth, []*Space{&space}, nil, 0)

	// Load user's bookmarked titles
	if auth != nil {
		bookmarkedTitles, err := loadBookmarkedTitles(conn, *auth, spaceID)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
		space.BookmarkedTitles = &bookmarkedTitles
	}

	// Load all-time top titles

	// Load all-time top tags

	// Load top content
	// Load current top content

	return &space, nil

}

func LoadSpaceMostRecentCheckinsByDate(conn *sql.DB, auth *ajax.Auth,
	spaceID uint, date *time.Time) {
	// Load most recent checkins up to date
}

func LoadSpaceMostRecentUserCheckinsByDate(conn *sql.DB, auth *ajax.Auth,
	spaceID uint, userID uint, date *time.Time) {
	// Load a specific user's checkins ordered by most recent up to date
}

// --------------------------------------------------
// bath load content from joined tables

func loadBookmarkedTitles(conn *sql.DB, auth ajax.Auth, spaceID uint) ([]*Space, error) {
	rows, err := conn.Query(`SELECT space.id, space.space_type, space.created_at, space.created_by
		FROM space
		JOIN user_space_bookmark ON user_space_bookmark.space_id = space.id
		WHERE user_space_bookmark.user_id = $1
		AND user_space_bookmark.space_id = $2
		AND space.space_type = $3
		ORDER BY user_space_bookmark.created_at DESC
		LIMIT $4`,
		auth.UserID, spaceID, SpaceTypeTitle, SubspacesPageLimit,
	)

	if err != nil {
		return nil, fmt.Errorf("loading bookmarked titles: %w", err)
	}

	var spaces = []*Space{}
	for rows.Next() {
		var space = Space{}
		err = rows.Scan(&space.ID, &space.SpaceType, &space.CreatedAt, &space.CreatedBy)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarked titles: %w", err)
		}
		spaces = append(spaces, &space)
	}

	loadSpaceDetails(conn, &auth, spaces, nil, 0)

	return spaces, nil

}

func loadSpacesContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load content for multiple spaces

	if hasSpaceOfType(spaces, SpaceTypeTitle) {
		var titleSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeTitle {
				titleSpaces = append(titleSpaces, space)
			}
		}
		loadTitleSpacesContent(conn, auth, titleSpaces,
			date, interval)
	}

	if hasSpaceOfType(spaces, SpaceTypeTag) {
		var tagSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeTag {
				tagSpaces = append(tagSpaces, space)
			}
		}
		loadTagSpacesContent(conn, auth, tagSpaces,
			date, interval)
	}

	if hasSpaceOfType(spaces, SpaceTypeText) {
		var textSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeText {
				textSpaces = append(textSpaces, space)
			}
		}
		loadTextSpacesContent(conn, auth, textSpaces,
			date, interval)
	}

	if hasSpaceOfType(spaces, SpaceTypeNaked) {
		var nakedTextSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeNaked {
				nakedTextSpaces = append(nakedTextSpaces, space)
			}
		}
		loadNakedTextSpacesContent(conn, auth, nakedTextSpaces,
			date, interval)
	}

	if hasSpaceOfType(spaces, SpaceTypeCheckin) {
		var checkinSpaces []*Space
		for _, space := range spaces {
			if space.SpaceType == SpaceTypeCheckin {
				checkinSpaces = append(checkinSpaces, space)
			}
		}
		loadSpaceDetails(conn, auth, checkinSpaces,
			date, interval)
	}

}

func hasSpaceOfType(spaces []*Space, spaceType string) bool {
	// Check if a space of a certain type exists in a list of spaces
	for _, space := range spaces {
		if space.SpaceType == spaceType {
			return true
		}
	}
	return false
}

func loadTitleSpacesContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load title content for multiple spaces

}

func loadTagSpacesContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load tag content for multiple spaces

}

func loadTextSpacesContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load text content for multiple spaces

}

func loadNakedTextSpacesContent(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load naked text

}

func loadSpaceDetails(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load checkin content for multiple spaces

	var checkinSpaces []*Space
	for _, space := range spaces {
		if space.CheckinSpaceID != nil {
			var checkinSpace = Space{
				ID: *space.CheckinSpaceID,
			}
			checkinSpaces = append(checkinSpaces, &checkinSpace)
		}
	}

	loadCheckinSpaceTitles(conn, auth, checkinSpaces, date, interval)

	for _, space := range spaces {
		if space.CheckinSpaceID != nil {
			for _, checkinSpace := range checkinSpaces {
				if checkinSpace.ID == *space.CheckinSpaceID {
					space.CheckinSpace = &checkinSpace
				}
			}
		}
	}

}

func loadCheckinSpaceTitles(conn *sql.DB, auth *ajax.Auth,
	spaces []*Space,
	date *time.Time, interval time.Duration) {
	// Load titles for checked in spaces (top title all time)
}
