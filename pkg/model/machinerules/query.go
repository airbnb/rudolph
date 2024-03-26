package machinerules

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// @deprecated
func GetPrimaryKeysByMachineIDWhereMarkedForDeletion(client dynamodb.QueryAPI, machineID string) (keys *[]dynamodb.PrimaryKey, err error) {
	pk := machineRulePK(machineID)

	input := awsdynamodb.QueryInput{
		ConsistentRead:         aws.Bool(false),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
			":boo": &types.AttributeValueMemberBOOL{
				Value: true,
			},
		},
		FilterExpression:     aws.String("DeleteOnNextSync = :boo"),
		ProjectionExpression: aws.String("PK, SK"),
	}

	output, err := client.Query(&input)

	if err != nil {
		return
	}
	err = attributevalue.UnmarshalListOfMaps(output.Items, &keys)

	return
}

// @deprecated
func GetMachineRules(client dynamodb.QueryAPI, machineID string) (items *[]MachineRuleRow, err error) {
	partitionKey := machineRulePK(machineID)

	keyConditionExpression := aws.String("PK = :pk")
	expressionAttributeValues := map[string]types.AttributeValue{
		":pk": &types.AttributeValueMemberS{Value: partitionKey},
	}

	input := &awsdynamodb.QueryInput{
		ConsistentRead:            aws.Bool(false),
		ExpressionAttributeValues: expressionAttributeValues,
		KeyConditionExpression:    keyConditionExpression,
	}

	result, err := client.Query(input)
	if err != nil {
		err = fmt.Errorf("failed to read rules from DynamoDB for partitionKey %q: %w", partitionKey, err)
		return
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal result from DynamoDB: %w", err)
		return
	}
	return
}
