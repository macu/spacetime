package treetime

import (
	"fmt"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodeTitle(db db.DBConn, auth *ajax.Auth, id uint) (string, error) {

	var title string

	err := db.QueryRow(`SELECT tree_node_content.text_content
		FROM tree_node_content
		WHERE tree_node_content.node_id = $1
		AND tree_node_content.content_type = $2
		LIMIT 1`,
		id,
		ContentTypeTitle,
	).Scan(&title)

	if err != nil {
		return "", fmt.Errorf("loading node title: %w", err)
	}

	return title, nil

}
