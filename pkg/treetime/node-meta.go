package treetime

import (
	"fmt"
	"treetime/pkg/utils/db"
)

func LoadNodeMeta(db db.DBConn, id uint) (string, *string, error) {

	var class string
	var key *string

	err := db.QueryRow(`SELECT tree_node.node_class, tree_node_meta.internal_key
		FROM tree_node
		LEFt JOIN tree_node_meta ON tree_node.id = tree_node_meta.node_id
		WHERE tree_node.id = $1`,
		id,
	).Scan(&class, &key)

	if err != nil {
		return "", nil, fmt.Errorf("loading node class: %w", err)
	}

	return class, key, nil

}
