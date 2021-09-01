package scan

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type scanApi func(in *awsdynamodb.ScanInput) (*awsdynamodb.ScanOutput, error)

func (g scanApi) Scan(in *awsdynamodb.ScanInput) (*awsdynamodb.ScanOutput, error) {
	return g(in)
}

type testItem struct {
	dynamodb.PrimaryKey
	Field1 string `dynamodbav:"Field1"`
}

func Test_ScanService(t *testing.T) {

	scanner := scanApi(
		func(in *awsdynamodb.ScanInput) (*awsdynamodb.ScanOutput, error) {
			assert.Equal(t, int32(1), *in.Limit)

			page1Key := dynamodb.PrimaryKey{PartitionKey: "Page1", SortKey: "Yes"}
			page2Key := dynamodb.PrimaryKey{PartitionKey: "EndOfPage1", SortKey: "Yes"}
			page3Key := dynamodb.PrimaryKey{PartitionKey: "EndOfPage2", SortKey: "Yes"}
			item1 := testItem{
				PrimaryKey: page1Key,
				Field1:     "value1",
			}
			item2 := testItem{
				PrimaryKey: page2Key,
				Field1:     "value2",
			}
			item3 := testItem{
				PrimaryKey: page3Key,
				Field1:     "value3",
			}

			page2KeyIn, _ := attributevalue.MarshalMap(page2Key)
			page3KeyIn, _ := attributevalue.MarshalMap(page3Key)
			item1in, _ := attributevalue.MarshalMap(item1)
			item2in, _ := attributevalue.MarshalMap(item2)
			item3in, _ := attributevalue.MarshalMap(item3)

			if in.ExclusiveStartKey == nil || len(in.ExclusiveStartKey) == 0 {
				return &awsdynamodb.ScanOutput{
					LastEvaluatedKey: page2KeyIn,
					Items: []map[string]types.AttributeValue{
						item1in,
					},
					Count: 1,
				}, nil
			}

			var exclusiveStartKey dynamodb.PrimaryKey
			_ = attributevalue.UnmarshalMap(in.ExclusiveStartKey, &exclusiveStartKey)

			if exclusiveStartKey.PartitionKey == "EndOfPage1" {
				return &awsdynamodb.ScanOutput{
					LastEvaluatedKey: page3KeyIn,
					Items: []map[string]types.AttributeValue{
						item2in,
					},
					Count: 1,
				}, nil
			}

			return &awsdynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					item3in,
				},
				Count: 1,
			}, nil
		},
	)

	service := ConcreteScanService{
		scanner: scanner,
	}

	in := awsdynamodb.ScanInput{
		Limit: aws.Int32(int32(1)),
	}

	numItems := 0
	callback := func(out *awsdynamodb.ScanOutput) error {
		numItems += int(out.Count)
		return nil
	}

	stop := func(out *awsdynamodb.ScanOutput) (bool, error) {
		return false, nil
	}

	err := service.ScanAll(
		in,
		callback,
		stop,
	)

	assert.Empty(t, err)
	assert.Equal(t, 3, numItems)
}
