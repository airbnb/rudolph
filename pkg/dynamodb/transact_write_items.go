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

type TransactWriteItemsAPI interface {
	TransactWriteItems(items []types.TransactWriteItem, idempotencyToken *string) (*dynamodb.TransactWriteItemsOutput, error)
	CreateTransactPutItem(item interface{}) (*types.TransactWriteItem, error)
	CreateTransactUpdateItem(key PrimaryKey, item interface{}) (*types.TransactWriteItem, error)
	CreateTransactDeleteItem(key PrimaryKey) (*types.TransactWriteItem, error)
}

func (dbc concreteDynamoDBClient) TransactWriteItems(items []types.TransactWriteItem, idempotencyToken *string) (*dynamodb.TransactWriteItemsOutput, error) {
	return transactWriteItems(&dbc.awsclient, items, idempotencyToken, dbc.timeout)
}

type dynamodbTransactWriteItemsAPI interface {
	TransactWriteItems(ctx context.Context, in *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

func transactWriteItems(api dynamodbTransactWriteItemsAPI, items []types.TransactWriteItem, idempotencyToken *string, timeout time.Duration) (*dynamodb.TransactWriteItemsOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	}

	if idempotencyToken != nil {
		input.ClientRequestToken = idempotencyToken
	}

	return api.TransactWriteItems(ctx, input)
}

func (dbc concreteDynamoDBClient) CreateTransactPutItem(item interface{}) (writeItem *types.TransactWriteItem, err error) {
	putItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	writeItem = &types.TransactWriteItem{
		Put: &types.Put{
			Item:      putItem,
			TableName: aws.String(dbc.tableName),
		},
	}
	return
}

func (dbc concreteDynamoDBClient) CreateTransactUpdateItem(primaryKey PrimaryKey, updateFields interface{}) (writeItem *types.TransactWriteItem, err error) {
	key, err := attributevalue.MarshalMap(primaryKey)
	if err != nil {
		return nil, err
	}

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
		Build()
	if err != nil {
		return nil, err
	}

	writeItem = &types.TransactWriteItem{
		Update: &types.Update{
			Key:                       key,
			UpdateExpression:          expr.Update(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			TableName:                 aws.String(dbc.tableName),
		},
	}
	return
}

func (dbc concreteDynamoDBClient) CreateTransactDeleteItem(key PrimaryKey) (writeItem *types.TransactWriteItem, err error) {
	keyInput, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, err
	}

	writeItem = &types.TransactWriteItem{
		Delete: &types.Delete{
			Key:       keyInput,
			TableName: aws.String(dbc.tableName),
		},
	}
	return
}
