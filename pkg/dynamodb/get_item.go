package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GetItemAPI interface {
	GetItem(key PrimaryKey, consistentRead bool) (*dynamodb.GetItemOutput, error)
}

func (dbc concreteDynamoDBClient) GetItem(key PrimaryKey, consistentRead bool) (*dynamodb.GetItemOutput, error) {
	return getItemWithUnmarshal(dbc.tableName, &dbc.awsclient, key, consistentRead, dbc.timeout)
}

type dynamodbGetItemAPI interface {
	GetItem(ctx context.Context, in *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func getItemWithUnmarshal(tableName string, api dynamodbGetItemAPI, key PrimaryKey, consistentRead bool, timeout time.Duration) (*dynamodb.GetItemOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var keyInput map[string]types.AttributeValue
	keyInput, _ = attributevalue.MarshalMap(key)

	params := &dynamodb.GetItemInput{
		Key:            keyInput,
		TableName:      aws.String(tableName),
		ConsistentRead: aws.Bool(consistentRead),
	}

	return api.GetItem(ctx, params)
}
