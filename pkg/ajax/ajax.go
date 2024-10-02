package ajax

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"treetime/pkg/auth"
	"treetime/pkg/env"
	"treetime/pkg/user"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
)

var ajaxHandlersAuthOptional = map[string]map[string]ajax.AjaxRouteAuthOptional{
	http.MethodGet: {
		// "/ajax/test": ajaxTest,
	},
	http.MethodPost: {
		"/ajax/login": auth.AjaxLogin,
	},
}

var ajaxHandlersAuthRequired = map[string]map[string]ajax.AjaxRouteAuthRequired{
	http.MethodGet: {
		// "/ajax/auth":     ajaxLoadAuthenticatedUser,
		// "/ajax/settings": ajaxUserSettings,
		// "/ajax/logout": ajaxLogoutHandler,
	},
	http.MethodPost: {
		// "/ajax/user/change-password": ajaxUserChangePassword,
		// "/ajax/user/save-settings":   ajaxUserSaveSettings,

		// util
		// "/ajax/render-markdown": ajaxRenderMarkdown,

		"/ajax/logout": auth.AjaxLogout,
	},
}

func AjaxHandler(db *sql.DB, userID *uint, w http.ResponseWriter, r *http.Request) {
	// var rt = NewResponseTracker(w)

	var handle = func(handler func() (interface{}, int)) {
		// Verify access to admin routes
		if strings.HasPrefix(r.URL.Path, "/ajax/admin") && (userID == nil || !user.CheckAdmin(*userID)) {
			logging.LogError(r, userID, fmt.Errorf("forbidden admin access"))
			w.WriteHeader(http.StatusForbidden)
			return
		}
		response, statusCode := handler()
		if statusCode >= 400 {
			w.WriteHeader(statusCode)
			// Send current version stamp
			w.Write([]byte("VersionStamp: " + env.GetCacheControlVersionStamp()))
			return
		}
		if response != nil {
			js, err := json.Marshal(response)
			if err != nil {
				logging.LogError(r, userID, fmt.Errorf("marshalling response: %w", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode) // WriteHeader is called after setting headers
			w.Write(js)
		} else {
			w.WriteHeader(statusCode)
		}
	}

	handlersAuthOptional, foundMethod := ajaxHandlersAuthOptional[r.Method]
	if foundMethod {
		handler, fouundPath := handlersAuthOptional[r.URL.Path]
		if fouundPath {
			handle(func() (interface{}, int) {
				return handler(db, userID, w, r)
			})
			return
		}
	}

	handlersAuthRequired, foundMethod := ajaxHandlersAuthRequired[r.Method]
	if foundMethod {
		handler, fouundPath := handlersAuthRequired[r.URL.Path]
		if fouundPath {
			if userID == nil {
				w.WriteHeader(http.StatusForbidden)
			} else {
				handle(func() (interface{}, int) {
					return handler(db, *userID, w, r)
				})
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
