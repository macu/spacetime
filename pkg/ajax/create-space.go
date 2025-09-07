package ajax

import (
	"database/sql"
	"net/http"
	"strings"

	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

func AjaxCreateSubspace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	parentId, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	label := types.NormalizeSpaces(r.FormValue("label")) // label is required
	if label == "" || !spacetime.ValidateLabel(label) {
		return nil, http.StatusBadRequest
	}

	title := types.NormalizeSpaces(r.FormValue("title")) // title is optional
	if title != "" && !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	// check throttle
	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if blocked {
		return nil, http.StatusTooManyRequests
	}

	if parentId != nil {
		// check if parent exists
		if exists, err := spacetime.CheckSpaceExists(db, *parentId); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !exists {
			return nil, http.StatusNotFound
		}
	}

	// check if the given label exists under the given parent
	if exists, err := spacetime.CheckSubspaceLabelExists(db, parentId, label); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		return nil, http.StatusConflict
	}

	space, err := spacetime.CreateSubspace(db, auth, parentId, label)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	if title != "" {
		_, err := spacetime.CreateTitle(db, auth, space.ID, title)
		if err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return space, http.StatusCreated

}

func AjaxCreateLinkSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// space to link required
	spaceID, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// title optional
	title := types.NormalizeSpaces(r.FormValue("title"))
	if title != "" && !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// check if link already exists
	if exists, err := spacetime.CheckLinkSpaceExists(db, parentID, spaceID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		return nil, http.StatusConflict
	}

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(db, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	space, err := spacetime.CreateSpaceLink(db, auth, parentID, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	if title != "" {
		_, err := spacetime.CreateTitle(db, auth, space.ID, title)
		if err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return space, http.StatusCreated

}

func AjaxCreateCheckin(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// check throttle
	if blocked, err := spacetime.CheckCreateCheckinThrottleBlock(db, auth, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if blocked {
		return nil, http.StatusTooManyRequests
	}

	space, err := spacetime.CreateCheckin(db, auth, parentID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateTitleSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	title := types.NormalizeSpaces(r.FormValue("title"))
	if title == "" || !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(db, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	// check if title exists
	if exists, err := spacetime.CheckTitleExists(db, parentID, title); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		return nil, http.StatusConflict
	}

	space, err := spacetime.CreateTitle(db, auth, parentID, title)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateTagSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	tag := types.NormalizeSpaces(strings.TrimSpace(r.FormValue("tag")))
	if !spacetime.ValidateTag(tag) {
		return nil, http.StatusBadRequest
	}

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(db, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	// check if tag exists
	if exists, err := spacetime.CheckTagExists(db, parentID, tag); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		return nil, http.StatusConflict
	}

	space, err := spacetime.CreateTag(db, auth, parentID, tag)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateTextSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent optional
	parentID, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// title optional
	title := types.NormalizeSpaces(strings.TrimSpace(r.FormValue("title")))
	if title != "" && !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	text := strings.TrimSpace(r.FormValue("text"))
	if text == "" || !spacetime.ValidateText(text) {
		return nil, http.StatusBadRequest
	}

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// if parentID given, check if parent exists
	if parentID != nil {
		if exists, err := spacetime.CheckSpaceExists(db, *parentID); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !exists {
			return nil, http.StatusNotFound
		}
	}

	space, err := spacetime.CreateText(db, auth, parentID, text)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	if title != "" {
		_, err = spacetime.CreateTitle(db, auth, space.ID, title)
		if err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return space, http.StatusCreated

}

func AjaxCreateNakedTextSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	return nil, http.StatusNotImplemented

}

func AjaxCreateStreamOfConsciousnessSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	return nil, http.StatusNotImplemented

}

func AjaxCloseStreamOfConsciousnessSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	return nil, http.StatusNotImplemented

}
