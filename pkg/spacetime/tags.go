package spacetime

import (
	"database/sql"
	"fmt"
	"net/http"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
	"spacetime/pkg/utils/types"
)

// payload to create functions
type FormTags struct {
	Tags       []string `json:"tags"`
	PinnedTags []string `json:"pinnedTags"`
}

func UnmarshalFormTags(r *http.Request) (FormTags, error) {
	var formTags FormTags
	var err error

	formTags.Tags, err = types.JSONtoStringArray(r.FormValue("tags"))
	if err != nil {
		return FormTags{}, fmt.Errorf("parse tags: %w", err)
	}
	formTags.PinnedTags, err = types.JSONtoStringArray(r.FormValue("pinnedTags"))
	if err != nil {
		return FormTags{}, fmt.Errorf("parse pinned tags: %w", err)
	}

	// Normalize and validate
	for i, tag := range formTags.Tags {
		tag = types.NormalizeSpaces(tag)
		if !ValidateTag(tag) {
			return FormTags{}, fmt.Errorf("invalid tag: %s", tag)
		}
		formTags.Tags[i] = tag
	}
	for i, tag := range formTags.PinnedTags {
		tag = types.NormalizeSpaces(tag)
		if !ValidateTag(tag) {
			return FormTags{}, fmt.Errorf("invalid pinned tag: %s", tag)
		}
		formTags.PinnedTags[i] = tag
	}

	// check for duplicates within and across tags and pinnedTags
	tagSet := make(map[string]bool)
	for _, tag := range formTags.Tags {
		if _, exists := tagSet[tag]; exists {
			return FormTags{}, fmt.Errorf("duplicate tag: %s", tag)
		}
		tagSet[tag] = true
	}
	for _, tag := range formTags.PinnedTags {
		if _, exists := tagSet[tag]; exists {
			return FormTags{}, fmt.Errorf("duplicate pinned tag: %s", tag)
		}
		tagSet[tag] = true
	}

	return formTags, nil
}

func LoadExistingTag(conn *sql.DB,
	parentID uint, tag string,
) (*Space, error) {

	// Load tag space

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeTag,
		Text:      &tag,
	}

	var args = []interface{}{}

	err := conn.QueryRow(`SELECT space.id, space.created_at, space.created_by,
		CASE WHEN user_space_config.space_id IS NULL THEN FALSE ELSE TRUE END AS pinned
		FROM space
		INNER JOIN tag_space ON tag_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = tag_space.text_id
		INNER JOIN user_space_config ON user_space_config.space_id = space.id
		WHERE space.parent_id = `+db.Arg(&args, parentID)+`
		AND space.space_type = `+db.Arg(&args, SpaceTypeTag)+`
		AND unique_text.text_value = `+db.Arg(&args, tag),
		args...,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy, &space.IsPinned)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select tag_space: %w", err)
	}

	return space, nil

}

func BatchCreateTags(tx *sql.Tx, auth ajax.Auth, spaceID uint, tags FormTags) ([]*Space, error) {

	var tagSpaces = []*Space{}

	for _, tag := range tags.PinnedTags {
		tagSpace, err := CreateTag(tx, auth, spaceID, tag)
		if err != nil {
			return nil, fmt.Errorf("create pinned tag: %w", err)
		}
		tagSpaces = append(tagSpaces, tagSpace)

		err = PinSpace(tx, auth, tagSpace)
		if err != nil {
			return nil, fmt.Errorf("pin tag: %w", err)
		}

		tagSpace.IsPinned = true
	}

	for _, tag := range tags.Tags {
		createdTag, err := CreateTag(tx, auth, spaceID, tag)
		if err != nil {
			return nil, fmt.Errorf("create tag: %w", err)
		}
		tagSpaces = append(tagSpaces, createdTag)
	}

	return tagSpaces, nil

}

func CreateTag(tx *sql.Tx, auth ajax.Auth, parentID uint, tag string) (*Space, error) {

	// Load unique_text ID
	// Check for existing tag space under parent
	// Create tag space if not exists
	// Check-in on tag space

	if !ValidateTag(tag) {
		return nil, fmt.Errorf("invalid tag: %s", tag)
	}

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeTag,
		Text:      &tag,
	}

	uniqueTextId, err := GetOrCreateUniqueTextId(tx, tag)
	if err != nil {
		return nil, err
	} else if uniqueTextId == nil {
		return nil, fmt.Errorf("unique text id is nil")
	}

	// Create space
	err = CreateSpace(tx, auth, space, &parentID, SpaceTypeTag)
	if err != nil {
		return nil, err
	}

	// Create tag_space
	_, err = tx.Exec(`INSERT INTO tag_space
		(space_id, parent_id, text_id)
		VALUES ($1, $2, $3)`,
		space.ID, parentID, *uniqueTextId,
	)
	if err != nil {
		return nil, fmt.Errorf("insert tag_space: %w", err)
	}

	return space, nil

}

func LoadTopTags(conn *sql.DB, spaces []*Space,
	offset uint, limit uint,
) error {
	// Load top tags for multiple spaces

	if len(spaces) == 0 {
		return nil
	}

	if limit > MaxSubspacesPageLimit {
		limit = MaxSubspacesPageLimit
	}

	for _, space := range spaces {

		tags, err := LoadMoreTags(conn, space.ID, offset, limit)
		if err != nil {
			return err
		}

		space.TopTags = tags

	}

	return nil

}

func LoadMoreTags(conn *sql.DB, parentId uint,
	offset uint, limit uint,
) (*[]*Space, error) {

	rows, err := conn.Query(`SELECT space.id, unique_text.text_value,
		(SELECT COUNT(*) FROM checkin
			WHERE checkin.space_id = space.id) AS checkin_count,
		CASE WHEN user_space_config.space_id IS NULL THEN FALSE ELSE TRUE END AS pinned
		FROM space
		INNER JOIN tag_space ON tag_space.space_id = space.id
		INNER JOIN unique_text ON unique_text.id = tag_space.text_id
		LEFT JOIN user_space_config ON user_space_config.space_id = space.id
		WHERE space.space_type = $1
		AND space.parent_id = $2
		GROUP BY space.id, unique_text.text_value,
			user_space_config.space_id, user_space_config.order_number
		ORDER BY
			CASE WHEN user_space_config.space_id IS NULL THEN FALSE ELSE TRUE END DESC,
			CASE WHEN user_space_config.order_number IS NULL THEN 0 ELSE user_space_config.order_number END ASC,
			checkin_count DESC
		OFFSET $3
		LIMIT $4`,
		SpaceTypeTag, parentId, offset, limit,
	)

	if err != nil {
		return nil, fmt.Errorf("loading top tags: %w", err)
	}

	defer rows.Close()

	var tags = []*Space{}

	for rows.Next() {
		var spaceID uint
		var text string
		var checkinCount uint
		var pinned bool
		err = rows.Scan(&spaceID, &text, &checkinCount, &pinned)
		if err != nil {
			return nil, fmt.Errorf("scanning top tags: %w", err)
		}
		var tag = &Space{
			ID:           spaceID,
			SpaceType:    SpaceTypeTag,
			Text:         &text,
			CheckinCount: checkinCount,
			IsPinned:     pinned,
		}
		tags = append(tags, tag)
	}

	return &tags, nil

}
