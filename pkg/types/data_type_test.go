package types

import (
	"errors"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// SensorData - Success Validation //
func TestTypes_DataType_Marshal_SensorData_Success(t *testing.T) {
	dataType, err := DataType("SensorData").MarshalText()
	assert.Empty(t, err)
	assert.Equal(t, dataType, []byte("SENSORDATA"))
}

func TestTypes_DataType_Unmarshal_DataTypeSensorData_Success(t *testing.T) {
	dataType := DataTypeSensorData
	err := dataType.UnmarshalText([]byte("SENSORDATA"))
	assert.Empty(t, err)
}

func TestTypes_DataType_Unmarshal_DataTypeSENSOR_STATE_Success(t *testing.T) {
	dataType := DataTypeSensorData
	err := dataType.UnmarshalText([]byte("SENSOR_DATA"))
	assert.Empty(t, err)
}

// SensorData - Failure Validation //
func TestTypes_DataType_Marshal_SensorData_Failure(t *testing.T) {
	_, err := DataType("SensorDatas").MarshalText()
	assert.NotEmpty(t, err)
}

func TestTypes_DataType_Unmarshal_DataTypeSensorData_Failure(t *testing.T) {
	dataType := DataTypeSensorData
	err := dataType.UnmarshalText([]byte("SENSORDATAS"))
	assert.NotEmpty(t, err)
}

// SyncState - Success Validation //
func TestTypes_DataType_Marshal_SyncState_Success(t *testing.T) {
	dataType, err := DataType("SyncState").MarshalText()
	assert.Empty(t, err)
	assert.Equal(t, dataType, []byte("SYNCSTATE"))
}

func TestTypes_DataType_Unmarshal_SyncState_Success(t *testing.T) {
	dataType := DataTypeSyncState
	err := dataType.UnmarshalText([]byte("SYNCSTATE"))
	assert.Empty(t, err)
}

func TestTypes_DataType_Unmarshal_SYNC_STATE_Success(t *testing.T) {
	dataType := DataTypeSyncState
	err := dataType.UnmarshalText([]byte("SYNC_STATE"))
	assert.Empty(t, err)
}

// RulesFeed - Success Validation //
func TestTypes_DataType_Marshal_RulesFeed_Success(t *testing.T) {
	dataType, err := DataType("RulesFeed").MarshalText()
	assert.Empty(t, err)
	assert.Equal(t, dataType, []byte("RULESFEED"))
}

func TestTypes_DataType_Unmarshal_RulesFeed_Success(t *testing.T) {
	dataType := DataTypeRulesFeed
	err := dataType.UnmarshalText([]byte("RULESFEED"))
	assert.Empty(t, err)
}

func TestTypes_DataType_Unmarshal_RULES_FEED_Success(t *testing.T) {
	dataType := DataTypeRulesFeed
	err := dataType.UnmarshalText([]byte("RULES_FEED"))
	assert.Empty(t, err)
}

func TestTypes_DataType_Marshal_SyncState_Failure(t *testing.T) {
	_, err := DataType("SYNCSTATE").MarshalText()
	assert.NotEmpty(t, err)
}

// RulesFeed - Failure Validation //
func TestTypes_DataType_Unmarshal_RulesFeed_Failure(t *testing.T) {
	dataType := DataTypeRulesFeed
	err := dataType.UnmarshalText([]byte("RULESFEEDS"))
	assert.NotEmpty(t, err)
}

func TestTypes_DataType_Marshal_RulesFeed_Failure(t *testing.T) {
	_, err := DataType("RULESFEED").MarshalText()
	assert.NotEmpty(t, err)
}

// SyncState - Failure Validation //
func TestTypes_DataType_Unmarshal_SyncState_Failure(t *testing.T) {
	dataType := DataTypeSensorData
	err := dataType.UnmarshalText([]byte("SYNCSTATESS"))
	assert.NotEmpty(t, err)
}

// MachineConfig - Success Validation //
func TestTypes_DataType_Marshal_MachineConfig_Success(t *testing.T) {
	dataType, err := DataType("MachineConfig").MarshalText()
	assert.Empty(t, err)
	assert.Equal(t, dataType, []byte("MACHINECONFIG"))
}

func TestTypes_DataType_Unmarshal_MachineConfig_Success(t *testing.T) {
	dataType := DataTypeMachineConfig
	err := dataType.UnmarshalText([]byte("MACHINECONFIG"))
	assert.Empty(t, err)
}

func TestTypes_DataType_Unmarshal_MACHINE_CONFIG_Success(t *testing.T) {
	dataType := DataTypeMachineConfig
	err := dataType.UnmarshalText([]byte("MACHINE_CONFIG"))
	assert.Empty(t, err)
}

// SyncState - Failure Validation //
func TestTypes_DataType_Marshal_MachineConfig_Failure(t *testing.T) {
	_, err := DataType("MACHINECONFIGS").MarshalText()
	assert.NotEmpty(t, err)
}

func TestTypes_DataType_Unmarshal_MachineConfig_Failure(t *testing.T) {
	dataType := DataTypeSensorData
	err := dataType.UnmarshalText([]byte("MACHINECONFIGS"))
	assert.NotEmpty(t, err)
}

// MarshalDynamoDBAttributeValue - Success Validation //
func TestTypes_DataType_MarshalDynamoDBAttributeValue_SensorState_Success(t *testing.T) {
	var dataType DataType = DataTypeSensorData

	av, err := dataType.MarshalDynamoDBAttributeValue()
	expectedAV := &awstypes.AttributeValueMemberS{Value: "SensorData"}
	assert.Empty(t, err)
	assert.Equal(t, av, expectedAV)
}

func TestTypes_DataType_MarshalDynamoDBAttributeValue_SyncState_Success(t *testing.T) {
	var dataType DataType = DataTypeSyncState

	av, err := dataType.MarshalDynamoDBAttributeValue()
	expectedAV := &awstypes.AttributeValueMemberS{Value: "SyncState"}
	assert.Empty(t, err)
	assert.Equal(t, av, expectedAV)
}

func TestTypes_DataType_MarshalDynamoDBAttributeValue_MachineConfig_Success(t *testing.T) {
	var dataType DataType = DataTypeMachineConfig

	av, err := dataType.MarshalDynamoDBAttributeValue()
	expectedAV := &awstypes.AttributeValueMemberS{Value: "MachineConfig"}
	assert.Empty(t, err)
	assert.Equal(t, av, expectedAV)
}

// MarshalDynamoDBAttributeValue - Failure Validation //
func TestTypes_DataType_MarshalDynamoDBAttributeValue_MachineConfigs_Failure(t *testing.T) {
	var dataType DataType = "MachineConfigs"

	_, err := dataType.MarshalDynamoDBAttributeValue()
	assert.NotEmpty(t, err)
	expectedErr := errors.New(`unknown data_type value "MachineConfigs"`)
	assert.Equal(t, err.Error(), expectedErr.Error())
}
