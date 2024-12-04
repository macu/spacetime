package spacetime

import "time"

type Space struct {
	ID        uint      `json:"id"`
	SpaceType string    `json:"space_type"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"`

	Text *string `json:"text,omitempty"` // tag, title, text

	FinalText  *string                 `json:"final_text,omitempty"`  // naked-text
	ReplayData *map[string]interface{} `json:"replay_data,omitempty"` // naked-text

	CheckinSpaceID *uint   `json:"checkin_space_id,omitempty"` // checkin
	CheckinSpace   **Space `json:"checkin_space,omitempty"`    // checkin

	UserBookmark *bool `json:"user_bookmark,omitempty"`

	BookmarkedTitles *[]*Space `json:"tagged_titles,omitempty"`
	TopTitles        *[]*Space `json:"top_titles,omitempty"`
	TopTags          *[]*Space `json:"top_tags,omitempty"`
	TopContent       *[]*Space `json:"top_content,omitempty"`

	AllTimeCheckinCount     *uint `json:"checkin_count_all,omitempty"`
	LastTwentyFourHourCount *uint `json:"checkin_count_24,omitempty"`
}
