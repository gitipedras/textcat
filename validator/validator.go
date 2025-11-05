package validator

import(
	"regexp"
	"strings"
	"textcat/models"
)

var UsernameRegex = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// validates a username
// if returns false then is invalid
func Username(username string) bool {
	if len(username) < 0 {
		return false
	}

	if username == "" {
		return false
	}
	return UsernameRegex.MatchString(username)
}

func Message(message string) bool {
    // check empty or only spaces
    if strings.TrimSpace(message) == "" {
        return false
    }

    // check length
    if len(message) > models.Config.MaxLength {
        return false
    }

    return true
}
