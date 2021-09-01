package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DeleteItemAPI interface {
	DeleteItem(key PrimaryKey) (*dynamodb.DeleteItemOutput, error)
}

func (dbc concreteDynamoDBClient) DeleteItem(key PrimaryKey) (*dynamodb.DeleteItemOutput, error) {
	return deleteItemFromDynamoDB(dbc.tableName, &dbc.awsclient, key, dbc.timeout)
}

type dynamodbDeleteItemAPI interface {
	DeleteItem(ctx context.Context, in *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

func deleteItemFromDynamoDB(tableName string, api dynamodbDeleteItemAPI, key PrimaryKey, timeout time.Duration) (*dynamodb.DeleteItemOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var keyInput map[string]types.AttributeValue
	keyInput, _ = attributevalue.MarshalMap(key)

	params := &dynamodb.DeleteItemInput{
		Key:       keyInput,
		TableName: aws.String(tableName),
	}

	return api.DeleteItem(ctx, params)
}
