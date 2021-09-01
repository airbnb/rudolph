package dynamodb

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

// Per recommendation from AWS docs: https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

type mockUpdateItemApi func(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)

func (m mockUpdateItemApi) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return m(ctx, params, optFns...)
}

type updateTestType struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	DataState string `dynamodbav:"DataState"`
}

func Test_UpdateItem(t *testing.T) {
	updatedItem := &updateTestType{
		DataState: "Updated",
	}

	dbbPKSK := PrimaryKey{
		PartitionKey: "AA",
		SortKey:      "BB",
	}

	output, err := updateItemToDynamoDB(
		"test_table",
		mockUpdateItemApi(func(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			expectKey, err := attributevalue.MarshalMap(dbbPKSK)
			if err != nil {
				return nil, err
			}

			assert.Equal(t, "test_table", *params.TableName)
			assert.Equal(t, expectKey, params.Key)

			return &dynamodb.UpdateItemOutput{}, nil
		}),
		PrimaryKey{
			PartitionKey: "AA",
			SortKey:      "BB",
		},
		updatedItem,
		1*time.Second,
	)

	assert.Empty(t, err)
	assert.Empty(t, output)
}

func Test_UpdateItem_Error(t *testing.T) {
	dbbPKSK := PrimaryKey{
		PartitionKey: "AA",
		SortKey:      "BB",
	}

	updatedItem := &updateTestType{
		DataState: "Updates",
	}

	output, err := updateItemToDynamoDB(
		"test_table",
		mockUpdateItemApi(func(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			expectKey, err := attributevalue.MarshalMap(dbbPKSK)
			if err != nil {
				return nil, err
			}

			assert.Equal(t, "test_table", *params.TableName)
			assert.Equal(t, expectKey, params.Key)

			return &dynamodb.UpdateItemOutput{}, nil
		}),
		PrimaryKey{
			PartitionKey: "AA",
			SortKey:      "BB",
		},
		updatedItem,
		1*time.Second,
	)

	assert.Empty(t, err)
	assert.Empty(t, output)
}
