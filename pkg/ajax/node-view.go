package ajax

import (
	"database/sql"
	"fmt"
	"net/http"

	"treetime/pkg/treetime"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/types"
)

const childNodesPageLimit = 20

func AjaxLoadNodeViewPage(db *sql.DB, userID *uint, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("id"))

	if err != nil {
		return nil, http.StatusBadRequest
	}

	header, err := treetime.LoadNodeHeaderByID(db, userID, id)

	if err != nil {
		logging.LogError(r, userID, fmt.Errorf("loading node header by ID: %w", err))
		return nil, http.StatusInternalServerError
	}

	if header == nil {
		return nil, http.StatusNotFound
	}

	var parentPath []treetime.NodeHeader

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
			parent_nodes`,
		id,
	)

	if err != nil {
		logging.LogError(r, userID, fmt.Errorf("loading parent nodes: %w", err))
		return nil, http.StatusInternalServerError
	}

	defer rows.Close()

	for rows.Next() {
		var parentHeader = &treetime.NodeHeader{}
		rows.Scan(&parentHeader.ID, &parentHeader.Class, &parentHeader.Key)
		parentHeader.Title, err = treetime.LoadNodeTitle(db, userID, parentHeader.ID)
		if err != nil {
			logging.LogError(r, userID, fmt.Errorf("loading parent node title: %w", err))
			return nil, http.StatusInternalServerError
		}
		// Prepend to path
		parentPath = append([]treetime.NodeHeader{*parentHeader}, parentPath...)
	}

	// Drop last node
	parentPath = parentPath[:len(parentPath)-1]

	return struct {
		Header treetime.NodeHeader   `json:"header"`
		Path   []treetime.NodeHeader `json:"path"`
	}{
		Header: *header,
		Path:   parentPath,
	}, http.StatusOK

}

func AjaxLoadNodeChildren(db *sql.DB, userID *uint, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("id"))

	if err != nil {
		return nil, http.StatusBadRequest
	}

	var offset uint

	if offsetString := r.FormValue("offset"); offsetString != "" {
		offset, err = types.AtoUint(offsetString)
		if err != nil {
			return nil, http.StatusBadRequest
		}
	}

	rows, err := db.Query(`SELECT
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
		ORDER BY tn.id
		OFFSET $2
		LIMIT $3`,
		id,
		offset,
		childNodesPageLimit,
	)

	if err != nil {
		logging.LogError(r, userID, fmt.Errorf("loading child nodes: %w", err))
		return nil, http.StatusInternalServerError
	}

	var children = make([]treetime.NodeHeader, 0)

	defer rows.Close()

	for rows.Next() {
		var childHeader = &treetime.NodeHeader{}
		rows.Scan(&childHeader.ID, &childHeader.Class, &childHeader.Key)
		childHeader.Title, err = treetime.LoadNodeTitle(db, userID, childHeader.ID)
		if err != nil {
			logging.LogError(r, userID, fmt.Errorf("loading child node title: %w", err))
			return nil, http.StatusInternalServerError
		}
		children = append(children, *childHeader)
	}

	return children, http.StatusOK

}
