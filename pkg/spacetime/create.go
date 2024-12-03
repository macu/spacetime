package spacetime

import (
	"database/sql"

	"spacetime/pkg/utils/ajax"
)

func CreateEmptySpace(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	return nil, nil

}

func CreateCheckin(conn *sql.DB, auth ajax.Auth, parentID uint, spaceID *uint) (*Space, error) {

	// Check if checkin already exists
	// Create checkin space if not exists
	// Check-in under existing checkin space if exists

	return nil, nil

}

func CreateTitleCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, title string) (*Space, error) {

	// Load unique_text ID
	// Check for existing title space under parent
	// Create title space if not exists
	// Check-in on title space

	return nil, nil

}

func CreateTagCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing tag space under parent
	// Create tag space if not exists
	// Check-in on tag space

	return nil, nil

}

func CreateTextCheckin(conn *sql.DB, auth ajax.Auth, parentID *uint, text string) (*Space, error) {

	// Load unique_text ID
	// Check for existing text space under parent
	// Create text space if not exists
	// Check-in on text space

	return nil, nil

}

func CreateNakedText(conn *sql.DB, auth ajax.Auth, parentID *uint, finalText, replayData string) (*Space, error) {

	// Create naked text space with given replay data

	return nil, nil

}

func CreateStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	// Create an open stream of consciousness space

	return nil, nil

}

func CloseStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, id *uint) error {

	// Mark stream of consciousness as "closed" by user who created it

	return nil

}

func CreateJSONAttribute(conn *sql.DB, auth ajax.Auth, parentID *uint, url, path string) (*Space, error) {

	// Check if space exists
	// Create if not exists

	return nil, nil

}
