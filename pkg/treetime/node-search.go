package treetime

import (
	"fmt"
	"regexp"
	"strings"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

func FindExistingNodes(db db.DBConn, auth *ajax.Auth, parentID *uint, class, query string) ([]NodeHeader, error) {

	var nodes = []NodeHeader{}

	query = strings.TrimSpace(query)
	if query == "" {
		return nodes, nil
	}

	// Prepapre for to_tsquery
	query = whitespaceRegex.ReplaceAllString(query, " | ")

	var args = []interface{}{class, query}

	var orderBy string

	if parentID != nil {
		args = append(args, *parentID)
		orderBy = `ORDER BY (tree_node.parent_id = $3) DESC,
		ts_rank_cd(tree_node_content.text_search, to_tsquery('pg_catalog.simple', $2)) DESC`
	} else {
		orderBy = "ORDER BY ts_rank_cd(tree_node_content.text_search, to_tsquery('pg_catalog.simple', $2)) DESC"
	}

	rows, err := db.Query(`
		SELECT tree_node.id,
			tree_node.node_class,
			tree_node_internal_key.internal_key
		FROM tree_node
		LEFT JOIN tree_node_internal_key ON tree_node.id = tree_node_internal_key.node_id
		INNER JOIN tree_node_content ON tree_node.id = tree_node_content.node_id
		WHERE tree_node.node_class = $1
		AND ts_rank(tree_node_content.text_search, to_tsquery('pg_catalog.simple', $2)) > 0
		`+orderBy+`
		LIMIT 20`,
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf("querying for existing nodes: %w", err)
	}

	defer rows.Close()

	var alreadySeenIDs = make(map[uint]bool)

	for rows.Next() {
		var node NodeHeader
		err = rows.Scan(&node.ID, &node.Class, &node.Key)
		if err != nil {
			return nil, fmt.Errorf("scanning node: %w", err)
		}
		if alreadySeenIDs[node.ID] {
			continue
		}
		alreadySeenIDs[node.ID] = true
		nodes = append(nodes, node)
	}

	err = LoadNodeTitles(db, auth, nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil

}
