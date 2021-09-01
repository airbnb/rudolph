package dynamodb

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// Per recommendation from AWS docs: https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

type mockGetObjectAPI func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)

func (m mockGetObjectAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m(ctx, params, optFns...)
}

type unmarshalType struct {
	Field1 string `dynamodbav:"field_1"`
	Field2 string `dynamodbav:"field_2"`
}

func Test_getItemWithUnmarshal_Success(t *testing.T) {

	var unmarshalTarget unmarshalType

	output, err := getItemWithUnmarshal(
		"test_table",
		mockGetObjectAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {

			expectKey := map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{
					Value: "AA",
				},
				"SK": &types.AttributeValueMemberS{
					Value: "BB",
				},
			}

			assert.Equal(t, "test_table", *params.TableName)
			assert.True(t, *params.ConsistentRead)
			assert.Equal(t, expectKey, params.Key)

			return &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"field_1": &types.AttributeValueMemberS{
						Value: "stringvalue",
					},
					"field_2": &types.AttributeValueMemberS{
						Value: "otherstringvalue",
					},
				},
			}, nil
		}),
		PrimaryKey{
			PartitionKey: "AA",
			SortKey:      "BB",
		},
		true,
		1*time.Second,
	)

	assert.Empty(t, err)

	attributevalue.UnmarshalMap(output.Item, &unmarshalTarget)
	assert.NotEmpty(t, unmarshalTarget)
	assert.Equal(t, "stringvalue", unmarshalTarget.Field1)
	assert.Equal(t, "otherstringvalue", unmarshalTarget.Field2)
}
