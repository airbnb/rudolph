package machinerules

import (
	"fmt"
	"testing"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type mockGetter func(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error)

func (m mockGetter) GetItem(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
	return m(key, consistentRead)
}

type mockUpdater func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error)

func (m mockUpdater) UpdateItem(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
	return m(key, item)
}

var _ dynamodb.GetItemAPI = mockGetter(nil)
var _ dynamodb.UpdateItemAPI = mockUpdater(nil)

func Test_RemoveMachineRule_OK(t *testing.T) {
	type test struct {
		globalRuleExists    bool
		expectedFinalPolicy types.Policy
	}

	cases := []test{
		{
			globalRuleExists:    false,
			expectedFinalPolicy: types.Remove,
		},
		{
			globalRuleExists:    true,
			expectedFinalPolicy: types.Blocklist,
		},
	}

	for _, testcase := range cases {
		machineID := "AAAA-BBBB-CCCC"
		var updatecalled bool
		err := RemoveMachineRule(
			mockGetter(
				func(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
					switch key.PartitionKey {
					case "MachineRules#AAAA-BBBB-CCCC":
						if key.SortKey == "AAA#SORTKEY" {
							return &awsdynamodb.GetItemOutput{
								Item: map[string]awsdynamodbtypes.AttributeValue{
									"PK":     &awsdynamodbtypes.AttributeValueMemberS{Value: key.PartitionKey},
									"SK":     &awsdynamodbtypes.AttributeValueMemberS{Value: key.SortKey},
									"Policy": &awsdynamodbtypes.AttributeValueMemberN{Value: "1"},
								},
							}, nil
						}

					case "GlobalRules":
						if key.SortKey == "AAA#SORTKEY" {
							if !testcase.globalRuleExists {
								return &awsdynamodb.GetItemOutput{}, nil
							}
							return &awsdynamodb.GetItemOutput{
								Item: map[string]awsdynamodbtypes.AttributeValue{
									"PK":     &awsdynamodbtypes.AttributeValueMemberS{Value: key.PartitionKey},
									"SK":     &awsdynamodbtypes.AttributeValueMemberS{Value: key.SortKey},
									"Policy": &awsdynamodbtypes.AttributeValueMemberN{Value: "2"},
								},
							}, nil
						}
					}
					assert.Fail(t, fmt.Sprintf("dynamodb:GetItem call unexpected with %+v", key))
					return &awsdynamodb.GetItemOutput{}, nil
				},
			),
			mockUpdater(
				func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
					assert.Equal(t, "MachineRules#AAAA-BBBB-CCCC", key.PartitionKey)
					assert.Equal(t, "AAA#SORTKEY", key.SortKey)

					testitem := item.(ruleRemovalRequest)
					assert.Equal(t, testcase.expectedFinalPolicy, testitem.Policy)
					assert.True(t, testitem.DeleteOnNextSync)

					updatecalled = true

					return &awsdynamodb.UpdateItemOutput{}, nil
				},
			),
			machineID,
			"AAA#SORTKEY",
		)

		assert.True(t, updatecalled)
		assert.Empty(t, err)
	}
}
