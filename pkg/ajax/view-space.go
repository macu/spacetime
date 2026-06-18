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

	id, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	includeTags := types.AtoBool(r.FormValue("includeTags"))
	includeSubspaces := types.AtoBool(r.FormValue("includeSubspaces"))
	includeParentPath := types.AtoBool(r.FormValue("includeParentPath"))

	filter, err := spacetime.ParseSpaceFilter(r.FormValue("filter"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.LoadSpace(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	err = spacetime.LoadCheckinCount(db,
		[]*spacetime.Space{space})
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if includeTags {
		err = spacetime.LoadTopTags(db,
			[]*spacetime.Space{space}, 0, spacetime.DefaultTagsLimit)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	if includeSubspaces {
		content, err := spacetime.LoadTopSubspaces(db, auth,
			&id, []string{}, 0, spacetime.MaxSubspacesPageLimit, filter)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}

		if includeTags {
			err = spacetime.LoadTopTags(db,
				content, 0, spacetime.DefaultTagsLimit)
			if err != nil {
				logging.LogError(r, auth, err)
				return nil, http.StatusInternalServerError
			}
		}

		space.TopSubspaces = &content
	}

	if includeParentPath {
		if space.ParentID == nil {
			space.ParentPath = &[]*spacetime.Space{}
		} else {
			path, err := spacetime.LoadParentPath(db, auth, *space.ParentID)
			if err != nil {
				logging.LogError(r, auth, err)
				return nil, http.StatusInternalServerError
			}

			space.ParentPath = &path
		}
	}

	return space, http.StatusOK

}

func AjaxLoadTextSpaceRecording(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.LoadSpace(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if space.SpaceType != spacetime.SpaceTypeText {
		return nil, http.StatusBadRequest
	}

	err = spacetime.LoadTextRecording(db, space)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return space.ReplayData, http.StatusOK
}

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

	includeTypes, err := types.AtoStringArray(r.FormValue("includeTypes"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	filter, err := spacetime.ParseSpaceFilter(r.FormValue("filter"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	spaces, err := spacetime.LoadTopSubspaces(db, auth, parentId, includeTypes, offset, limit, filter)
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
