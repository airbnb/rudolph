package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Unixtimestamp(t *testing.T) {
	then, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	assert.Equal(t, int64(946684800), Unixtimestamp(then))

	than := FromUnixtimestamp(Unixtimestamp(then))
	assert.Equal(t, int64(946684800), Unixtimestamp(than))
}

func Test_RFC3339(t *testing.T) {
	then, err := ParseRFC3339("2000-01-01T00:00:00-07:00")
	assert.Empty(t, err)

	assert.Equal(t, "2000-01-01T07:00:00Z", then.UTC().Format(time.RFC3339))
	assert.Equal(t, "2000-01-01T07:00:00Z", RFC3339(then.UTC()))
}
