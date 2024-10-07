package treetime

import (
	"fmt"
	"treetime/pkg/utils/db"
)

type NodeContentType string

const (
	NodeContentTypeTitle NodeContentType = "title"
	NodeContentTypeBody  NodeContentType = "body"
)

func LoadNodeTitle(db db.DBConn, userId *uint, nodeId uint) (string, error) {
	var title string

	err := db.QueryRow(`SELECT tree_node_content.text_content
		FROM tree_node_content
		WHERE tree_node_content.node_id = $1
		AND tree_node_content.content_type = $2
		LIMIT 1`,
		nodeId, NodeContentTypeTitle,
	).Scan(&title)

	if err != nil {
		return "", fmt.Errorf("loading node title: %w", err)
	}

	return title, nil
}
