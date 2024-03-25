package types

import (
	"errors"
	"fmt"
	"regexp"
)

var sha256Regexp = regexp.MustCompile(`^[a-f0-9]{64}$`)

func ValidateSha256(sha256 string) error {
	if !sha256Regexp.MatchString(sha256) {
		return errors.New(fmt.Sprintf("invalid sha256: %s", sha256))
	}
	return nil
}
