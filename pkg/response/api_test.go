package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIResponse_StatusCode(t *testing.T) {
	// Test status codes
	resp, _ := APIResponse(200, nil)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAPIResponse_Headers(t *testing.T) {
	// Test Content-Type header
	resp, _ := APIResponse(200, nil)
	assert.Equal(t, map[string]string{"Content-Type": "application/json"}, resp.Headers)
}

func TestAPIResponse_Body(t *testing.T) {
	// Test json marshalling
	resp, _ := APIResponse(200, map[string]string{"key1": "value1", "key2": "value2"})
	assert.Equal(t, `{"key1":"value1","key2":"value2"}`, resp.Body)
}

func TestAPIResponse_Whoops(t *testing.T) {
	// in cases when the object cannot be serialized to json
	resp, err := APIResponse(200, map[string]complex64{"key1": 1 + 4i})
	assert.Equal(t, `value could not be serialized to JSON`, resp.Body)
	assert.IsType(t, &json.UnsupportedTypeError{}, err)
}
