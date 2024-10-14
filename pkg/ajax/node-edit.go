package ajax

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/types"
)

func AjaxNodeLoadPath(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("id"))

	if err != nil {
		return nil, http.StatusBadRequest
	}

	path, err := treetime.LoadNodePath(db, &auth, id)

	if err != nil {
		logging.LogError(r, &auth, fmt.Errorf("loading node path: %w", err))
		return nil, http.StatusInternalServerError
	}

	return struct {
		Path []treetime.NodeHeader `json:"path"`
	}{
		Path: path,
	}, http.StatusOK

}

func AjaxNodeFindExisting(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var parentIDString = r.FormValue("parentId")
	var title = strings.TrimSpace(r.FormValue("title"))
	var description = strings.TrimSpace(r.FormValue("description"))
	var class = strings.TrimSpace(r.FormValue("class"))

	if title == "" || !treetime.IsValidNodeClass(class) {
		return nil, http.StatusBadRequest
	}

	var parentID *uint
	var err error
	if parentIDString != "" {
		var pID uint
		pID, err = types.AtoUint(parentIDString)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		parentID = &pID
	}

	var nodes []treetime.NodeHeader
	nodes, err = treetime.FindExistingNodes(db, &auth, parentID, class, title, description)

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

	var title = strings.TrimSpace(r.FormValue("title"))
	var class = strings.TrimSpace(r.FormValue("class"))
	var description = strings.TrimSpace(r.FormValue("description"))

	if title == "" || !treetime.IsValidNodeClass(class) {
		return nil, http.StatusBadRequest
	}

	if !treetime.CheckContentLengthLimit(class, treetime.ContentTypeTitle, title) {
		return nil, http.StatusBadRequest
	}

	if !treetime.CheckContentLengthLimit(class, treetime.ContentTypeBody, description) {
		return nil, http.StatusBadRequest
	}

	var err error

	var parentIDString = r.FormValue("parentId")
	var parentID *uint

	if parentIDString != "" {
		var pID uint
		pID, err = types.AtoUint(parentIDString)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		parentID = &pID
	}

	var node *treetime.NodeHeader
	err = db.InTransaction(r, conn, func(tx *sql.Tx) error {
		node, err = treetime.CreateNode(tx, auth, parentID, class, title, description)
		return err
	})

	if err != nil {
		logging.LogError(r, &auth, fmt.Errorf("creating node: %w", err))
		return nil, http.StatusInternalServerError
	}

	return node, http.StatusOK

}

func AjaxCreateNodeContent(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	return nil, http.StatusNotImplemented

}
