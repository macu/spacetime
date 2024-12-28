package spacetime

import (
	"database/sql"

	"spacetime/pkg/utils/ajax"
)

func CreateStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, parentID *uint) (*Space, error) {

	// Create an open stream of consciousness space
	// (will hold a series of naked texts created by author)

	return nil, nil

}

func CloseStreamOfConsciousness(conn *sql.DB, auth ajax.Auth, id uint) error {

	// Mark stream of consciousness as "closed" by author

	return nil

}
