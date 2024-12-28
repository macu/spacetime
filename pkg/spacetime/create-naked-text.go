package spacetime

import (
	"database/sql"

	"spacetime/pkg/utils/ajax"
)

func CreateNakedText(conn *sql.DB, auth ajax.Auth, parentID uint, finalText, replayData string) (*Space, error) {

	// Create naked text space with given replay data

	return nil, nil

}
