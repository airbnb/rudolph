package sensordata

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

// GetSensorData gets a previously stored sensor data
// BUG(derek.wang) In a NoSQL database like DynamoDB the whole idea of save-then-get is an antipattern.
//
//	It requires this call to be strongly consistent with doubles the cost.
//	We should think about a different way of designing the data model.
func GetSensorData(client dynamodb.GetItemAPI, machineID string) (sensorData *SensorData, err error) {
	pk, sk := MachineIDSensorDataPKSK(machineID)

	output, err := client.GetItem(dynamodb.PrimaryKey{
		PartitionKey: pk,
		SortKey:      sk,
	}, false)

	if err != nil {
		return
	}

	if len(output.Item) == 0 {
		return
	}

	err = attributevalue.UnmarshalMap(output.Item, &sensorData)

	if err != nil {
		err = fmt.Errorf("succeeded GetItem but failed to unmarshalMap into output interface: %w", err)
		return
	}

	return
}
