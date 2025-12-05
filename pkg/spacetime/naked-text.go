package spacetime

import (
	"database/sql"

	"spacetime/pkg/utils/ajax"
)

type NakedTextDelta struct {
	Timestamp uint `json:"ts"`

	Event string `json:"e"` // "change", "select"

	// added text (or blank if removed)
	Text *string `json:"t,omitempty"`

	// selections
	SelectStart *uint `json:"ss,omitempty"`
	SelectEnd   *uint `json:"se,omitempty"`
}

type NakedText []NakedTextDelta

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

		hasEvent := delta.Event != ""
		if !hasEvent {
			return false
		}
		if delta.Event != "change" && delta.Event != "select" {
			return false
		}

		hasAddText := delta.Text != nil

		hasSelect := delta.SelectStart != nil && delta.SelectEnd != nil
		hasPartialSelect := (delta.SelectStart != nil || delta.SelectEnd != nil) && !hasSelect

		if !hasAddText && !hasSelect {
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
