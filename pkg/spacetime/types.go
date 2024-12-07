package spacetime

import "time"

type Space struct {
	ID           uint      `json:"id"`
	SpaceType    string    `json:"spaceType"`
	CreatedAt    time.Time `json:"createdAt"`
	CreatedBy    uint      `json:"createdBy"`
	CheckinTotal uint      `json:"checkinTotal"`

	Text *string `json:"text,omitempty"` // tag, title, text

	FinalText  *string                 `json:"finalText,omitempty"`  // naked-text
	ReplayData *map[string]interface{} `json:"replayData,omitempty"` // naked-text

	CheckinSpaceID *uint   `json:"checkinSpaceId,omitempty"` // checkin
	CheckinSpace   **Space `json:"checkinSpace,omitempty"`   // checkin

	UserBookmark *bool `json:"userBookmark,omitempty"`

	LastUserTitle **Space   `json:"lastUserTitle,omitempty"` // last title by user check-in
	TopTitles     *[]*Space `json:"topTitles,omitempty"`
	TopTags       *[]*Space `json:"topTags,omitempty"`
	TopContent    *[]*Space `json:"topContent,omitempty"`

	AllTimeCheckinCount     *uint `json:"checkinCountAll,omitempty"`
	LastTwentyFourHourCount *uint `json:"checkinCount24,omitempty"`

	StreamClosedAt **time.Time `json:"streamClosedAt,omitempty"`
	StreamTexts    *[]*Space   `json:"streamTexts,omitempty"`
}
