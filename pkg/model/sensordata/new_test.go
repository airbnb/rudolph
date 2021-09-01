package sensordata

import (
	"fmt"
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/types"
	rudolphtype "github.com/airbnb/rudolph/pkg/types"
	"github.com/stretchr/testify/assert"
)

var timeProvider = clock.ConcreteTimeProvider{}

func Test_CreateSyncState(t *testing.T) {
	type test struct {
		comment              string
		hostname             string
		machineID            string
		serialNumber         string
		requestCleanSync     bool
		osBuild              string
		osVersion            string
		santaVersion         string
		binaryRuleCount      int
		certRuleCount        int
		transitiveRuleCount  int
		compilerRuleCount    int
		ruleCount            int
		primaryUser          string
		expectedTime         string
		expectedExpiresAfter int64
		expectedDataType     types.DataType
	}

	machineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	pk, sk := MachineIDSensorDataPKSK(machineID)

	expected := test{
		comment:              fmt.Sprintf("%s %s", "Testing", machineID),
		hostname:             "macbook.pro.localhost",
		machineID:            machineID,
		osBuild:              "20A21",
		osVersion:            "12.34",
		santaVersion:         "2021.1",
		serialNumber:         "123456789ABC",
		requestCleanSync:     false,
		binaryRuleCount:      4,
		certRuleCount:        3,
		transitiveRuleCount:  2,
		compilerRuleCount:    1,
		ruleCount:            10,
		primaryUser:          "john_doe",
		expectedTime:         clock.RFC3339(timeProvider.Now()),
		expectedExpiresAfter: clock.Unixtimestamp(timeProvider.Now().UTC().AddDate(0, 0, 90)),
		expectedDataType:     rudolphtype.DataTypeSensorData,
	}

	sensorData := NewSensorData(
		timeProvider,
		expected.machineID,
		expected.serialNumber,
		expected.osVersion,
		expected.osBuild,
		expected.requestCleanSync,
		expected.primaryUser,
		expected.certRuleCount,
		expected.binaryRuleCount,
		expected.compilerRuleCount,
		expected.transitiveRuleCount,
	)
	assert.Equal(t, expected.machineID, sensorData.MachineID)
	assert.Equal(t, expected.serialNumber, sensorData.SerialNum)
	assert.Equal(t, expected.requestCleanSync, sensorData.RequestCleanSync)
	assert.Equal(t, expected.osBuild, sensorData.OSBuild)
	assert.Equal(t, expected.osVersion, sensorData.OSVersion)
	assert.Equal(t, expected.binaryRuleCount, sensorData.BinaryRuleCount)
	assert.Equal(t, expected.certRuleCount, sensorData.CertificateRuleCount)
	assert.Equal(t, expected.compilerRuleCount, sensorData.CompilerRuleCount)
	assert.Equal(t, expected.transitiveRuleCount, sensorData.TransitiveRuleCount)
	assert.Equal(t, expected.ruleCount, sensorData.RuleCount)
	assert.Equal(t, expected.primaryUser, sensorData.PrimaryUser)
	assert.Equal(t, pk, sensorData.PartitionKey)
	assert.Equal(t, sk, sensorData.SortKey)
	assert.Equal(t, expected.expectedDataType, sensorData.DataType)
}
