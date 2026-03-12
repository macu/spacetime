package ajax

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

func AjaxCreateBranch(db *sql.DB, auth ajax.Auth,
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
	if exists, err := spacetime.CheckBranchLabelExists(db, parentId, label); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		// TODO Return the existing branch space
		return nil, http.StatusConflict
	}

	space, err := spacetime.CreateBranchSpace(db, auth, parentId, label)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
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

	// check if link already exists
	if exists, err := spacetime.CheckLinkSpaceExists(db, parentID, spaceID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		// TODO Return the existing link
		return nil, http.StatusConflict
	}

	space, err := spacetime.CreateSpaceLink(db, auth, parentID, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
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

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(db, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	err = spacetime.CreateCheckin(db, auth, parentID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusCreated

}

func AjaxCreateTagSpace(db *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	tag := types.NormalizeSpaces(r.FormValue("tag"))
	if tag == "" || !spacetime.ValidateTag(tag) {
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
		// TODO Return the existing tag
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
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// title optional
	var titlePtr *string
	title := types.NormalizeSpaces(r.FormValue("title"))
	if title != "" {
		titlePtr = &title
		if !spacetime.ValidateTitle(title) {
			return nil, http.StatusBadRequest
		}
	}

	text := strings.TrimSpace(r.FormValue("text"))
	if text == "" || !spacetime.ValidateText(text) {
		return nil, http.StatusBadRequest
	}

	// Normalize line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(db, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// Check if parent exists
	if exists, err := spacetime.CheckSpaceExists(db, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	recording := r.FormValue("recording")
	var recordingData *spacetime.NakedText
	var startedAtTime *time.Time
	if recording != "" {
		err = json.Unmarshal([]byte(recording), &recordingData)
		if err != nil {
			return nil, http.StatusBadRequest
		}

		startedAt, err := types.AtoInt64(r.FormValue("startedAt")) // Unix timestamp
		if err != nil {
			return nil, http.StatusBadRequest
		}
		startedAtTimeValue := time.Unix(startedAt, 0)
		startedAtTime = &startedAtTimeValue

		if !spacetime.ValidateNakedText(*recordingData, text) {
			logging.LogError(r, &auth, fmt.Errorf("invalid naked text recording"))
			return nil, http.StatusBadRequest
		}
	}

	space, err := spacetime.CreateText(db, auth, parentID, text, titlePtr, recordingData, startedAtTime)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

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

	return nil, http.StatusNotImplemented

}
