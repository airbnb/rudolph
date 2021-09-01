package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PutItemAPI interface {
	PutItem(item interface{}) (*dynamodb.PutItemOutput, error)
}

func (dbc concreteDynamoDBClient) PutItem(item interface{}) (*dynamodb.PutItemOutput, error) {
	return putItem(dbc.tableName, &dbc.awsclient, item, dbc.timeout)
}

type dynamodbPutItemAPI interface {
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func putItem(tableName string, api dynamodbPutItemAPI, item interface{}, timeout time.Duration) (*dynamodb.PutItemOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	putItemItem, err := attributevalue.MarshalMap(item)

	if err != nil {
		return nil, err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      putItemItem,
	}

	return api.PutItem(ctx, params)
}
