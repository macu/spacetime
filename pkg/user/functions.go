package user

import (
	"database/sql"
	"fmt"
	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
	"time"
)

func CheckAdmin(db db.DBConn, userID uint) bool {
	var userRole string
	err := db.QueryRow(`SELECT role FROM user_account WHERE id = $1`, userID).Scan(&userRole)
	if err != nil {
		return false
	}
	return userRole == string(RoleAdmin)
}

func BookmarkSpace(db db.DBConn, userID uint, spaceID uint, bookmark bool) error {
	if bookmark {
		_, err := db.Exec(`INSERT INTO user_bookmark
			(user_id, space_id, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, space_id) DO UPDATE SET created_at = $3`,
			userID, spaceID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to bookmark space: %w", err)
		}
	} else {
		_, err := db.Exec(`DELETE FROM user_bookmark
			WHERE user_id = $1 AND space_id = $2`,
			userID, spaceID)
		if err != nil {
			return fmt.Errorf("failed to unbookmark space: %w", err)
		}
	}

	return nil
}

func GetBookmarkedSpaces(conn *sql.DB, auth ajax.Auth,
	offset uint, limit uint,
	includeParentPath bool, includeTags bool,
) ([]*spacetime.Space, error) {

	var args []interface{}
	var spaces = []*spacetime.Space{}

	rows, err := conn.Query(`SELECT space.id,
		space.space_type, space.created_at, space.created_by,
		unique_text.text_value AS label,
		user_account.handle, user_account.display_name,
		TRUE AS bookmarked,
		user_bookmark.created_at AS bookmark_created_at,
		EXISTS(SELECT 1 FROM user_space_config
			WHERE user_space_config.space_id = space.id
		) AS is_pinned,
		(SELECT COUNT(*) FROM checkin
			WHERE checkin.space_id = space.id) AS checkin_count
		FROM space
		INNER JOIN user_bookmark ON user_bookmark.space_id = space.id
		LEFT JOIN branch_space ON branch_space.space_id = space.id
		LEFT JOIN unique_text ON unique_text.id = branch_space.label_text_id
		LEFT JOIN user_account ON user_account.id = space.created_by
		LEFT JOIN user_space_config ON user_space_config.space_id = space.id
		WHERE user_bookmark.user_id = `+db.Arg(&args, auth.UserID)+`
		ORDER BY user_bookmark.created_at DESC
		LIMIT `+db.Arg(&args, limit)+`
		OFFSET `+db.Arg(&args, offset),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf("loading bookmarks: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var space = &spacetime.Space{}
		err = rows.Scan(&space.ID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy,
			&space.Label,
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.UserBookmark, &space.BookmarkCreatedAt,
			&space.IsPinned,
			&space.CheckinCount,
		)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarks: %w", err)
		}
		spaces = append(spaces, space)
	}

	if includeParentPath {
		for _, space := range spaces {
			parentPath, err := spacetime.LoadParentPath(conn, &auth, space.ID)
			if err != nil {
				return nil, fmt.Errorf("loading parent path for space %d: %w", space.ID, err)
			}
			space.ParentPath = &parentPath
		}
	}

	if includeTags {
		err := spacetime.LoadTags(conn, spaces,
			0, spacetime.DefaultTagsLimit,
			&spacetime.SpaceFilter{
				Mode: spacetime.SpaceFilterModeTopSubspaces,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("loading tags for spaces: %w", err)
		}

	}

	err = spacetime.LoadSpaceContent(conn, &auth, spaces, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	return spaces, nil

}
