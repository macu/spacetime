package colour

import (
	"regexp"
)

var validColourPattern = regexp.MustCompile(`^rgb\((\d{1,3}), (\d{1,3}), (\d{1,3})\)$`)

func IsValidColour(colour string) bool {
	return validColourPattern.MatchString(colour)
}
