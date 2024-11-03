package ajax

import (
	"database/sql"
	"fmt"
	"net/http"

	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
)

func AjaxLoadLangauges(db *sql.DB, auth *ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	nodes, err := treetime.LoadLangs(db)
	if err != nil {
		logging.LogError(r, auth, fmt.Errorf("loading lang nodes: %w", err))
		return nil, http.StatusInternalServerError
	}

	return nodes, http.StatusOK

}
