package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UpdateItemAPI interface {
	UpdateItem(key PrimaryKey, item interface{}) (*dynamodb.UpdateItemOutput, error)
}

func (dbc concreteDynamoDBClient) UpdateItem(key PrimaryKey, item interface{}) (*dynamodb.UpdateItemOutput, error) {
	return updateItemToDynamoDB(dbc.tableName, &dbc.awsclient, key, item, dbc.timeout)
}

type dynamodbUpdateItemAPI interface {
	UpdateItem(ctx context.Context, in *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

func updateItemToDynamoDB(tableName string, api dynamodbUpdateItemAPI, updateKey PrimaryKey, updateFields interface{}, timeout time.Duration) (*dynamodb.UpdateItemOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	key, err := attributevalue.MarshalMap(updateKey)
	if err != nil {
		return nil, err
	}

	// Conditional expression that requires the PK must match the PartitionKey provided
	conditions := expression.Equal(expression.Name("PK"), expression.Value(updateKey.PartitionKey))

	updateItems, err := attributevalue.MarshalMap(updateFields)
	if err != nil {
		return nil, err
	}

	var updates expression.UpdateBuilder
	for key, value := range updateItems {
		updates = updates.Set(expression.Name(key), expression.Value(value))
	}

	expr, err := expression.NewBuilder().
		WithUpdate(updates).
		WithCondition(conditions).
		Build()

	if err != nil {
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(tableName),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueAllNew,
	}

	return api.UpdateItem(ctx, input)
}
