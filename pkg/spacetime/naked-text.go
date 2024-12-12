package spacetime

type NakedTextDelta struct {
	Timestmap uint `json:"ts"`

	// key presses (like backspace)
	Key *uint `json:"k,omitempty"`

	// added text (one char at a time)
	AddText *string `json:"t,omitempty"`

	// cursor positioning
	Cursor *uint `json:"c,omitempty"`

	// selections
	SelectStart *uint `json:"ss,omitempty"`
	SelectEnd   *uint `json:"se,omitempty"`

	// replacements/paste (used with text/key)
	ReplaceStart *uint `json:"rs,omitempty"`
	ReplaceEnd   *uint `json:"re,omitempty"`
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
		if delta.AddText != nil {
			// Check valid text

			// - valid text will be single characters
			if len(*delta.AddText) != 1 {
				return false
			}

			// valid text will be printable (tabs allowed)
			if (*delta.AddText)[0] < 32 && (*delta.AddText)[0] != 9 {
				return false
			}

			if delta.Key != nil || delta.Cursor != nil ||
				delta.SelectStart != nil || delta.SelectEnd != nil {
				return false
			}
		}

		if delta.Key != nil {
			// Check valid key

			// - valid keys will be newlines and backspaces
			if *delta.Key != 8 && *delta.Key != 13 {
				return false
			}

			if delta.Cursor != nil || delta.SelectStart != nil || delta.SelectEnd != nil ||
				delta.ReplaceStart != nil || delta.ReplaceEnd != nil ||
				delta.AddText != nil {
				return false
			}
		}

		if delta.Cursor != nil {
			if delta.Key != nil || delta.SelectStart != nil || delta.SelectEnd != nil ||
				delta.ReplaceStart != nil || delta.ReplaceEnd != nil ||
				delta.AddText != nil {
				return false
			}
		}

		if delta.SelectStart != nil || delta.SelectEnd != nil {
			if delta.Key != nil || delta.Cursor != nil ||
				delta.SelectStart == nil || delta.SelectEnd == nil ||
				delta.ReplaceStart != nil || delta.ReplaceEnd != nil ||
				delta.AddText != nil {
				return false
			}
		}

		if delta.ReplaceStart != nil || delta.ReplaceEnd != nil {
			if delta.Cursor != nil ||
				delta.SelectStart != nil || delta.SelectEnd != nil ||
				delta.ReplaceStart == nil || delta.ReplaceEnd == nil ||
				(delta.AddText == nil && delta.Key == nil) {
				return false
			}
		}
	}

	return true

}
