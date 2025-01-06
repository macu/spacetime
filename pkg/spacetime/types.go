package spacetime

import "time"

type Space struct {
	ID        uint      `json:"id"`
	ParentID  *uint     `json:"parentId"`
	SpaceType string    `json:"spaceType"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy uint      `json:"createdBy"`

	TotalSubspaces uint `json:"totalSubspaces"`

	Text *string `json:"text,omitempty"` // tag, title, text

	FinalText  *string    `json:"finalText,omitempty"`  // naked-text
	ReplayData *NakedText `json:"replayData,omitempty"` // naked-text

	StreamClosedAt **time.Time `json:"streamClosedAt,omitempty"`
	StreamTexts    *[]*Space   `json:"streamTexts,omitempty"`

	LinkSpaceID **uint  `json:"linkSpaceId,omitempty"` // space-link
	LinkSpace   **Space `json:"linkSpace,omitempty"`   // space-link

	AuthorHandle      **string `json:"authorHandle,omitempty"`
	AuthorDisplayName *string  `json:"authorDisplayName,omitempty"`

	UserBookmark *bool `json:"userBookmark,omitempty"`

	UserTitle     **Space   `json:"userTitle,omitempty"` // last title by user check-in
	OriginalTitle **Space   `json:"originalTitle,omitempty"`
	TopTitle      **Space   `json:"topTitle,omitempty"`
	TopTags       *[]*Space `json:"topTags,omitempty"`
	TopSubspaces  *[]*Space `json:"topSubspaces,omitempty"`

	ParentPath *[]*Space `json:"parentPath,omitempty"`
}
