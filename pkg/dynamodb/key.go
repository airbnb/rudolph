package dynamodb

type PrimaryKey struct {
	PartitionKey string `dynamodbav:"PK"`
	SortKey      string `dynamodbav:"SK"`
}
