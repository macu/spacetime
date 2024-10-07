package ajax

import (
	"database/sql"
	"fmt"
	"net/http"

	"treetime/pkg/treetime"
	"treetime/pkg/utils/logging"
)

func AjaxDashboard(db *sql.DB, userID *uint, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var treetimeNode, err = treetime.LoadNodeHeaderByKey(db, userID, "treetime")

	if err != nil {
		logging.LogError(r, userID, fmt.Errorf("loading category header by key: %w", err))
		return nil, http.StatusInternalServerError
	}

	return struct {
		Node treetime.NodeHeader `json:"treetimeNode"`
	}{
		Node: *treetimeNode,
	}, http.StatusOK

}
