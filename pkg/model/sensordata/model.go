package sensordata

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
)

const (
	sensorDataPKPrefix  = "Machine#"
	sensorDataDefaultSK = "Current"

	sensorDataExpiresAfterInDays int = 90
)

// SensorData encapsulation for a DDB row encapsulating data that is uploaded in preflight
type SensorData struct {
	dynamodb.PrimaryKey
	MachineID            string         `dynamodbav:"MachineID"`
	SerialNum            string         `dynamodbav:"SerialNum"`
	OSVersion            string         `dynamodbav:"OSVersion"`
	OSBuild              string         `dynamodbav:"OSBuild"`
	RequestCleanSync     bool           `dynamodbav:"RequestCleanSync"`
	PrimaryUser          string         `dynamodbav:"PrimaryUser"`
	RuleCount            int            `dynamodbav:"RuleCount"`
	CertificateRuleCount int            `dynamodbav:"CertificateRuleCount"`
	BinaryRuleCount      int            `dynamodbav:"BinaryRuleCount"`
	CompilerRuleCount    int            `dynamodbav:"CompilerRuleCount"`
	TransitiveRuleCount  int            `dynamodbav:"TransitiveRuleCount"`
	Time                 string         `dynamodbav:"Time"`
	ExpiresAfter         int64          `dynamodbav:"ExpiresAfter,omitempty"`
	DataType             types.DataType `dynamodbav:"DataType"`
}

// MachineIDSensorDataPKSK returns the partition and sort keys for a machine id
func MachineIDSensorDataPKSK(machineID string) (string, string) {
	return SensorDataPK(machineID), sensorDataDefaultSK
}

func SensorDataPK(machineID string) string {
	return fmt.Sprintf("%s%s", sensorDataPKPrefix, machineID)
}

func GetSensorDataExpiresAfter(timeProvider clock.TimeProvider) int64 {
	return clock.Unixtimestamp(timeProvider.Now().UTC().AddDate(0, 0, sensorDataExpiresAfterInDays))
}

func GetDataType() types.DataType {
	return types.DataTypeSensorData
}
