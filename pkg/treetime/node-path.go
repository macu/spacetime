package treetime

import (
	"fmt"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func LoadNodePath(db db.DBConn, auth *ajax.Auth, id uint, orderRootFirst bool) ([]NodeHeader, error) {
	var path = make([]NodeHeader, 0)

	rows, err := db.Query(`WITH RECURSIVE node_path AS (
			SELECT
				tree_node.id,
				tree_node.node_class,
				tree_node.is_deleted,
				tree_node.parent_id,
				tree_node.owner_type,
				tree_node.created_by
			FROM
				tree_node
			WHERE
				tree_node.id = $1 -- Starting node ID

			UNION ALL

			SELECT
				tn.id,
				tn.node_class,
				tn.is_deleted,
				tn.parent_id,
				tn.owner_type,
				tn.created_by
			FROM
				tree_node tn
			INNER JOIN
				node_path pn
			ON
				tn.id = pn.parent_id
		)
		SELECT
			id,
			node_class,
			is_deleted,
			owner_type,
			created_by
		FROM
			node_path`,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("loading node path: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var node = NodeHeader{
			Creator: &NodeCreator{},
		}
		err = rows.Scan(&node.ID, &node.Class,
			&node.IsDeleted, &node.OwnerType, &node.Creator.ID)
		if err != nil {
			return nil, fmt.Errorf("scanning node path: %w", err)
		}
		path = append(path, node)
	}

	if orderRootFirst {
		// Reverse path
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
	}

	err = LoadContentForNodes(db, auth, path)
	if err != nil {
		return nil, fmt.Errorf("loading content for nodes: %w", err)
	}

	return path, nil

}

func LoadNodeParentPath(db db.DBConn, auth *ajax.Auth, id uint, orderRootFirst bool) ([]NodeHeader, error) {

	var path = make([]NodeHeader, 0)

	rows, err := db.Query(`WITH RECURSIVE parent_nodes AS (
			SELECT
				tree_node.id,
				tree_node.node_class,
				tree_node.is_deleted,
				tree_node.parent_id,
				tree_node.owner_type,
				tree_node.created_by
			FROM
				tree_node
			WHERE
				tree_node.id = $1 -- Starting node ID

			UNION ALL

			SELECT
				tn.id,
				tn.node_class,
				tn.is_deleted,
				tn.parent_id,
				tn.owner_type,
				tn.created_by
			FROM
				tree_node tn
			INNER JOIN
				parent_nodes pn
			ON
				tn.id = pn.parent_id
		)
		SELECT
			id,
			node_class,
			is_deleted,
			owner_type,
			created_by
		FROM
			parent_nodes
		WHERE
			id != $1`,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("loading parent nodes: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var node = NodeHeader{
			Creator: &NodeCreator{},
		}
		err = rows.Scan(&node.ID, &node.Class,
			&node.IsDeleted, &node.OwnerType, &node.Creator.ID)
		if err != nil {
			return nil, fmt.Errorf("scanning parent nodes: %w", err)
		}
		path = append(path, node)
	}

	if orderRootFirst {
		// Reverse path
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
	}

	err = LoadContentForNodes(db, auth, path)
	if err != nil {
		return nil, fmt.Errorf("loading content for parent nodes: %w", err)
	}

	return path, nil

}
