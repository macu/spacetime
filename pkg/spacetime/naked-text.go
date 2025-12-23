package spacetime

import (
	"fmt"
	"spacetime/pkg/utils/logging"
	"strings"
)

type NakedTextDelta struct {
	Timestamp uint `json:"ts"`

	EventType string `json:"et"` // "change", "select", "cursor"

	// added text (or blank if removed)
	Text *string `json:"t,omitempty"`

	// modification, selection, and cursor positions
	SelectStart uint `json:"ss"`
	SelectEnd   uint `json:"se"`
}

type NakedText []NakedTextDelta

func ValidateNakedText(recording NakedText, finalText string) bool {

	// Ensure has count
	if len(recording) == 0 || len(recording) > NakedTextMaxDeltas {
		logging.LogError(nil, nil, fmt.Errorf("naked text recording has invalid number of deltas"))
		return false
	}

	// Ensure first delta at timestamp 0
	if recording[0].Timestamp != 0 {
		logging.LogError(nil, nil, fmt.Errorf("naked text first delta timestamp not zero"))
		return false
	}

	// Ensure timestamps increment
	for i := 1; i < len(recording); i++ {
		if recording[i].Timestamp < recording[i-1].Timestamp {
			logging.LogError(nil, nil, fmt.Errorf("naked text delta timestamps not increasing"))
			return false
		}
	}

	// Ensure full data is available for each type of delta
	var currentText string
	var totalDeltas uint
	for _, delta := range recording {

		totalDeltas++

		switch delta.EventType {

		case "change":
			if delta.Text == nil || len(*delta.Text) > TextMaxLength {
				logging.LogError(nil, nil, fmt.Errorf("invalid change delta text"))
				return false
			}
			// Apply change
			currentText = currentText[:delta.SelectStart] +
				*delta.Text +
				currentText[delta.SelectEnd:]
			logging.LogNotice(nil, fmt.Sprintf("Applied change delta %v, new text: %s\n", delta, currentText))
			// Check length
			if len(currentText) > TextMaxLength {
				logging.LogError(nil, nil, fmt.Errorf("text exceeds max length"))
				return false
			}
			// Count added text towards total deltas
			totalDeltas += uint(len(*delta.Text))
			if totalDeltas > NakedTextMaxDeltas {
				logging.LogError(nil, nil, fmt.Errorf("naked text exceeds max deltas"))
				return false
			}

		case "select":
			if delta.Text != nil {
				logging.LogError(nil, nil, fmt.Errorf("select delta with text"))
				return false
			}
			if delta.SelectStart == delta.SelectEnd {
				logging.LogError(nil, nil, fmt.Errorf("select delta with no selection"))
				return false
			}
			if delta.SelectStart > delta.SelectEnd ||
				delta.SelectEnd > uint(len(currentText)) {
				logging.LogError(nil, nil, fmt.Errorf("invalid select delta"))
				return false
			}

		case "cursor":
			if delta.SelectStart != delta.SelectEnd {
				return false
			}
			if delta.SelectStart > uint(len(currentText)) {
				return false
			}

		default:
			return false
		}

	}

	// Ensure final text matches
	if strings.TrimSpace(currentText) != strings.TrimSpace(finalText) {
		logging.LogError(nil, nil, fmt.Errorf("final text does not match, %s, %s", currentText, finalText))
		return false
	}

	return true

}
