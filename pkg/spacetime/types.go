package spacetime

import "time"

type Space struct {
	ID        uint      `json:"id"`
	SpaceType string    `json:"spaceType"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy uint      `json:"createdBy"`

	TotalSubspaces uint `json:"totalSubspaces"`

	Text *string `json:"text,omitempty"` // tag, title, text

	FinalText  *string    `json:"finalText,omitempty"`  // naked-text
	ReplayData *NakedText `json:"replayData,omitempty"` // naked-text

	StreamClosedAt **time.Time `json:"streamClosedAt,omitempty"`
	StreamTexts    *[]*Space   `json:"streamTexts,omitempty"`

	CheckinSpaceID **uint  `json:"checkinSpaceId,omitempty"` // checkin
	CheckinSpace   **Space `json:"checkinSpace,omitempty"`   // checkin

	AuthorHandle      **string `json:"authorHandle,omitempty"`
	AuthorDisplayName *string  `json:"authorDisplayName,omitempty"`

	UserBookmark *bool `json:"userBookmark,omitempty"`

	UserTitles   *[]*Space `json:"userTitles,omitempty"` // last titles by user check-in
	TopTitles    *[]*Space `json:"topTitles,omitempty"`
	TopTags      *[]*Space `json:"topTags,omitempty"`
	TopSubspaces *[]*Space `json:"topSubspaces,omitempty"`
}
