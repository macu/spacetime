package ajax

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"treetime/pkg/treetime"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/types"
)

func AjaxFindExistingNodes(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var title = strings.TrimSpace(r.FormValue("title"))
	var class = strings.TrimSpace(r.FormValue("class"))

	if title == "" || class == "" {
		return nil, http.StatusBadRequest
	}

	return nil, http.StatusNotImplemented

}

func AjaxCreateNode(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var title = strings.TrimSpace(r.FormValue("title"))
	var class = strings.TrimSpace(r.FormValue("class"))
	var content = strings.TrimSpace(r.FormValue("content"))

	if title == "" || class == "" {
		return nil, http.StatusBadRequest
	}

	var err error

	var parentIDString = r.FormValue("parentID")
	var parentID uint

	if parentIDString != "" {
		parentID, err = types.AtoUint(parentIDString)
		if err != nil {
			return nil, http.StatusBadRequest
		}
	}

	createClass, err := treetime.StringToClass(class)

	if err != nil {
		return nil, http.StatusBadRequest
	}

	node, err := treetime.CreateNode(db, auth, parentID, createClass, title, content)

	if err != nil {
		logging.LogError(r, &auth, fmt.Errorf("creating node: %w", err))
		return nil, http.StatusInternalServerError
	}

	return node, http.StatusOK

}

func AjaxSaveNode(db *sql.DB, auth ajax.Auth, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	return nil, http.StatusNotImplemented

}
