package spacetime

import (
	"encoding/json"
	"fmt"
	"time"
)

type Space struct {
	ID        uint      `json:"id"`
	ParentID  *uint     `json:"parentId"`
	SpaceType string    `json:"spaceType"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy uint      `json:"createdBy"`

	IsPinned bool `json:"isPinned"`

	// TotalSubspaces uint `json:"totalSubspaces"`
	CheckinCount uint `json:"checkinCount"`

	Text  *string  `json:"text,omitempty"` // tag, text
	Title **string `json:"title,omitempty"`

	Label *string `json:"label,omitempty"`

	FinalText    *string    `json:"finalText,omitempty"`    // naked-text
	ReplayData   *NakedText `json:"recording,omitempty"`    // naked-text
	HasRecording *bool      `json:"hasRecording,omitempty"` // naked-text
	StartedAt    *time.Time `json:"startedAt,omitempty"`    // naked-text

	StreamClosedAt **time.Time `json:"streamClosedAt,omitempty"`
	StreamTexts    *[]*Space   `json:"streamTexts,omitempty"`

	LinkSpaceID **uint  `json:"linkSpaceId,omitempty"` // space-link
	LinkSpace   **Space `json:"linkSpace,omitempty"`   // space-link

	AuthorHandle      **string `json:"authorHandle,omitempty"`
	AuthorDisplayName *string  `json:"authorDisplayName,omitempty"`

	UserBookmark      *bool      `json:"userBookmark,omitempty"`
	BookmarkCreatedAt *time.Time `json:"bookmarkCreatedAt,omitempty"`
	IncludedInParent  *bool      `json:"includedInParent,omitempty"`

	Tags      *[]*Space `json:"tags,omitempty"`
	Subspaces *[]*Space `json:"subspaces,omitempty"`

	ParentPath *[]*Space `json:"parentPath,omitempty"`
}

// simple search filter for now
type SpaceFilter struct {
	Mode        string     `json:"mode"`             // "top-subspaces", "most-recent", "pinned"
	Date        *time.Time `json:"date,omitempty"`   // null for 'now'
	Window      *string    `json:"window,omitempty"` // "day", "week", "month", "year"; null for all-time
	PinnedFirst bool       `json:"pinnedFirst"`
}

func (f *SpaceFilter) Clone() *SpaceFilter {
	if f == nil {
		return nil
	}
	bytes, _ := json.Marshal(f)
	var clone SpaceFilter
	json.Unmarshal(bytes, &clone)
	return &clone
}

const SpaceFilterModeTopSubspaces = "top-subspaces"
const SpaceFilterModeMostRecent = "most-recent"
const SpaceFilterModePinned = "pinned"

func ParseSpaceFilter(jsonString string) (*SpaceFilter, error) {
	if jsonString == "" {
		// default filter
		return &SpaceFilter{Mode: SpaceFilterModeTopSubspaces}, nil
	}

	var filter SpaceFilter
	err := json.Unmarshal([]byte(jsonString), &filter)
	if err != nil {
		return nil, err
	}

	// Verify mode
	switch filter.Mode {
	case SpaceFilterModeTopSubspaces, SpaceFilterModeMostRecent, SpaceFilterModePinned:
		// valid modes
	default:
		return nil, fmt.Errorf("invalid filter mode: %s", filter.Mode)
	}

	return &filter, nil
}

type TypesFilter struct {
	Types   []string
	Exclude bool
}
