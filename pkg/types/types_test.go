package types

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

var (
	av dynamodb.AttributeValue
)

func TestTypesClientMode_Marshal(t *testing.T) {
	t.Run("Monitor", func(t *testing.T) {
		clientMode, err := ClientMode(1).MarshalText()
		assert.Empty(t, err)
		assert.Equal(t, clientMode, []byte("MONITOR"))
	})

	t.Run("Lockdown", func(t *testing.T) {
		clientMode, err := ClientMode(2).MarshalText()
		assert.Empty(t, err)
		assert.Equal(t, clientMode, []byte("LOCKDOWN"))
	})

	t.Run("Error", func(t *testing.T) {
		_, err := ClientMode(3).MarshalText()
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, "unknown client_mode 3")
	})
}

func TestTypesClientMode_Unmarshal(t *testing.T) {
	t.Run("Monitor", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("MONITOR"))
		assert.Empty(t, err)
		assert.Equal(t, clientMode, ClientMode(1))
	})

	t.Run("Lockdown", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("LOCKDOWN"))
		assert.Empty(t, err)
		assert.Equal(t, clientMode, ClientMode(2))
	})

	t.Run("ERROR", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("NOPE"))
		assert.NotEmpty(t, err)
	})

	t.Run("ERROR MONITORS", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("MONITORS"))
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, `unknown client_mode value "MONITORS"`)
	})
}
