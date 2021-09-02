package ruledownload

import (
	"testing"
	"time"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type mockUpdater func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error)

func (m mockUpdater) UpdateItem(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
	return m(key, item)
}

type mockGetter func(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error)

func (m mockGetter) GetItem(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
	return m(key, consistentRead)
}

var _ dynamodb.GetItemAPI = mockGetter(nil)
var _ dynamodb.UpdateItemAPI = mockUpdater(nil)

type mockTimeProvider struct{}

func (m mockTimeProvider) Now() time.Time {
	t, _ := clock.ParseRFC3339("2000-01-01T00:00:00Z")
	return t
}

func Test_ConcreteRuledownloadCursorService_ConstructCursor(t *testing.T) {
	type test struct {
		syncState      map[string]awsdynamodbtypes.AttributeValue
		expectStrategy ruledownloadStrategy
		expectPK       string
		expectSK       string
	}

	cases := []test{
		{
			syncState: map[string]awsdynamodbtypes.AttributeValue{
				"PK":        &awsdynamodbtypes.AttributeValueMemberS{Value: "Whatever"},
				"SK":        &awsdynamodbtypes.AttributeValueMemberS{Value: "Doesnt matter"},
				"CleanSync": &awsdynamodbtypes.AttributeValueMemberBOOL{Value: true},
				"BatchSize": &awsdynamodbtypes.AttributeValueMemberN{Value: "17"},
			},
			expectStrategy: ruledownloadStrategyClean,
		},
		{
			syncState: map[string]awsdynamodbtypes.AttributeValue{
				"PK":             &awsdynamodbtypes.AttributeValueMemberS{Value: "Whatever"},
				"SK":             &awsdynamodbtypes.AttributeValueMemberS{Value: "Doesnt matter"},
				"CleanSync":      &awsdynamodbtypes.AttributeValueMemberBOOL{Value: false},
				"BatchSize":      &awsdynamodbtypes.AttributeValueMemberN{Value: "17"},
				"FeedSyncCursor": &awsdynamodbtypes.AttributeValueMemberS{Value: "2000-01-01T02:00:00Z"},
			},
			expectStrategy: ruledownloadStrategyIncremental,
			expectPK:       "RulesFeed",
			expectSK:       "2000-01-01T02:00:00Z",
		},
	}

	for _, testcase := range cases {
		machineID := "AAAA-BBBB-CCCC-DDDD"
		service := concreteRuledownloadCursorService{
			timer: mockTimeProvider{},
			getter: mockGetter(
				func(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
					assert.True(t, consistentRead)
					assert.Equal(t, key.PartitionKey, "Machine#AAAA-BBBB-CCCC-DDDD")
					assert.Equal(t, key.SortKey, "SyncState")

					return &awsdynamodb.GetItemOutput{Item: testcase.syncState}, nil
				},
			),
			updater: mockUpdater(
				func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
					assert.Equal(t, key.PartitionKey, "Machine#AAAA-BBBB-CCCC-DDDD")
					assert.Equal(t, key.SortKey, "SyncState")

					// This is tested elsewhere

					return &awsdynamodb.UpdateItemOutput{}, nil
				},
			),
		}

		req := RuledownloadRequest{
			Cursor: nil,
		}
		cursor, err := service.ConstructCursor(req, machineID)

		assert.Empty(t, err)
		assert.Equal(t, testcase.expectStrategy, cursor.Strategy)
		assert.Equal(t, 17, cursor.BatchSize)
		assert.Equal(t, 1, cursor.PageNumber)
		if testcase.expectPK != "" {
			assert.Equal(t, testcase.expectPK, cursor.PartitionKey)
		} else {
			assert.Empty(t, cursor.PartitionKey)
		}
		if testcase.expectSK != "" {
			assert.Equal(t, testcase.expectSK, cursor.SortKey)
		} else {
			assert.Empty(t, cursor.SortKey)
		}
	}
}
