package spacetime

type NakedTextDelta struct {
	Timestmap uint `json:"timestamp"`

	// keystrokes
	Key *uint `json:"key,omitempty"`

	// cursor positioning
	Cursor *uint `json:"cursor,omitempty"`

	// selections
	SelectStart *uint `json:"select_start,omitempty"`
	SelectEnd   *uint `json:"select_end,omitempty"`

	// replacements/paste
	ReplaceStart *uint   `json:"replace_start,omitempty"`
	ReplaceEnd   *uint   `json:"replace_end,omitempty"`
	ReplaceText  *string `json:"replace_text,omitempty"`
}

type NakedText []NakedTextDelta

func ValidateNakedText(text NakedText) bool {

	// Ensure has count
	if len(text) == 0 {
		return false
	}

	// Ensure first delta at timestamp 0
	if text[0].Timestmap != 0 {
		return false
	}

	// Ensure timestmaps increment
	for i := 1; i < len(text); i++ {
		if text[i].Timestmap <= text[i-1].Timestmap {
			return false
		}
	}

	// Ensure full data is available for each type of delta
	for _, delta := range text {
		if delta.Key != nil {
			if delta.Cursor != nil || delta.SelectStart != nil || delta.SelectEnd != nil ||
				delta.ReplaceStart != nil || delta.ReplaceEnd != nil || delta.ReplaceText != nil {
				return false
			}
		}

		if delta.Cursor != nil {
			if delta.Key != nil || delta.SelectStart != nil || delta.SelectEnd != nil ||
				delta.ReplaceStart != nil || delta.ReplaceEnd != nil || delta.ReplaceText != nil {
				return false
			}
		}

		if delta.SelectStart != nil || delta.SelectEnd != nil {
			if delta.Key != nil || delta.Cursor != nil ||
				delta.SelectStart == nil || delta.SelectEnd == nil ||
				delta.ReplaceStart != nil || delta.ReplaceEnd != nil || delta.ReplaceText != nil {
				return false
			}
		}

		if delta.ReplaceStart != nil || delta.ReplaceEnd != nil || delta.ReplaceText != nil {
			if delta.Key != nil || delta.Cursor != nil ||
				delta.SelectStart != nil || delta.SelectEnd != nil ||
				delta.ReplaceStart == nil || delta.ReplaceEnd == nil || delta.ReplaceText == nil {
				return false
			}
		}
	}

	return true

}
