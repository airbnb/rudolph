package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes_ValidateMachineID_Success(t *testing.T) {
	err := ValidateMachineID("858CBF28-5EAA-58A3-A155-BB4E90C3B5DD")
	assert.Empty(t, err)
}

func TestTypes_ValidateMachineID_Failure(t *testing.T) {
	err := ValidateMachineID("858CBF28-5EAA-58A3-A155-BB4E90C3B5DDS")
	assert.NotEmpty(t, err)
}
