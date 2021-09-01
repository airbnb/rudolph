package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type QueryAPI interface {
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

func (dbc concreteDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return query(dbc.tableName, &dbc.awsclient, input)
}

type dynamodbQueryAPI interface {
	Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// query returns [limit] number of rows in the DDB, for the requested [partitionKey], ordered by the sortkey.
// If you provide a [cursor] it will start the query from that cursor.
// If there are additional pages to paginate over, will return a cursor. Else, the [nextCursor] is nil.
func query(tableName string, api dynamodbQueryAPI, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// expressionAttributeValues := map[string]types.AttributeValue{
	// 	":pk": &types.AttributeValueMemberS{
	// 		Value: partitionKey,
	// 	},
	// }
	// keyConditionExpression := aws.String("PK = :pk")

	// var exclusiveStartKey map[string]types.AttributeValue
	// input := &dynamodb.QueryInput{
	// 	TableName:                 aws.String(tableName),
	// 	ConsistentRead:            aws.Bool(consistentRead),
	// 	ExpressionAttributeValues: expressionAttributeValues,
	// 	KeyConditionExpression:    keyConditionExpression,
	// 	ExclusiveStartKey:         exclusiveStartKey,
	// 	FilterExpression: ,
	// 	Limit:                     limit,
	// }

	// // var exclusiveStartKey map[string]types.AttributeValue{}
	// exclusiveStartKey, _ = attributevalue.MarshalMap(cursor)

	input.TableName = aws.String(tableName)

	return api.Query(ctx, input)
}
