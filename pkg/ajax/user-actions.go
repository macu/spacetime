package ajax

import (
	"database/sql"
	"net/http"

	"spacetime/pkg/spacetime"
	"spacetime/pkg/user"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

func AjaxBookmark(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	spaceID, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	bookmark := types.AtoBool(r.FormValue("bookmark"))

	// check if space exists
	if exists, err := spacetime.CheckSpaceExists(db, spaceID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	if err = user.BookmarkSpace(db, auth.UserID, spaceID, bookmark); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusOK

}

func AjaxPinSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	spaceID, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	pinned := types.AtoBool(r.FormValue("pinned"))

	// check if space exists
	space, err := spacetime.GetSpace(db, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if space == nil {
		return nil, http.StatusNotFound
	} else if space.ParentID == nil {
		return nil, http.StatusBadRequest
	}

	if allowed, err := spacetime.AllowsPinToParent(db, auth, auth.UserID, *space.ParentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !allowed {
		return nil, http.StatusForbidden
	}

	if pinned {
		err = spacetime.PinSpace(db, auth, space)
		if err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		}
	} else {
		err = spacetime.UnpinSpace(db, &auth, spaceID)
		if err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return nil, http.StatusOK

}

func AjaxMovePinnedSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	spaceID, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	newIndex, err := types.AtoUint(r.FormValue("newIndex"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// check if space exists
	space, err := spacetime.GetSpace(db, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if space == nil {
		return nil, http.StatusNotFound
	}

	if allowed, err := spacetime.AllowsPinToParent(db, auth, auth.UserID, *space.ParentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !allowed {
		return nil, http.StatusForbidden
	}

	err = spacetime.ReorderPin(db, &auth, space, newIndex)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusOK

}
