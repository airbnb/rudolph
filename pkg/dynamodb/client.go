package dynamodb

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// https://aws.amazon.com/blogs/developer/aws-sdk-for-go-version-2-general-availability/

type DynamoDBClient interface {
	DeleteItemAPI
	GetItemAPI
	PutItemAPI
	UpdateItemAPI
	QueryAPI
	TransactWriteItemsAPI
	ScanAPI
}

type concreteDynamoDBClient struct {
	awsclient dynamodb.Client
	tableName string
	timeout   time.Duration
}

func GetClient(inputTableName string, region string) DynamoDBClient {
	return GetClientWithTimeout(inputTableName, region, defaultTimeout)
}

func GetClientWithTimeout(inputTableName string, region string, timeout time.Duration) DynamoDBClient {
	// Create Amazon S3 API client using path style addressing.
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return concreteDynamoDBClient{
		awsclient: *client,
		tableName: inputTableName,
		timeout:   defaultTimeout,
	}
}

const (
	defaultTimeoutMS = 5000
	defaultTimeout   = defaultTimeoutMS * time.Millisecond
)
