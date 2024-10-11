package treetime

import (
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

const childNodesPageLimit = 20

func LoadNodeChildren(conn db.DBConn, auth *ajax.Auth, id uint, offset uint, query *NodeSearchParams) ([]NodeHeader, error) {

	var children = make([]NodeHeader, 0)

	var args = []interface{}{id, offset, childNodesPageLimit}

	var limitToClassPart string
	if query != nil && query.LimitToClass != "" {
		limitToClassPart = "AND tn.node_class = " +
			db.ArgPlaceholder(string(query.LimitToClass), &args)
	}

	var excludeClassPart string
	if query != nil && query.ExcludeClass != "" {
		excludeClassPart = "AND tn.node_class != " +
			db.ArgPlaceholder(string(query.ExcludeClass), &args)
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
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var childHeader = &NodeHeader{}
		rows.Scan(&childHeader.ID, &childHeader.Class, &childHeader.Key)
		childHeader.Title, err = LoadNodeTitle(conn, auth, childHeader.ID)
		if err != nil {
			return nil, err
		}
		children = append(children, *childHeader)
	}

	return children, nil

}
