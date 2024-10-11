package ajax

import (
	"database/sql"
	"fmt"
	"net/http"

	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
)

func AjaxDashboard(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var treetimeNode, err = treetime.LoadNodeHeaderByKey(db, auth, treetime.NodeKeyTreeTime)

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
