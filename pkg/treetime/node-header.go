package treetime

import (
	"database/sql"
	"fmt"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodeHeaderByKey(db db.DBConn, auth *ajax.Auth, internalKey string) (*NodeHeader, error) {

	var header = &NodeHeader{}

	var err = db.QueryRow(`SELECT tree_node.id, tree_node.node_class, tree_node.is_deleted,
		tree_node_meta.internal_key
		FROM tree_node_meta
		INNER JOIN tree_node ON tree_node_meta.node_id = tree_node.id
		WHERE tree_node_meta.internal_key = $1`,
		internalKey,
	).Scan(&header.ID, &header.Class, &header.IsDeleted, &header.Key)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("node header not found by key: %s", internalKey)
		}
		return nil, fmt.Errorf("loading node header by key: %w", err)
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

func LoadNodeHeaderByID(db db.DBConn, auth *ajax.Auth, id uint) (*NodeHeader, error) {

	var header = &NodeHeader{}

	var err = db.QueryRow(`SELECT tree_node.id, tree_node.node_class, tree_node.is_deleted,
		tree_node_meta.internal_key
		FROM tree_node
		LEFT JOIN tree_node_meta ON tree_node.id = tree_node_meta.node_id
		WHERE tree_node.id = $1`,
		id,
	).Scan(&header.ID, &header.Class, &header.IsDeleted, &header.Key)

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
