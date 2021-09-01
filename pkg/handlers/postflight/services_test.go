package postflight

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type mockUpdater func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error)

func (m mockUpdater) UpdateItem(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
	return m(key, item)
}

var _ dynamodb.UpdateItemAPI = mockUpdater(nil)

func Test_ConcreteSyncStateUpdater_OK(t *testing.T) {
	cur, _ := clock.ParseRFC3339("2000-01-01T00:00:00Z")
	updater := concreteSyncStateUpdater{
		updater: mockUpdater(
			func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
				assert.Equal(t, "Machine#AAAA-BBBB-CCCC", key.PartitionKey)
				assert.Equal(t, "SyncState", key.SortKey)

				thisItem := item.(syncstate.UpdatePostflightItem)
				assert.Equal(t, "2000-01-01T00:00:00Z", thisItem.PostflightAt)
				assert.Equal(t, "1999-12-31T23:50:00Z", thisItem.FeedSyncCursor) // FIXME (derek.wang) This isn't necessarily "good" nor "correct", just "lazy"
				assert.Equal(t, cur.AddDate(0, 0, 90).Unix(), thisItem.ExpiresAfter)

				return &awsdynamodb.UpdateItemOutput{}, nil
			},
		),
		timeProvider: clock.FrozenTimeProvider{Current: cur},
	}

	err := updater.updatePostflightDate("AAAA-BBBB-CCCC")
	assert.Empty(t, err)
}

type mockQueryer func(input *awsdynamodb.QueryInput) (*awsdynamodb.QueryOutput, error)
type mockDeleter func(key dynamodb.PrimaryKey) (*awsdynamodb.DeleteItemOutput, error)

func (m mockQueryer) Query(input *awsdynamodb.QueryInput) (*awsdynamodb.QueryOutput, error) {
	return m(input)
}
func (m mockDeleter) DeleteItem(key dynamodb.PrimaryKey) (*awsdynamodb.DeleteItemOutput, error) {
	return m(key)
}

var _ dynamodb.QueryAPI = mockQueryer(nil)
var _ dynamodb.DeleteItemAPI = mockDeleter(nil)

func Test_RuleDeleter_OK(t *testing.T) {
	destroyer := concreteRuleDestroyer{
		queryer: mockQueryer(
			func(input *awsdynamodb.QueryInput) (*awsdynamodb.QueryOutput, error) {
				assert.False(t, *input.ConsistentRead)
				assert.Equal(t, "PK = :pk", *input.KeyConditionExpression)
				assert.Equal(t, "DeleteOnNextSync = :boo", *input.FilterExpression)
				assert.Equal(t, "PK, SK", *input.ProjectionExpression)

				pk := input.ExpressionAttributeValues[":pk"].(*types.AttributeValueMemberS)
				assert.NotEmpty(t, pk)
				assert.Equal(t, "MachineRules#AAAA-BBBB-CCCC", pk.Value)

				boo := input.ExpressionAttributeValues[":boo"].(*types.AttributeValueMemberBOOL)
				assert.NotEmpty(t, boo)
				assert.True(t, boo.Value)

				return &awsdynamodb.QueryOutput{}, nil
			},
		),
		deleter: mockDeleter(
			func(key dynamodb.PrimaryKey) (*awsdynamodb.DeleteItemOutput, error) {
				assert.Equal(t, "Machine#AAAA-BBBB-CCCC", key.PartitionKey)
				assert.Equal(t, "SyncState", key.SortKey)

				return &awsdynamodb.DeleteItemOutput{}, nil
			},
		),
	}

	err := destroyer.destroyMachineRulesMarkedForDeletion("AAAA-BBBB-CCCC")
	assert.Empty(t, err)
}
