package treetime

import (
	"fmt"
	"regexp"
	"strings"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

const findExistingLimit = 20

func FindExistingNodes(db db.DBConn, auth *ajax.Auth, parentID *uint, class, query string) ([]NodeHeader, error) {

	var nodes = []NodeHeader{}

	query = strings.TrimSpace(query)
	if query == "" {
		return nodes, nil
	}

	// Prepapre for to_tsquery
	query = whitespaceRegex.ReplaceAllString(query, " | ")

	var args = []interface{}{class, query, findExistingLimit}

	var orderBy string
	if parentID != nil {
		args = append(args, *parentID)
		orderBy = `(tree_node.parent_id = $3) DESC, ordered_ranks.max_rank DESC`
	} else {
		orderBy = "ordered_ranks.max_rank DESC"
	}

	rows, err := db.Query(`WITH content_ranks AS (
			SELECT tree_node.id,
				ts_rank_cd(tree_node_content.text_search, to_tsquery('pg_catalog.simple', $2)) AS rank
			FROM tree_node
			INNER JOIN tree_node_content ON tree_node.id = tree_node_content.node_id
			WHERE tree_node.node_class = $1
			AND ts_rank(tree_node_content.text_search, to_tsquery('pg_catalog.simple', $2)) > 0
		), ordered_ranks AS (
		SELECT id,
			rank,
			MAX(rank) AS max_rank
		FROM content_ranks
		)
		SELECT tree_node.id,
			tree_node.node_class,
			tree_node_meta.internal_key
		FROM tree_node
		LEFT JOIN tree_node_meta ON tree_node.id = tree_node_meta.node_id
		INNER JOIN ordered_ranks ON tree_node.id = ordered_ranks.id
		WHERE tree_node.is_deleted = FALSE
		ORDER BY `+orderBy+`
		LIMIT $3`,
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

	err = LoadContentForNodes(db, auth, nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil

}
