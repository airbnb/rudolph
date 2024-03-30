package globalrules

import (
	"errors"
	"fmt"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PingDatabase(client dynamodb.QueryAPI) (err error) {
	_, _, err = GetPaginatedGlobalRules(client, 1, nil)
	return
}

func GetPaginatedGlobalRules(
	client dynamodb.QueryAPI,
	limit int,
	exclusiveStartKey *dynamodb.PrimaryKey,
) (
	items *[]GlobalRuleRow,
	lastEvaluatedKey *dynamodb.PrimaryKey,
	err error,
) {
	partitionKey := globalRulesPK

	if limit <= 0 {
		err = errors.New("invalid limit/batchsize specified")
		return
	}

	keyConditionExpression := aws.String("PK = :pk")
	expressionAttributeValues := map[string]types.AttributeValue{
		":pk": &types.AttributeValueMemberS{Value: partitionKey},
	}
	var exclusiveStartKeyInput map[string]types.AttributeValue
	if exclusiveStartKey != nil {
		exclusiveStartKeyInput, err = attributevalue.MarshalMap(exclusiveStartKey)
		if err != nil {
			err = fmt.Errorf("failed to marshall exclusiveStartKey: %w", err)
			return
		}
	}

	input := &awsdynamodb.QueryInput{
		ConsistentRead:            aws.Bool(false),
		ExpressionAttributeValues: expressionAttributeValues,
		KeyConditionExpression:    keyConditionExpression,
		ExclusiveStartKey:         exclusiveStartKeyInput,
		Limit:                     aws.Int32(int32(limit)),
	}

	// log.Printf("Executing DynamoDB Query:\n%+v", input)

	result, err := client.Query(input)
	if err != nil {
		err = fmt.Errorf("failed to read rules from DynamoDB for partitionKey %q: %w", partitionKey, err)
		return
	}

	if result.LastEvaluatedKey != nil {
		err = attributevalue.UnmarshalMap(result.LastEvaluatedKey, &lastEvaluatedKey)
		if err != nil {
			err = fmt.Errorf("failed to unmarshall LastEvaluatedKey: %w", err)
			return
		}
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal result from DynamoDB: %w", err)
		return
	}

	// To support legacy SHA256 types, we must transform the datasets before returning
	for _, item := range *items {
		if item.SHA256 != "" && item.Identifier == "" {
			item.Identifier = item.SHA256
		}
	}
	return
}
