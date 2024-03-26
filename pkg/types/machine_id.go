package types

import (
	"fmt"
	"regexp"
)

var machineRegexp = regexp.MustCompile(`^[A-F0-9]{8}\-[A-F0-9]{4}\-[A-F0-9]{4}\-[A-F0-9]{4}\-[A-F0-9]{12}$`)

// ValidateMachineID returns an error if a machineID is not a properly formatted UUID string
func ValidateMachineID(machineID string) error {
	if !machineRegexp.MatchString(machineID) {
		return fmt.Errorf("invalid machineID: %s", machineID)
	}
	return nil
}
