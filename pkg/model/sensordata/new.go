package sensordata

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	rudolphtypes "github.com/airbnb/rudolph/pkg/types"
)

func NewSensorData(
	timeProvider clock.TimeProvider,
	machineID string,
	serialNumber string,
	osVersion string,
	osBuild string,
	santaVersion string,
	clientMode rudolphtypes.ClientMode,
	requestCleanSync bool,
	primaryUser string,
	certRuleCount int,
	binaryRuleCount int,
	cdHashRuleCount int,
	teamIDRuleCount int,
	signingIDRuleCount int,
	compilerRuleCount int,
	transitiveRuleCount int,
) SensorData {
	pk, sk := MachineIDSensorDataPKSK(machineID)
	var totalRuleCount int
	totalRuleCount += certRuleCount + binaryRuleCount + compilerRuleCount + transitiveRuleCount + cdHashRuleCount + teamIDRuleCount + signingIDRuleCount
	return SensorData{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: pk,
			SortKey:      sk,
		},
		MachineID:            machineID,
		SerialNum:            serialNumber,
		OSVersion:            osVersion,
		OSBuild:              osBuild,
		SantaVersion:         santaVersion,
		ClientMode:           clientMode,
		RequestCleanSync:     requestCleanSync,
		PrimaryUser:          primaryUser,
		RuleCount:            totalRuleCount,
		BinaryRuleCount:      binaryRuleCount,
		CertificateRuleCount: certRuleCount,
		CDHashRuleCount:      cdHashRuleCount,
		TeamIDRuleCount:      teamIDRuleCount,
		SigningIDRuleCount:   signingIDRuleCount,
		CompilerRuleCount:    compilerRuleCount,
		TransitiveRuleCount:  transitiveRuleCount,
		Time:                 clock.RFC3339(timeProvider.Now()),
		ExpiresAfter:         GetSensorDataExpiresAfter(timeProvider),
		DataType:             GetDataType(),
	}
}
