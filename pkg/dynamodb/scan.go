package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ScanAPI interface {
	Scan(in *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

func (dbc concreteDynamoDBClient) Scan(in *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return scan(dbc.tableName, &dbc.awsclient, in, dbc.timeout)
}

type dynamodbScanAPI interface {
	Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func scan(tableName string, api dynamodbScanAPI, in *dynamodb.ScanInput, timeout time.Duration) (*dynamodb.ScanOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	in.TableName = &tableName

	return api.Scan(ctx, in)
}
