package treetime

import (
	"fmt"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

const childNodesPageLimit = 20

func LoadNodeChildren(conn db.DBConn, auth *ajax.Auth, id uint, offset uint, query *NodeSearchParams) ([]NodeHeader, error) {

	var children = make([]NodeHeader, 0)

	var args = []interface{}{id, offset, childNodesPageLimit}

	var limitToClassPart string
	if query != nil && query.LimitToClass != "" {
		limitToClassPart = "AND tn.node_class = " + db.Arg(&args, query.LimitToClass)
	}

	var excludeClassPart string
	if query != nil && query.ExcludeClass != "" {
		excludeClassPart = "AND tn.node_class != " + db.Arg(&args, query.ExcludeClass)
	}

	rows, err := conn.Query(`SELECT
			tn.id,
			tn.node_class,
			tnik.internal_key
		FROM
			tree_node tn
		LEFT JOIN
			tree_node_internal_key tnik
		ON
			tn.id = tnik.node_id
		WHERE
			tn.parent_id = $1
			`+limitToClassPart+`
			`+excludeClassPart+`
		ORDER BY tn.id
		OFFSET $2
		LIMIT $3`,
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf("loading node children: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var childHeader = &NodeHeader{}
		err = rows.Scan(&childHeader.ID, &childHeader.Class, &childHeader.Key)
		if err != nil {
			return nil, fmt.Errorf("scanning node header: %w", err)
		}
		children = append(children, *childHeader)
	}

	err = LoadNodeTitles(conn, auth, children)
	if err != nil {
		return nil, err
	}

	return children, nil

}
