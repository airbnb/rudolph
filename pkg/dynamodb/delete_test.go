package dynamodb

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// Per recommendation from AWS docs: https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

type mockDeleteObjectAPI func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)

func (m mockDeleteObjectAPI) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m(ctx, params, optFns...)
}

func Test_deleteItemFromDynamoDB_Success(t *testing.T) {

	_, err := deleteItemFromDynamoDB(
		"test_table_2",
		mockDeleteObjectAPI(func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			expectKey := map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{
					Value: "CCC",
				},
				"SK": &types.AttributeValueMemberS{
					Value: "DDD",
				},
			}

			assert.Equal(t, "test_table_2", *params.TableName)
			assert.Equal(t, expectKey, params.Key)

			return &dynamodb.DeleteItemOutput{}, nil
		}),
		PrimaryKey{
			PartitionKey: "CCC",
			SortKey:      "DDD",
		},
		1*time.Second,
	)

	assert.Empty(t, err)
}
