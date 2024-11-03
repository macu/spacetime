package treetime

import (
	"database/sql"
	"fmt"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodeClass(db db.DBConn, id uint) (string, error) {

	var class string

	err := db.QueryRow(`SELECT tree_node.node_class
		FROM tree_node
		WHERE tree_node.id = $1`,
		id,
	).Scan(&class)

	if err != nil {
		return "", fmt.Errorf("loading node class: %w", err)
	}

	return class, nil

}

func LoadNodeHeaderByID(db db.DBConn, auth *ajax.Auth, id uint) (*NodeHeader, error) {

	var header = &NodeHeader{}

	var err = db.QueryRow(`SELECT tree_node.id, tree_node.node_class,
		tree_node.is_deleted, tree_node.owner_type, tree_node.created_by
		FROM tree_node
		WHERE tree_node.id = $1`,
		id,
	).Scan(&header.ID, &header.Class,
		&header.IsDeleted, &header.OwnerType, &header.CreatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("node header not found by id: %d", id)
		}
		return nil, fmt.Errorf("loading node header by id: %w", err)
	}

	err = LoadContentForNode(db, auth, header)
	if err != nil {
		return nil, err
	}

	err = LoadNodeTags(db, auth, header)
	if err != nil {
		return nil, err
	}

	return header, nil

}
