package ajax

import (
	"database/sql"
	"net/http"

	"spacetime/pkg/spacetime"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/types"
)

// load space and optionally path/subspaces/tags
func AjaxLoadSpace(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	includeParentPath := types.AtoBool(r.FormValue("includeParentPath"))

	includeTags := types.AtoBool(r.FormValue("includeTags"))

	filter, err := spacetime.ParseSpaceFilter(r.FormValue("filter"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.LoadSpace(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	var subspacesTypesFilter *spacetime.TypesFilter
	if filter.Mode != spacetime.SpaceFilterModePinned {
		// include tags when viewing pinned mode
		// otherwise exclude tags
		subspacesTypesFilter = &spacetime.TypesFilter{
			Exclude: true,
			Types:   []string{spacetime.SpaceTypeTag},
		}
	}

	subspaces, err := spacetime.LoadSubspaces(db, auth,
		&id, 0, spacetime.MaxSubspacesPageLimit, filter, subspacesTypesFilter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if includeTags {
		// include tags on parent and subspaces
		err = spacetime.LoadTags(db, auth,
			append([]*spacetime.Space{space}, subspaces...),
			0, spacetime.DefaultTagsLimit, filter)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	space.Subspaces = &subspaces

	if includeParentPath {
		if space.ParentID == nil {
			space.ParentPath = &[]*spacetime.Space{}
		} else {
			path, err := spacetime.LoadParentPath(db, auth, *space.ParentID)
			if err != nil {
				logging.LogError(r, auth, err)
				return nil, http.StatusInternalServerError
			}

			space.ParentPath = &path
		}
	}

	return space, http.StatusOK

}

// Load tags and subspaces when filters change
func AjaxReloadSpace(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	filter, err := spacetime.ParseSpaceFilter(r.FormValue("filter"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	includeTags := types.AtoBool(r.FormValue("includeTags"))

	exists, err := spacetime.CheckSpaceExists(db, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}
	if !exists {
		return nil, http.StatusNotFound
	}

	var subspacesTypesFilter *spacetime.TypesFilter
	if filter.Mode != spacetime.SpaceFilterModePinned {
		// include tags when viewing pinned mode
		// otherwise exclude tags
		subspacesTypesFilter = &spacetime.TypesFilter{
			Exclude: true,
			Types:   []string{spacetime.SpaceTypeTag},
		}
	}

	subspaces, err := spacetime.LoadSubspaces(db, auth,
		&id, 0, spacetime.MaxSubspacesPageLimit, filter, subspacesTypesFilter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	tagTypesFilter := &spacetime.TypesFilter{
		Types: []string{spacetime.SpaceTypeTag},
	}

	tags, err := spacetime.LoadSubspaces(db, auth,
		&id, 0, spacetime.MaxSubspacesPageLimit, filter, tagTypesFilter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if includeTags {
		// load tags on subspaces
		err = spacetime.LoadTags(db, auth, subspaces, 0,
			spacetime.DefaultTagsLimit, filter)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return struct {
		Subspaces []*spacetime.Space `json:"subspaces"`
		Tags      []*spacetime.Space `json:"tags"`
	}{
		Subspaces: subspaces,
		Tags:      tags,
	}, http.StatusOK

}

// load non-tag subspaces
func AjaxLoadSubspaces(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	parentId, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	offset, err := types.AtoUint(r.FormValue("offset"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	limit, err := types.AtoUint(r.FormValue("limit"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	includeTags := types.AtoBool(r.FormValue("includeTags"))

	filter, err := spacetime.ParseSpaceFilter(r.FormValue("filter"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	var subspacesTypesFilter *spacetime.TypesFilter
	if filter.Mode != spacetime.SpaceFilterModePinned {
		// include tags when viewing pinned mode
		// otherwise exclude tags
		subspacesTypesFilter = &spacetime.TypesFilter{
			Exclude: true,
			Types:   []string{spacetime.SpaceTypeTag},
		}
	}

	subspaces, err := spacetime.LoadSubspaces(db, auth,
		parentId, offset, limit, filter, subspacesTypesFilter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if includeTags {
		err = spacetime.LoadTags(db, auth,
			subspaces, 0, spacetime.DefaultTagsLimit, filter)
		if err != nil {
			logging.LogError(r, auth, err)
			return nil, http.StatusInternalServerError
		}
	}

	return subspaces, http.StatusOK

}

// load tags
func AjaxLoadTags(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	parentId, err := types.AtoUintNilIfEmpty(r.FormValue("parentId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	offset, err := types.AtoUint(r.FormValue("offset"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	limit, err := types.AtoUint(r.FormValue("limit"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	filter, err := spacetime.ParseSpaceFilter(r.FormValue("filter"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	// include only tags
	subspacesTypesFilter := &spacetime.TypesFilter{
		Types: []string{spacetime.SpaceTypeTag},
	}

	subspaces, err := spacetime.LoadSubspaces(db, auth,
		parentId, offset, limit, filter, subspacesTypesFilter)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return subspaces, http.StatusOK

}

func AjaxLoadTextSpaceRecording(db *sql.DB, auth *ajax.Auth,
	w http.ResponseWriter, r *http.Request,
) (interface{}, int) {

	id, err := types.AtoUint(r.FormValue("spaceId"))
	if err != nil {
		return nil, http.StatusBadRequest
	}

	space, err := spacetime.LoadSpace(db, auth, id)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	if space.SpaceType != spacetime.SpaceTypeText {
		return nil, http.StatusBadRequest
	}

	err = spacetime.LoadTextRecording(db, space)
	if err != nil {
		logging.LogError(r, auth, err)
		return nil, http.StatusInternalServerError
	}

	return space.ReplayData, http.StatusOK
}
