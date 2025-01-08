package spacetime

import (
	"database/sql"
	"encoding/json"

	"spacetime/pkg/utils/ajax"
)

type NakedTextDelta struct {
	Timestamp uint `json:"ts"`

	// key presses (like backspace)
	Key *uint `json:"k,omitempty"`

	// added text (one char at a time)
	AddText *rune `json:"t,omitempty"`

	// cursor positioning
	Cursor *uint `json:"c,omitempty"`

	// selections
	SelectStart *uint `json:"ss,omitempty"`
	SelectEnd   *uint `json:"se,omitempty"`
}

type NakedText []NakedTextDelta

func (d *NakedTextDelta) MarshalJSON() ([]byte, error) {
	type Alias NakedTextDelta
	return json.Marshal(&struct {
		AddText *string `json:"t,omitempty"`
		*Alias
	}{
		AddText: func() *string {
			if d.AddText != nil {
				s := string(*d.AddText)
				return &s
			}
			return nil
		}(),
		Alias: (*Alias)(d),
	})
}

func (d *NakedTextDelta) UnmarshalJSON(data []byte) error {
	type Alias NakedTextDelta
	aux := &struct {
		AddText *string `json:"t,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if aux.AddText != nil && len(*aux.AddText) > 0 {
		r := []rune(*aux.AddText)
		d.AddText = &r[0]
	}
	return nil
}

func ValidateNakedText(text NakedText) bool {

	// Ensure has count
	if len(text) == 0 || len(text) > NakedTextMaxDeltas {
		return false
	}

	// Ensure first delta at timestamp 0
	if text[0].Timestamp != 0 {
		return false
	}

	// Ensure timestamps increment
	for i := 1; i < len(text); i++ {
		if text[i].Timestamp >= text[i-1].Timestamp {
			return false
		}
	}

	// Ensure full data is available for each type of delta
	for _, delta := range text {

		hasAddText := delta.AddText != nil
		hasKey := delta.Key != nil
		hasCursor := delta.Cursor != nil
		hasSelect := delta.SelectStart != nil && delta.SelectEnd != nil
		hasPartialSelect := (delta.SelectStart != nil || delta.SelectEnd != nil) && !hasSelect

		if !hasAddText && !hasKey && !hasCursor && !hasSelect {
			return false
		}

		if (hasAddText || hasKey || hasSelect) && hasCursor {
			return false
		}

		if hasPartialSelect {
			return false
		}

	}

	return true

}

func CreateNakedText(conn *sql.DB, auth ajax.Auth,
	parentID uint, finalText, replayData string,
) (*Space, error) {

	// Create naked text space with given replay data

	return nil, nil

}
