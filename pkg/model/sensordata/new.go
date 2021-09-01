package sensordata

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
)

func NewSensorData(timeProvider clock.TimeProvider,
	machineID string,
	serialNumber string,
	osVersion string,
	osBuild string,
	requestCleanSync bool,
	primaryUser string,
	certRuleCount int,
	binaryRuleCount int,
	compilerRuleCount int,
	transitiveRuleCount int,
) SensorData {
	pk, sk := MachineIDSensorDataPKSK(machineID)
	return SensorData{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: pk,
			SortKey:      sk,
		},
		MachineID:            machineID,
		SerialNum:            serialNumber,
		OSVersion:            osVersion,
		OSBuild:              osBuild,
		RequestCleanSync:     requestCleanSync,
		PrimaryUser:          primaryUser,
		RuleCount:            certRuleCount + binaryRuleCount + compilerRuleCount + transitiveRuleCount,
		BinaryRuleCount:      binaryRuleCount,
		CertificateRuleCount: certRuleCount,
		CompilerRuleCount:    compilerRuleCount,
		TransitiveRuleCount:  transitiveRuleCount,
		Time:                 clock.RFC3339(timeProvider.Now()),
		ExpiresAfter:         GetSensorDataExpiresAfter(timeProvider),
		DataType:             GetDataType(),
	}
}
