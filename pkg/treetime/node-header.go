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

	var header = &NodeHeader{
		Creator: &NodeCreator{},
	}

	var err = db.QueryRow(`SELECT tree_node.id, tree_node.node_class,
		tree_node.is_deleted, tree_node.owner_type,
		tree_node.created_at,
		tree_node.created_by, user_account.display_name
		FROM tree_node
		LEFT JOIN user_account
		ON user_account.id = tree_node.created_by
		WHERE tree_node.id = $1`,
		id,
	).Scan(&header.ID, &header.Class,
		&header.IsDeleted, &header.OwnerType,
		&header.CreatedAt,
		&header.Creator.ID, &header.Creator.DisplayName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("node header not found by id: %d", id)
		}
		return nil, fmt.Errorf("loading node header by id: %w", err)
	}

	err = LoadContentForNode(db, auth, header)
	if err != nil {
		return nil, fmt.Errorf("loading content for node: %w", err)
	}

	path, err := LoadNodeParentPath(db, auth, header.ID, true)
	if err != nil {
		return nil, fmt.Errorf("loading parent nodes: %w", err)
	}
	header.Path = &path

	err = LoadNodeTags(db, auth, header)
	if err != nil {
		return nil, fmt.Errorf("loading tags for node: %w", err)
	}

	return header, nil

}
