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

	exists, err := spacetime.CheckSpaceExists(db, spaceID)
	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	if !exists {
		return nil, http.StatusNotFound
	}

	bookmark := types.AtoBool(r.FormValue("bookmark"))

	err = user.BookmarkSpace(db, auth.UserID, spaceID, bookmark)

	if err != nil {
		logging.LogError(r, &auth, err)
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusOK

}
