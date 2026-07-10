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
	"spacetime/pkg/utils/db"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

func AjaxCreateBranch(conn *sql.DB, auth ajax.Auth,
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

	// whether to pin to parent
	pin := types.AtoBool(r.FormValue("pin"))
	if pin && parentId == nil {
		// cannot pin if no parent
		return nil, http.StatusBadRequest
	}

	formTags, err := spacetime.UnmarshalFormTags(r)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusBadRequest
	}

	// check throttle
	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(conn, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if blocked {
		return nil, http.StatusTooManyRequests
	}

	// check pinning permissions
	if pin {
		if allowed, err := spacetime.AllowsPinToParent(conn, auth, auth.UserID, *parentId); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !allowed {
			return nil, http.StatusForbidden
		}
	}
	if len(formTags.PinnedTags) > 0 {
		if allowed, err := spacetime.AllowsPinningUnderCreateSpace(conn, auth, auth.UserID, *parentId, spacetime.SpaceTypeBranch); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !allowed {
			return nil, http.StatusForbidden
		}
	}

	if parentId != nil {
		// check if parent exists
		if exists, err := spacetime.CheckSpaceExists(conn, *parentId); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !exists {
			return nil, http.StatusNotFound
		}
	}

	// check if the given label exists under the given parent
	if exists, err := spacetime.CheckBranchLabelExists(conn, parentId, label); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		// TODO Return the existing branch space
		return nil, http.StatusConflict
	}

	var space *spacetime.Space

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		var err error

		space, err = spacetime.CreateBranchSpace(tx, auth, parentId, label)
		if err != nil {
			return fmt.Errorf("create branch space: %w", err)
		}

		tags, err := spacetime.BatchCreateTags(tx, auth, space.ID, formTags)
		if err != nil {
			return fmt.Errorf("batch create tags: %w", err)
		}
		space.Tags = &tags

		if pin {
			if err := spacetime.PinSpace(tx, auth, space); err != nil {
				return fmt.Errorf("pin space: %w", err)
			}
		}

		return nil

	})

	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateLinkSpace(conn *sql.DB, auth ajax.Auth,
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

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(conn, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(conn, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	// check if link already exists
	if exists, err := spacetime.CheckLinkSpaceExists(conn, parentID, spaceID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		// TODO Return the existing link
		return nil, http.StatusConflict
	}

	space, err := spacetime.CreateSpaceLink(conn, auth, parentID, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateCheckin(conn *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	// parent required
	parentID, err := types.AtoUint(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// check throttle
	if blocked, err := spacetime.CheckCreateCheckinThrottleBlock(conn, auth, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if blocked {
		return nil, http.StatusTooManyRequests
	}

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(conn, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	err = spacetime.CreateCheckin(conn, auth, parentID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusCreated

}

func AjaxCreateTagSpace(conn *sql.DB, auth ajax.Auth,
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

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(conn, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	pin := types.AtoBool(r.FormValue("pin"))
	if pin && parentID == 0 {
		// cannot pin if no parent
		return nil, http.StatusBadRequest
	}

	// check if parent exists
	if exists, err := spacetime.CheckSpaceExists(conn, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	// check if tag exists
	if exists, err := spacetime.CheckTagExists(conn, parentID, tag); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if exists {
		// TODO Return the existing tag
		return nil, http.StatusConflict
	}

	// check owner permissions
	if pin {
		if allowed, err := spacetime.AllowsPinToParent(conn, auth, auth.UserID, parentID); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !allowed {
			return nil, http.StatusForbidden
		}
	}

	var space *spacetime.Space

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		var err error

		space, err = spacetime.CreateTag(tx, auth, parentID, tag)
		if err != nil {
			return fmt.Errorf("create tag space: %w", err)
		}

		if pin {
			if err := spacetime.PinSpace(tx, auth, space); err != nil {
				return fmt.Errorf("pin space: %w", err)
			}
		}

		return nil

	})

	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateTextSpace(conn *sql.DB, auth ajax.Auth,
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

	// whether to pin to parent
	pin := types.AtoBool(r.FormValue("pin"))
	if pin && parentID == 0 {
		// cannot pin if no parent
		return nil, http.StatusBadRequest
	}

	formTags, err := spacetime.UnmarshalFormTags(r)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusBadRequest
	}

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(conn, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	// Check if parent exists
	if exists, err := spacetime.CheckSpaceExists(conn, parentID); err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	} else if !exists {
		return nil, http.StatusNotFound
	}

	// check pinning permissions
	if pin {
		if allowed, err := spacetime.AllowsPinToParent(conn, auth, auth.UserID, parentID); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !allowed {
			return nil, http.StatusForbidden
		}
	}
	if len(formTags.PinnedTags) > 0 {
		if allowed, err := spacetime.AllowsPinningUnderCreateSpace(conn, auth, auth.UserID, parentID, spacetime.SpaceTypeText); err != nil {
			logging.LogError(r, &auth, err)
			return nil, http.StatusInternalServerError
		} else if !allowed {
			return nil, http.StatusForbidden
		}
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

	var space *spacetime.Space

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		var err error

		space, err = spacetime.CreateText(tx, auth, parentID, text, titlePtr, recordingData, startedAtTime)
		if err != nil {
			return fmt.Errorf("create text space: %w", err)
		}

		if pin {
			if err = spacetime.PinSpace(tx, auth, space); err != nil {
				return fmt.Errorf("pin space: %w", err)
			}
		}

		tags, err := spacetime.BatchCreateTags(tx, auth, space.ID, formTags)
		if err != nil {
			return fmt.Errorf("batch create tags: %w", err)
		}
		space.Tags = &tags

		return nil

	})

	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return space, http.StatusCreated

}

func AjaxCreateStreamOfConsciousnessSpace(conn *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	blocked, err := spacetime.CheckCreateSpaceThrottleBlock(conn, auth)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}
	if blocked {
		return nil, http.StatusTooManyRequests
	}

	return nil, http.StatusNotImplemented

}

func AjaxCloseStreamOfConsciousnessSpace(conn *sql.DB, auth ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	return nil, http.StatusNotImplemented

}
