package types

import "regexp"

var emailPattern *regexp.Regexp

func ValidateEmailAddress(email string) bool {
	if emailPattern == nil {
		emailPattern = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	}
	return emailPattern.MatchString(email)
}
