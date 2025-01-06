package ajax

import (
	"database/sql"
	"net/http"
	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

// Load subspaces ordered by most checkins all time.
func AjaxLoadTopSubspaces(db *sql.DB, auth *ajax.Auth,
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

	includeTags := types.AtoBool(r.FormValue("includeTags"))

	spaces, err := spacetime.LoadTopSubspaces(db, auth, parentId, offset, limit, nil, nil)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if auth != nil {
		err = spacetime.LoadLastUserTitles(db, *auth, spaces)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	err = spacetime.LoadOriginalTitles(db, spaces)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	err = spacetime.LoadTopTitles(db, spaces)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if includeTags {
		err = spacetime.LoadTopTags(db, spaces, 0, spacetime.DefaultTagsLimit)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return spaces, http.StatusOK

}

// Load titles ordered by most subspaces.
func AjaxLoadTopTitles(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	parentId, err := types.AtoUint(r.FormValue("parentId"))
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

	titles, err := spacetime.LoadMoreTitles(db, parentId, offset, limit)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return titles, http.StatusOK

}

// Load tags ordered by most subspaces.
func AjaxLoadTopTags(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	parentId, err := types.AtoUint(r.FormValue("parentId"))
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

	tags, err := spacetime.LoadMoreTags(db, parentId, offset, limit)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return tags, http.StatusOK

}
