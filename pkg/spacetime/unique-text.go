package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/db"
)

func GetUniqueTextId(conn db.DBConn, text string) (*uint, error) {
	var uniqueTextId *uint
	err := conn.QueryRow(`SELECT id FROM unique_text WHERE text = $1`, text).Scan(&uniqueTextId)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("query unique_text: %w", err)
	}
	return uniqueTextId, nil
}

func CreateUniqueText(conn db.DBConn, text string) (*uint, error) {
	var uniqueTextId *uint
	err := conn.QueryRow(`INSERT INTO unique_text (text)
		VALUES ($1)
		RETURNING id`, text).Scan(&uniqueTextId)
	if err != nil {
		return nil, fmt.Errorf("insert unique_text: %w", err)
	}
	return uniqueTextId, nil
}
