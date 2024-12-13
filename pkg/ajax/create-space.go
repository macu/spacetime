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

func AjaxCreateEmptySpace(db *sql.DB, auth ajax.Auth,
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

	parentId, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	title := strings.TrimSpace(r.FormValue("title")) // optional

	if title != "" && !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.CreateEmptySpace(db, auth, parentId)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	if title != "" {
		titleSpace, err := spacetime.CreateTitleCheckin(db, auth, space.ID, title)
		if err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		}
		space.UserTitles = &[]*spacetime.Space{titleSpace}
		space.TopTitles = &[]*spacetime.Space{titleSpace}
		space.TopSubspaces = &[]*spacetime.Space{titleSpace}
		space.TotalSubspaces = 1
	}

	return space, http.StatusCreated

}

func AjaxCreateCheckinSpace(db *sql.DB, auth ajax.Auth,
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

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// space optional
	spaceID, err := types.AtoUintNilIfEmpty(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.CreateCheckin(db, auth, parentID, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateTitleSpace(db *sql.DB, auth ajax.Auth,
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

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	title := strings.TrimSpace(r.FormValue("title"))

	if !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.CreateTitleCheckin(db, auth, parentID, title)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusNotImplemented

}

func AjaxCreateTagSpace(db *sql.DB, auth ajax.Auth,
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

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	tag := strings.TrimSpace(r.FormValue("tag"))

	if !spacetime.ValidateTag(tag) {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.CreateTagCheckin(db, auth, parentID, tag)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusNotImplemented

}

func AjaxCreateTextSpace(db *sql.DB, auth ajax.Auth,
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

	// parent optional
	parentID, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	title := strings.TrimSpace(r.FormValue("title"))

	if title != "" && !spacetime.ValidateTitle(title) {
		return nil, http.StatusBadRequest
	}

	text := strings.TrimSpace(r.FormValue("text"))

	if !spacetime.ValidateText(text) {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.CreateTextCheckin(db, auth, parentID, text)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	_, err = spacetime.CreateTitleCheckin(db, auth, space.ID, title)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusNotImplemented

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
