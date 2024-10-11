package treetime

import (
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodeParentPath(db db.DBConn, auth *ajax.Auth, id uint) ([]NodeHeader, error) {
	var path = make([]NodeHeader, 0)

	rows, err := db.Query(`WITH RECURSIVE parent_nodes AS (
			SELECT
				tree_node.id,
				tree_node.node_class,
				tree_node_internal_key.internal_key,
				tree_node.parent_id
			FROM
				tree_node
			LEFT JOIN
				tree_node_internal_key
			ON
				tree_node.id = tree_node_internal_key.node_id
			WHERE
				tree_node.id = $1 -- Starting node ID

			UNION ALL

			SELECT
				tn.id,
				tn.node_class,
				tnik.internal_key,
				tn.parent_id
			FROM
				tree_node tn
			LEFT JOIN
				tree_node_internal_key tnik
			ON
				tn.id = tnik.node_id
			INNER JOIN
				parent_nodes pn
			ON
				tn.id = pn.parent_id
		)
		SELECT
			id,
			node_class,
			internal_key
		FROM
			parent_nodes
		WHERE
			id != $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var parentHeader = &NodeHeader{}
		rows.Scan(&parentHeader.ID, &parentHeader.Class, &parentHeader.Key)
		parentHeader.Title, err = LoadNodeTitle(db, auth, parentHeader.ID)
		if err != nil {
			return nil, err
		}
		// Prepend to path
		path = append([]NodeHeader{*parentHeader}, path...)
	}

	return path, nil

}
