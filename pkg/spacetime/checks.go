package spacetime

import (
	"database/sql"
	"fmt"
	"spacetime/pkg/utils/types"
)

func CheckSpaceExists(conn *sql.DB, spaceID uint) (bool, error) {

	var exists bool

	var err = conn.QueryRow(`SELECT EXISTS (
		SELECT 1
		FROM space
		WHERE id = $1
	)`, spaceID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check space exists: %w", err)
	}

	return exists, nil

}

func CheckSubspaceLabelExists(conn *sql.DB, parentId *uint, label string) (bool, error) {

	labelUnitextId, err := GetUniqueTextId(conn, label)
	if err != nil {
		return false, fmt.Errorf("get label unique text id: %w", err)
	} else if labelUnitextId == nil {
		// label does not exist, so subspace with this label cannot exist
		return false, nil
	}

	var exists bool

	if parentId == nil {
		// check root spaces
		err = conn.QueryRow(`SELECT EXISTS (
			SELECT 1 FROM space
			INNER JOIN subspace
				ON space.id = subspace.space_id
			WHERE space.parent_id IS NULL
			AND subspace.label_text_id = $1
		)`, labelUnitextId).Scan(&exists)
	} else {
		// check subspaces
		err = conn.QueryRow(`SELECT EXISTS (
			SELECT 1 FROM space
			INNER JOIN subspace
				ON space.id = subspace.space_id
			WHERE space.parent_id = $1
			AND subspace.label_text_id = $2
		)`, *parentId, labelUnitextId).Scan(&exists)
	}

	if err != nil {
		return false, fmt.Errorf("check subspace label exists: %w", err)
	}

	return exists, nil

}

func CheckLinkSpaceExists(conn *sql.DB, parentID, linkSpaceID uint) (bool, error) {

	var exists bool

	var err = conn.QueryRow(`SELECT EXISTS (
		SELECT 1
		FROM space
		INNER JOIN link_space
			ON space.id = link_space.space_id
		WHERE space.parent_id = $1
		AND link_space.link_space_id = $2
	)`, parentID, linkSpaceID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check link space exists: %w", err)
	}

	return exists, nil

}

func ValidateLabel(label string) bool {
	if len(label) == 0 || len(label) > LabelMaxLength {
		return false
	}

	if types.HasUnprintableChar(label) || types.HasNewlines(label) {
		return false
	}

	// Check for newlines and invalid characters
	for _, c := range label {
		if c < 32 {
			return false
		}
	}

	return true
}

func CheckTitleExists(conn *sql.DB, parentID uint, title string) (bool, error) {

	uniqueTextId, err := GetUniqueTextId(conn, title)
	if err != nil {
		return false, fmt.Errorf("get title unique text id: %w", err)
	} else if uniqueTextId == nil {
		// title does not exist, so title space cannot exist
		return false, nil
	}

	var exists bool

	err = conn.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM space
		INNER JOIN title_space
			ON space.id = title_space.space_id
		WHERE space.parent_id = $1
		AND title_space.text_id = $2
	)`, parentID, *uniqueTextId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check title exists: %w", err)
	}

	return exists, nil

}

func CheckTagExists(conn *sql.DB, parentID uint, tag string) (bool, error) {

	tagUniqueTextId, err := GetUniqueTextId(conn, tag)
	if err != nil {
		return false, fmt.Errorf("get tag unique text id: %w", err)
	} else if tagUniqueTextId == nil {
		// tag does not exist, so tag space cannot exist
		return false, nil
	}

	var exists bool

	err = conn.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM space
		INNER JOIN tag_space
			ON space.id = tag_space.space_id
		WHERE space.parent_id = $1
		AND tag_space.text_id = $2
	)`, parentID, *tagUniqueTextId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check tag exists: %w", err)
	}

	return exists, nil

}

func ValidateTitle(title string) bool {
	if len(title) == 0 || len(title) > TitleMaxLength {
		return false
	}

	// Check for newlines and invalid characters
	for _, c := range title {
		if c < 32 {
			return false
		}
	}

	return true
}

func ValidateTag(tag string) bool {
	if len(tag) == 0 || len(tag) > TagMaxLength {
		return false
	}

	// Check for newlines and invalid characters
	for _, c := range tag {
		if c < 32 {
			return false
		}
	}

	return true
}

func ValidateText(text string) bool {
	if len(text) == 0 || len(text) > TextMaxLength {
		return false
	}

	// Check for invalid characters
	// Allow newlines and tabs
	for _, c := range text {
		if c < 32 {
			if c != '\n' && c != '\r' && c != '\t' {
				return false
			}
		}
	}

	return true
}
