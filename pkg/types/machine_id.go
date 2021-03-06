package types

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var machineRegexp = regexp.MustCompile(`^[A-F0-9]{8}\-[A-F0-9]{4}\-[A-F0-9]{4}\-[A-F0-9]{4}\-[A-F0-9]{12}$`)

// ValidateMachineID returns an error if a machineID is not a properly formatted UUID string
func ValidateMachineID(machineID string) error {
	if !machineRegexp.MatchString(machineID) {
		return errors.New(fmt.Sprintf("invalid machineID: %s", machineID))
	}
	return nil
}
