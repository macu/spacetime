package treetime

import (
	"database/sql"
	"fmt"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodeHeaderByKey(db db.DBConn, auth *ajax.Auth, internalKey string) (*NodeHeader, error) {

	var header = &NodeHeader{}

	var err = db.QueryRow(`SELECT tree_node.id, tree_node.node_class,
		tree_node_internal_key.internal_key
		FROM tree_node_internal_key
		INNER JOIN tree_node ON tree_node_internal_key.node_id = tree_node.id
		WHERE tree_node_internal_key.internal_key = $1`,
		internalKey,
	).Scan(&header.ID, &header.Class, &header.Key)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("loading node header by key: %w", err)
	}

	header.Title, err = LoadNodeTitle(db, auth, header.ID)

	if err != nil {
		return nil, err
	}

	header.Tags, err = LoadNodeTags(db, auth, header.ID)

	if err != nil {
		return nil, err
	}

	return header, nil

}

func LoadNodeHeaderByID(db db.DBConn, auth *ajax.Auth, id uint) (*NodeHeader, error) {

	var header = &NodeHeader{}

	var err = db.QueryRow(`SELECT tree_node.id, tree_node.node_class,
		tree_node_internal_key.internal_key
		FROM tree_node
		LEFT JOIN tree_node_internal_key ON tree_node.id = tree_node_internal_key.node_id
		WHERE tree_node.id = $1`,
		id,
	).Scan(&header.ID, &header.Class, &header.Key)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("loading node header by id: %w", err)
	}

	header.Title, err = LoadNodeTitle(db, auth, header.ID)

	if err != nil {
		return nil, err
	}

	header.Tags, err = LoadNodeTags(db, auth, header.ID)

	if err != nil {
		return nil, err
	}

	return header, nil

}
