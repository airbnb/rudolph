package dynamodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type mockQueryAPI func(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)

func (m mockQueryAPI) Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return m(ctx, in, optFns...)
}

// Per recommendation from AWS docs: https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/
func Test_query(t *testing.T) {
	output, err := query(
		"test_table",
		mockQueryAPI(func(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			assert.Equal(t, "test_table", *in.TableName)
			assert.Equal(t, "PK = :PK AND SK = :SK", aws.ToString(in.KeyConditionExpression))
			assert.Equal(t, aws.Int32(5), in.Limit)
			assert.Equal(
				t,
				map[string]string{
					"#PK":     "PK",
					"#SK":     "SK",
					"#Number": "Number",
				},
				in.ExpressionAttributeNames,
			)
			assert.Equal(
				t,
				map[string]types.AttributeValue{},
				in.ExpressionAttributeValues,
			)

			countItems := 3
			items := make([]map[string]types.AttributeValue, countItems)

			for i := range items {
				items[i] = map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{
						Value: fmt.Sprintf("PK#AA#%d", i),
					},
					"SK": &types.AttributeValueMemberS{
						Value: fmt.Sprintf("SK#BB#%d", i),
					},
					"Number": &types.AttributeValueMemberN{
						Value: fmt.Sprintf("%d", i),
					},
				}
			}

			return &dynamodb.QueryOutput{
				Items:            items,
				Count:            int32(countItems),
				LastEvaluatedKey: map[string]types.AttributeValue{},
			}, nil
		}),
		&dynamodb.QueryInput{
			TableName:              aws.String("test_table"),
			KeyConditionExpression: aws.String("PK = :PK AND SK = :SK"),
			ExpressionAttributeNames: map[string]string{
				"#PK":     "PK",
				"#SK":     "SK",
				"#Number": "Number",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{},
			Limit:                     aws.Int32(5),
		},
	)

	if err != nil {
		t.Errorf("Error was not expected: %v", err)
	}

	if output == nil {
		t.Errorf("Output was not expected: %v", output)
	}

	var items []map[string]interface{}

	err = attributevalue.UnmarshalListOfMaps(output.Items, &items)
	if err != nil {
		t.Errorf("failed to unmarshal result from DynamoDB: %s", err.Error())
	}

	assert.NotEmpty(t, items)
	assert.Equal(t, 3, len(items))
	for i, item := range items {
		assert.Equal(t, fmt.Sprintf("PK#AA#%d", i), item["PK"])
		assert.Equal(t, fmt.Sprintf("SK#BB#%d", i), item["SK"])
		assert.Equal(t, float64(i), item["Number"])
	}
}
