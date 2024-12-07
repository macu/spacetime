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

	if len(title) > spacetime.TitleMaxLength {
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
		space.LastUserTitle = &titleSpace
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

	return nil, http.StatusNotImplemented

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

	return nil, http.StatusNotImplemented

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

	return nil, http.StatusNotImplemented

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

	return nil, http.StatusNotImplemented

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
