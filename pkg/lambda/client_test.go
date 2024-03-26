package lambda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	client := GetClient(
		"rudolph-test",
		"dev",
		"us-east-1",
	)

	assert.NotNil(t, client)
}
