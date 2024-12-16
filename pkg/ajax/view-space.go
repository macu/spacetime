package ajax

import (
	"database/sql"
	"net/http"

	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

func AjaxLoadSpace(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("id"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	includeSubspaces := types.AtoBool(r.FormValue("includeSubspaces"))

	space, err := spacetime.LoadSpace(db, auth, id, includeSubspaces, nil, nil)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusOK

}

// Load subspaces ordered by most checkins all time.
func AjaxLoadSubspacesByCheckinTotal(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	parentId, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	offset, err := types.AtoUint(r.FormValue("offset"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	limit, err := types.AtoUint(r.FormValue("limit"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	spaces, err := spacetime.LoadTopSubspaces(db, auth, parentId, offset, limit, nil, nil)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return spaces, http.StatusOK

}
