package spacetime

import (
	"database/sql"
	"fmt"
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

func GetSpaceMeta(conn *sql.DB, spaceID uint) (*uint, string, error) {

	var parentID *uint
	var spaceType string

	var err = conn.QueryRow(`SELECT parent_id, space_type
		FROM space
		WHERE id = $1`, spaceID).Scan(&parentID, &spaceType)

	if err != nil {
		return nil, "", fmt.Errorf("get space parent ID: %w", err)
	}

	return parentID, spaceType, nil

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
			if c != '\n' && c != '\t' {
				return false
			}
		}
	}

	return true
}
