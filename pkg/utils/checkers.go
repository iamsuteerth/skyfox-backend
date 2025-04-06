package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	maxEmailLength = 254
)

func IsValidEmail(email string) bool {
	if len(email) > maxEmailLength || len(email) < 3 {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domain := parts[1]

	if len(localPart) == 0 || len(domain) == 0 {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}
	domainParts := strings.Split(domain, ".")
	tld := domainParts[len(domainParts)-1]
	return len(tld) >= 2
}
