package ajax

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/types"
)

func AjaxLoadCreateNode(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	parentID, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	class := r.FormValue("class")
	if !treetime.IsValidNodeClass(class) {
		return nil, http.StatusBadRequest
	}

	path := []treetime.NodeHeader{}

	if parentID != nil {
		path, err = treetime.LoadNodePath(db, &auth, *parentID)
		if err != nil {
			logging.LogError(r, &auth, fmt.Errorf("loading node path: %w", err))
			return nil, http.StatusInternalServerError
		}
	}

	createAllowed, err := treetime.CheckCreateNodeAllowed(db, parentID, class)
	if err != nil {
		logging.LogError(r, &auth, fmt.Errorf("verifying create node allowed: %w", err))
		return nil, http.StatusInternalServerError
	}

	return struct {
		Path          []treetime.NodeHeader `json:"path"`
		CreateAllowed bool                  `json:"createAllowed"`
	}{
		Path:          path,
		CreateAllowed: createAllowed,
	}, http.StatusOK

}

func AjaxFindExistingNode(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	parentID, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	query := strings.TrimSpace(r.FormValue("query"))
	class := strings.TrimSpace(r.FormValue("class"))
	if query == "" || !treetime.IsValidNodeClass(class) {
		return nil, http.StatusBadRequest
	}

	nodes, err := treetime.FindExistingNodes(db, &auth, parentID, class, query)
	if err != nil {
		logging.LogError(r, &auth, fmt.Errorf("finding existing nodes: %w", err))
		return nil, http.StatusInternalServerError
	}

	return struct {
		Nodes []treetime.NodeHeader `json:"nodes"`
	}{
		Nodes: nodes,
	}, http.StatusOK

}

func AjaxCreateNode(conn *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	parentID, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	class := strings.TrimSpace(r.FormValue("class"))
	if !treetime.IsValidNodeClass(class) {
		return nil, http.StatusBadRequest
	}

	title := treetime.FormatTitle(r.FormValue("title"))
	if !treetime.CheckContentLength(class, treetime.ContentTypeTitle, title) {
		return nil, http.StatusBadRequest
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if !treetime.CheckContentLength(class, treetime.ContentTypeBody, body) {
		return nil, http.StatusBadRequest
	}

	node, err := treetime.CreateNode(conn, auth, parentID, class, title, body)
	if err != nil {
		logging.LogError(r, &auth, fmt.Errorf("creating node: %w", err))
		return nil, http.StatusInternalServerError
	}

	return node, http.StatusOK

}
