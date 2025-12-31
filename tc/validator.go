package tc

import (
	"regexp"
)

var validChars = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

func IsValidUsername(s string) bool {
	return validChars.MatchString(s)
}