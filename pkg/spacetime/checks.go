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
