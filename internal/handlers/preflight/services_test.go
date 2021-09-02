package preflight

import (
	"fmt"
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/sensordata"
	rudolphtypes "github.com/airbnb/rudolph/pkg/types"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

type mockPutter func(item interface{}) (*awsdynamodb.PutItemOutput, error)

func (m mockPutter) PutItem(item interface{}) (*awsdynamodb.PutItemOutput, error) { return m(item) }

var _ dynamodb.PutItemAPI = mockPutter(nil)

func Test_concreteSensorDataSaver_OK(t *testing.T) {
	type test struct {
		comment              string
		hostname             string
		machineID            string
		serialNumber         string
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
		expectedDataType     rudolphtypes.DataType
	}

	timeProvider := clock.ConcreteTimeProvider{}
	machineID := "AAAAAAAA-A00A-1234-1234-5864377B4831"
	pk, sk := sensordata.MachineIDSensorDataPKSK(machineID)

	expected := test{
		comment:              fmt.Sprintf("%s %s", "Testing", machineID),
		hostname:             "macbook.pro.localhost",
		machineID:            machineID,
		osBuild:              "20A21",
		osVersion:            "12.34",
		santaVersion:         "2021.1",
		serialNumber:         "123456789ABC",
		binaryRuleCount:      4,
		certRuleCount:        3,
		transitiveRuleCount:  2,
		compilerRuleCount:    1,
		ruleCount:            10,
		primaryUser:          "john_doe",
		expectedTime:         clock.RFC3339(timeProvider.Now()),
		expectedExpiresAfter: clock.Unixtimestamp(timeProvider.Now().UTC().AddDate(0, 0, 90)),
		expectedDataType:     rudolphtypes.DataTypeSensorData,
	}

	stateTrackingService := concreteStateTrackingService{
		putter: mockPutter(
			func(item interface{}) (*awsdynamodb.PutItemOutput, error) {
				sensorData := item.(sensordata.SensorData)
				assert.Equal(t, expected.machineID, sensorData.MachineID)
				assert.Equal(t, expected.osBuild, sensorData.OSBuild)
				assert.Equal(t, expected.osVersion, sensorData.OSVersion)
				assert.Equal(t, expected.binaryRuleCount, sensorData.BinaryRuleCount)
				assert.Equal(t, expected.certRuleCount, sensorData.CertificateRuleCount)
				assert.Equal(t, expected.compilerRuleCount, sensorData.CompilerRuleCount)
				assert.Equal(t, expected.transitiveRuleCount, sensorData.TransitiveRuleCount)
				assert.Equal(t, expected.ruleCount, sensorData.RuleCount)
				assert.Equal(t, expected.primaryUser, sensorData.PrimaryUser)
				assert.Equal(t, expected.serialNumber, sensorData.SerialNum)
				assert.Equal(t, expected.expectedDataType, sensorData.DataType)
				// assert.Equal(t, expected.santaVersion, sensorData.SantaVersion)
				assert.Equal(t, pk, sensorData.PartitionKey)
				assert.Equal(t, sk, sensorData.SortKey)
				// assert.Equal(t, "osversion", sensorData.Hostname) // Missing??
				// assert.Equal(t, types.Lockdown, sensorData.ClientMode) // missing??

				return &awsdynamodb.PutItemOutput{}, nil
			},
		),
		timeProvider: timeProvider,
	}
	request := &PreflightRequest{
		OSBuild:              expected.osBuild,
		OSVersion:            expected.osVersion,
		Hostname:             expected.hostname,
		SantaVersion:         expected.santaVersion,
		ClientMode:           rudolphtypes.Lockdown,
		BinaryRuleCount:      expected.binaryRuleCount,
		CertificateRuleCount: expected.certRuleCount,
		TransitiveRuleCount:  expected.transitiveRuleCount,
		CompilerRuleCount:    expected.compilerRuleCount,
		PrimaryUser:          expected.primaryUser,
		SerialNumber:         expected.serialNumber,
	}

	err := stateTrackingService.saveSensorDataFromPreflightRequest(machineID, request)

	assert.Empty(t, err)
}
