package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CreateCheckin(conn *sql.DB, auth ajax.Auth, parentID uint) (*Space, error) {

	var space = Space{}

	if err := CreateSpace(conn, auth, &space, &parentID, SpaceTypeCheckin); err != nil {
		return nil, fmt.Errorf("create checkin: %w", err)
	}

	return &space, nil

}
