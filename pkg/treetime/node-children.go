package treetime

import (
	"fmt"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

const childNodesPageLimit = 20

func LoadNodeChildren(conn db.DBConn, auth *ajax.Auth, id uint, offset uint, query *NodeSearchParams) ([]NodeHeader, uint, error) {

	var children = make([]NodeHeader, 0)

	var args = []interface{}{id, offset, childNodesPageLimit}

	var limitToClassPart string
	if query != nil && query.LimitToClass != "" {
		limitToClassPart = "AND tree_node.node_class = " + db.Arg(&args, query.LimitToClass)
	}

	var excludeClassPart string
	if query != nil && query.ExcludeClass != "" {
		excludeClassPart = "AND tree_node.node_class != " + db.Arg(&args, query.ExcludeClass)
	}

	rows, err := conn.Query(`SELECT tree_node.id, tree_node.node_class, tree_node.created_at,
		tree_node.owner_type, tree_node.created_by, user_account.display_name
		FROM tree_node
		LEFT JOIN user_account
		ON user_account.id = tree_node.created_by
		WHERE
			tree_node.parent_id = $1
			AND tree_node.is_deleted = FALSE
			`+limitToClassPart+`
			`+excludeClassPart+`
		ORDER BY tree_node.id
		OFFSET $2
		LIMIT $3`,
		args...,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("loading node children: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var childHeader = &NodeHeader{
			Creator: &NodeCreator{},
		}
		err = rows.Scan(&childHeader.ID, &childHeader.Class, &childHeader.CreatedAt,
			&childHeader.OwnerType, &childHeader.Creator.ID, &childHeader.Creator.DisplayName)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning node header: %w", err)
		}
		children = append(children, *childHeader)
	}

	var total uint
	err = conn.QueryRow(`SELECT COUNT(*)
		FROM tree_node
		WHERE
			parent_id = $1
			AND is_deleted = FALSE
			`+limitToClassPart+`
			`+excludeClassPart,
		id,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting node children: %w", err)
	}

	// TODO Load path for any non-direct children

	err = LoadContentForNodes(conn, auth, children)
	if err != nil {
		return nil, 0, err
	}

	return children, total, nil

}
