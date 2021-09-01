package globalrules

import (
	"log"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

func PingDatabase(client dynamodb.QueryAPI) (err error) {
	_, _, err = GetPaginatedGlobalRules(client, 1, nil)
	return
}

func GetPaginatedGlobalRules(client dynamodb.QueryAPI, limit int, exclusiveStartKey *dynamodb.PrimaryKey) (items *[]GlobalRuleRow, lastEvaluatedKey *dynamodb.PrimaryKey, err error) {
	partitionKey := globalRulesPK

	if limit <= 0 {
		err = errors.New("Invalid limit/batchsize specified")
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
			err = errors.Wrap(err, "failed to marshall exclusiveStartKey")
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
		err = errors.Wrapf(err, "failed to read rules from DynamoDB for partitionKey %q", partitionKey)
		return
	}

	if result.LastEvaluatedKey != nil {
		err = attributevalue.UnmarshalMap(result.LastEvaluatedKey, &lastEvaluatedKey)
		if err != nil {
			err = errors.Wrap(err, "failed to unmarshall LastEvaluatedKey")
			return
		}
		log.Printf("    lastEvaluatedKey: %+v", lastEvaluatedKey)
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal result from DynamoDB")
		return
	}
	log.Printf("    got %d items from query.", len(*items))
	return
}
