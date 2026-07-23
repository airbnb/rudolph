package kinesis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// An eventupload with no forwardable events must not call Kinesis PutRecords,
// which rejects an empty Records slice with "InvalidParameter: minimum field
// size of 1". Send should short-circuit to a no-op success. kinesisService is
// intentionally left nil: the guard must return before it is dereferenced, so
// this also proves the empty check runs first.
func Test_Send_NoRecords_IsNoop(t *testing.T) {
	c := client{kinesisStream: "test-stream"}

	assert.NoError(t, c.Send("machine-id", KinesisEvents{Items: nil}))
	assert.NoError(t, c.Send("machine-id", KinesisEvents{Items: []interface{}{}}))
}
