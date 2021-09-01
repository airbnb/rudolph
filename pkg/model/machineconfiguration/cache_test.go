package machineconfiguration

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/stretchr/testify/assert"
)

func Test_Cache(t *testing.T) {
	c := GetCache(clock.Y2K{})

	config := MachineConfiguration{
		CleanSync: true,
	}

	assert.False(t, c.Has("a"))
	assert.False(t, c.Has("b"))
	assert.False(t, c.Has("c"))

	c.Set("a", &config)
	c.Set("b", nil)

	assert.True(t, c.Has("a"))
	assert.True(t, c.Has("b"))
	assert.False(t, c.Has("c"))

	fetchedConfig := c.Get("a")
	assert.True(t, fetchedConfig.CleanSync)

	otherFetched := c.Get("b")
	assert.Nil(t, otherFetched)
}

func Test_CacheExpiration(t *testing.T) {
	timeMachine := clock.TimeMachine{
		Current: clock.Y2KTime(),
	}

	// Has 1 hour cache expiration hardcoded into this
	c := GetCache(&timeMachine)

	assert.False(t, c.Has("a"))

	c.Set("a", nil)

	assert.True(t, c.Has("a"))

	newTime, _ := clock.ParseRFC3339("2000-01-01T01:00:00Z")
	timeMachine.Travel(newTime)

	assert.False(t, c.Has("a"))
}
