package syncstate

import (
	"fmt"
	"testing"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type getSyncState func(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error)

func (get getSyncState) GetItem(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
	return get(key, consistentRead)
}

func Test_GetIntendedConfig(t *testing.T) {
	machineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	pkSyncState := fmt.Sprintf("%s%s", "Machine#", machineID)
	type test struct {
		machineID         string
		pk                string
		sk                string
		dbError           bool
		expectedError     string
		expectedBatchSize int
		expectedCleanSync bool
		expectedDataType  types.DataType
	}

	cases := []test{
		{
			machineID:         machineID,
			pk:                pkSyncState,
			sk:                "SyncState",
			dbError:           false,
			expectedBatchSize: 31,
			expectedCleanSync: false,
			expectedDataType:  types.DataTypeSyncState,
		},
		{
			machineID:     machineID,
			pk:            pkSyncState,
			sk:            "SyncState",
			dbError:       true,
			expectedError: "failed to retrieve sync state items",
		},
	}

	for _, test := range cases {
		dataTypeAV, _ := test.expectedDataType.MarshalDynamoDBAttributeValue()
		dynamodb := getSyncState(
			func(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
				if test.dbError {
					return nil, errors.New("failed to retrieve sync state items")
				}

				switch key.PartitionKey {
				case test.pk:
					if key.SortKey == test.sk {
						return &awsdynamodb.GetItemOutput{
							Item: map[string]awstypes.AttributeValue{
								"CleanSync": &awstypes.AttributeValueMemberBOOL{Value: false},
								"BatchSize": &awstypes.AttributeValueMemberN{Value: "31"},
								"DataType":  dataTypeAV,
							},
						}, nil
					}
				}

				return &awsdynamodb.GetItemOutput{}, nil
			},
		)
		result, err := GetByMachineID(dynamodb, test.machineID)
		if test.expectedError != "" {
			assert.NotEmpty(t, err)
			assert.Equal(t, test.expectedError, err.Error())
		}

		if test.expectedBatchSize != 0 {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedBatchSize, result.BatchSize)
		}

		if test.expectedDataType != "" {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedDataType, result.DataType)
		}

		if test.expectedCleanSync != false {
			assert.NotEmpty(t, result)
			assert.Equal(t, test.expectedCleanSync, result.CleanSync)
		}
	}
}
