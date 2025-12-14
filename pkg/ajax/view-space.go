package ajax

import (
	"database/sql"
	"encoding/json"
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

	filterJSON := r.FormValue("filter")
	filter := spacetime.SpaceFilter{}
	if filterJSON != "" {
		err = json.Unmarshal([]byte(filterJSON), &filter)
		if err != nil {
			return nil, http.StatusBadRequest
		}
	}

	space, err := spacetime.LoadSpace(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	err = spacetime.LoadSubspaceCount(db,
		[]*spacetime.Space{space}, &filter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if auth != nil {
		err = spacetime.LoadLastUserTitles(db, *auth,
			[]*spacetime.Space{space}, &filter)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	err = spacetime.LoadOriginalTitles(db,
		[]*spacetime.Space{space}, &filter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	err = spacetime.LoadTopTitles(db,
		[]*spacetime.Space{space}, &filter)
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
			&id, 0, spacetime.MaxSubspacesPageLimit, nil, nil)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}

		if auth != nil {
			err = spacetime.LoadLastUserTitles(db, *auth,
				content, &filter)
			if err != nil {
				logging.LogError(r, auth, err)
				return nil, http.StatusInternalServerError
			}
		}

		err = spacetime.LoadOriginalTitles(db, content, &filter)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}

		err = spacetime.LoadTopTitles(db, content, &filter)
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

	spaces, err := spacetime.LoadTopSubspaces(db, auth, parentId, offset, limit, nil, nil)
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
