package machinerules

import (
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
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

	// log.Printf("DDB Query Input:\n%+v", input)

	output, err := client.Query(&input)

	// log.Printf("Error:\n%+v", err)
	// log.Printf("DDB Query Output:\n%+v", output)
	// log.Printf("Discovered %d items", len(output.Items))

	if err != nil {
		return
	}
	err = attributevalue.UnmarshalListOfMaps(output.Items, &keys)

	// log.Printf("Keys:\n%+v", *keys)

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

	// log.Printf("Executing DynamoDB Query:\n%+v", input)

	result, err := client.Query(input)
	if err != nil {
		err = errors.Wrapf(err, "failed to read rules from DynamoDB for partitionKey %q", partitionKey)
		return
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal result from DynamoDB")
		return
	}
	// log.Printf("    got %d items from query.", len(*items))
	return
}
