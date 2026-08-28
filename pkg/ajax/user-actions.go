package ajax

import (
	"database/sql"
	"net/http"
	"regexp"

	authFunctions "spacetime/pkg/auth"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"

	"spacetime/pkg/spacetime"
)

func AjaxVote(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	spaceID, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	voteValue, err := types.AtoInt(r.FormValue("voteValue"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	if voteValue < -1 || voteValue > 1 {
		return nil, http.StatusBadRequest
	}

	// check if space exists
	if exists, err := spacetime.CheckSpaceExists(db, spaceID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	voteSum, err := spacetime.AddVote(db, auth.UserID, spaceID, voteValue)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return struct {
		VoteSum int `json:"voteSum"`
	}{
		VoteSum: voteSum,
	}, http.StatusOK

}

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

	if err = spacetime.BookmarkSpace(db, auth.UserID, spaceID, bookmark); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusOK

}

func AjaxBookmarks(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	offset, err := types.AtoUint(r.FormValue("offset"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	limit, err := types.AtoUint(r.FormValue("limit"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	includeParentPath := types.AtoBool(r.FormValue("includeParentPath"))
	includeTags := types.AtoBool(r.FormValue("includeTags"))

	includeLinkedInParentId, err := types.AtoUintNilIfEmpty(r.FormValue("includeLinkedInParentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	bookmarks, err := spacetime.GetBookmarkedSpaces(db, auth, offset, limit,
		includeParentPath, includeTags, includeLinkedInParentId)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return bookmarks, http.StatusOK

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

func AjaxUserSpaceId(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	id := r.FormValue("id")

	if id == "" {
		return nil, http.StatusBadRequest
	}

	var userID uint
	var err error

	// if numeric, look up user by ID
	isNumeric := regexp.MustCompile(`^\d+$`).MatchString(id)
	if isNumeric {
		userID, err = types.AtoUint(id)
		if err != nil {
			return nil, http.StatusBadRequest
		}
	} else {
		// otherwise, look up user by handle
		user, err := authFunctions.GetUserByHandle(db, id)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		} else if user == nil {
			return nil, http.StatusNotFound
		}
		userID = user.ID
	}

	spaceID, err := spacetime.GetUserSpaceID(db, userID)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	} else if spaceID == 0 {
		return nil, http.StatusNotFound
	}

	return struct {
		SpaceID uint `json:"spaceId"`
	}{
		SpaceID: spaceID,
	}, http.StatusOK

}
