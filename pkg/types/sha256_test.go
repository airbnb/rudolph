package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes_Validate_ValidateSha256_Success_LowerCase(t *testing.T) {
	err := ValidateSha256("850524a7218b2370b43f56b0e62ad6a92d8d6e1c5e9efbc578a68284d9da3545")
	assert.Empty(t, err)
}

func TestTypes_Validate_ValidateSha256_Failure_MixedCases(t *testing.T) {
	err := ValidateSha256("850524a7218b2370b43f56b0e62ad6a92d8D6e1c5e9efbc578a68284d9da3545")
	assert.NotEmpty(t, err)
}

func TestTypes_Validate_ValidateSha256_SuccessFailure(t *testing.T) {
	err := ValidateSha256("850524a7218b2370b43f56b0e62ad6a92d8d6e1c5e9efbc578a68284d9Da35466")
	assert.NotEmpty(t, err)
}
