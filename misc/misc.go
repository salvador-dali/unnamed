package misc

import (
	"net/mail"
	"strings"
)

const (
	passwordMinLen = 8
)

func IsPasswordValid(str string) bool {
	return len(str) >= passwordMinLen
}

func IsIdValid(id int) bool {
	return id > 0
}

func ValidateString(str string, maxLen int) (string, bool) {
	str = strings.TrimSpace(str)
	if len(str) == 0 || len(str) > maxLen {
		return "", false
	}

	return str, true
}

func ValidateEmail(str string) (string, bool) {
	e, err := mail.ParseAddress(strings.ToLower(strings.TrimSpace(str)))
	if err != nil {
		return "", false
	}

	return e.Address, true
}
