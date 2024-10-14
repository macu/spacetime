package treetime

import (
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodePath(db db.DBConn, auth *ajax.Auth, id uint, orderRootFirst bool) ([]NodeHeader, error) {
	var path = make([]NodeHeader, 0)

	rows, err := db.Query(`WITH RECURSIVE node_path AS (
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
				node_path pn
			ON
				tn.id = pn.parent_id
		)
		SELECT
			id,
			node_class,
			internal_key
		FROM
			node_path`,
		id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var node = NodeHeader{}
		err = rows.Scan(&node.ID, &node.Class, &node.Key)
		if err != nil {
			return nil, err
		}
		path = append(path, node)
	}

	if orderRootFirst {
		// Reverse path
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
	}

	err = LoadNodeTitles(db, auth, path)
	if err != nil {
		return nil, err
	}

	return path, nil

}

func LoadNodeParentPath(db db.DBConn, auth *ajax.Auth, id uint, orderRootFirst bool) ([]NodeHeader, error) {
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
		var node = NodeHeader{}
		err = rows.Scan(&node.ID, &node.Class, &node.Key)
		if err != nil {
			return nil, err
		}
		path = append(path, node)
	}

	if orderRootFirst {
		// Reverse path
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
	}

	err = LoadNodeTitles(db, auth, path)
	if err != nil {
		return nil, err
	}

	return path, nil

}
