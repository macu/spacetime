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

	space, err := spacetime.LoadSpace(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	err = spacetime.LoadSubspaceCount(db,
		[]*spacetime.Space{space})
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if auth != nil {
		err = spacetime.LoadLastUserTitles(db, *auth,
			[]*spacetime.Space{space})
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	err = spacetime.LoadOriginalTitles(db,
		[]*spacetime.Space{space})
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	err = spacetime.LoadTopTitles(db,
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
			&id, 0, spacetime.MaxSubspacesPageLimit, nil, nil)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}

		if auth != nil {
			err = spacetime.LoadLastUserTitles(db, *auth,
				content)
			if err != nil {
				logging.LogError(r, auth, err)
				return nil, http.StatusInternalServerError
			}
		}

		err = spacetime.LoadOriginalTitles(db, content)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}

		err = spacetime.LoadTopTitles(db, content)
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

			err = spacetime.LoadOriginalTitles(db, path)
			if err != nil {
				logging.LogError(r, auth, err)
				return nil, http.StatusInternalServerError
			}

			space.ParentPath = &path
		}
	}

	return space, http.StatusOK

}
