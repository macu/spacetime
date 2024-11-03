package ajax

import (
	"database/sql"
	"fmt"
	"net/http"

	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
)

const treeTimeCategoryID = 1 // first category created in init script

func AjaxDashboard(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var treetimeNode, err = treetime.LoadNodeHeaderByID(db, auth, treeTimeCategoryID)

	if err != nil {
		logging.LogError(r, auth, fmt.Errorf("loading category header by key: %w", err))
		return nil, http.StatusInternalServerError
	}

	return struct {
		Node treetime.NodeHeader `json:"treetimeNode"`
	}{
		Node: *treetimeNode,
	}, http.StatusOK

}
