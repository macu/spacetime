package ajax

import (
	"database/sql"
	"fmt"
	"net/http"

	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/types"
)

func AjaxLoadNodeViewPage(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("id"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	node, err := treetime.LoadNodeHeaderByID(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, fmt.Errorf("loading node header by ID: %w", err))
		return nil, http.StatusInternalServerError
	}
	if node == nil {
		return nil, http.StatusNotFound
	}

	return node, http.StatusOK

}

func AjaxLoadNodeChildren(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("id"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	class := r.FormValue("class")
	if class != "" {
		if !treetime.IsValidNodeClass(class) {
			return nil, http.StatusBadRequest
		}
	}

	var offset uint
	if offsetString := r.FormValue("offset"); offsetString != "" {
		offset, err = types.AtoUint(offsetString)
		if err != nil {
			return nil, http.StatusBadRequest
		}
	}

	children, total, err := treetime.LoadNodeChildren(db, auth, id, offset,
		&treetime.NodeSearchParams{
			LimitToClass: class,
		})
	if err != nil {
		logging.LogError(r, auth, fmt.Errorf("loading child nodes: %w", err))
		return nil, http.StatusInternalServerError
	}

	return struct {
		Nodes []treetime.NodeHeader `json:"nodes"`
		Total uint                  `json:"total"`
	}{
		Nodes: children,
		Total: total,
	}, http.StatusOK

}
