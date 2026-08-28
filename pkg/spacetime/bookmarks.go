package spacetime

import (
	"database/sql"
	"fmt"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
	"time"
)

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
	includeLinkedInParentId *uint,
) ([]*Space, error) {

	var args []interface{}
	var spaces = []*Space{}

	var includedInParentIdField string
	if includeLinkedInParentId != nil {
		includedInParentIdField = `EXISTS(SELECT 1 FROM link_space
			WHERE link_space.parent_id = ` + db.Arg(&args, *includeLinkedInParentId) + `
			AND link_space.link_space_id = space.id
		) AS included_in_parent_id`
	} else {
		includedInParentIdField = `NULL AS included_in_parent_id`
	}

	rows, err := conn.Query(`SELECT space.id, space.parent_id, space.space_type,
		space.created_at, space.created_by,
		unique_text.text_value AS label,
		user_account.handle, user_account.display_name,
		TRUE AS bookmarked,
		user_bookmark.created_at AS bookmark_created_at,
		EXISTS(SELECT 1 FROM user_space_config
			WHERE user_space_config.space_id = space.id
		) AS is_pinned,
		COALESCE((
			SELECT vote_value FROM space_vote
			WHERE space_vote.user_id = `+db.Arg(&args, auth.UserID)+`
			AND space_vote.space_id = space.id
			AND created_at >= `+db.Arg(&args, time.Now().Add(-VoteWindow))+`
			ORDER BY created_at DESC
			LIMIT 1
		), 0) AS current_vote,
		(SELECT SUM(vote_value) FROM space_vote
			WHERE space_vote.space_id = space.id
		) AS vote_sum,
		`+includedInParentIdField+`,
		(SELECT link_space.link_space_id FROM link_space
			WHERE link_space.space_id = space.id
			LIMIT 1
		) AS link_space_id
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
		var space = &Space{}
		err = rows.Scan(&space.ID, &space.ParentID, &space.SpaceType,
			&space.CreatedAt, &space.CreatedBy,
			&space.Label,
			&space.AuthorHandle, &space.AuthorDisplayName,
			&space.UserBookmark, &space.BookmarkCreatedAt,
			&space.IsPinned,
			&space.CurrentVote,
			&space.VoteSum,
			&space.IncludedInParent,
			&space.LinkSpaceID,
		)
		if err != nil {
			return nil, fmt.Errorf("loading bookmarks: %w", err)
		}
		spaces = append(spaces, space)
	}

	if includeParentPath {
		for _, space := range spaces {
			if space.ParentID == nil {
				space.ParentPath = &[]*Space{}
				continue
			}
			parentPath, err := LoadParentPath(conn, &auth, *space.ParentID)
			if err != nil {
				return nil, fmt.Errorf("loading parent path for space %d: %w", space.ID, err)
			}
			space.ParentPath = &parentPath
		}
	}

	if includeTags {
		err := LoadTags(conn, &auth, spaces,
			0, DefaultTagsLimit,
			&SpaceFilter{
				Mode: SpaceFilterModeTopSubspaces,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("loading tags for spaces: %w", err)
		}

	}

	err = LoadSpaceContent(conn, &auth, spaces, true)
	if err != nil {
		return nil, fmt.Errorf("loading space details: %w", err)
	}

	return spaces, nil

}
