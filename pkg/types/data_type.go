package types

import (
	"fmt"

	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DataType identifies the current DynamoDB data model
type DataType string

const (
	DataTypeSensorData    DataType = "SensorData"
	DataTypeSyncState     DataType = "SyncState"
	DataTypeGlobalConfig  DataType = "GlobalConfig"
	DataTypeMachineConfig DataType = "MachineConfig"
	DataTypeRulesFeed     DataType = "RulesFeed"
)

// UnmarshalText
func (dt *DataType) UnmarshalText(text []byte) error {
	switch mode := string(text); mode {
	case "SENSOR_DATA":
		fallthrough
	case "SENSORDATA":
		*dt = DataTypeSensorData
	case "RULES_FEED":
		fallthrough
	case "RULESFEED":
		*dt = DataTypeRulesFeed
	case "SYNC_STATE":
		fallthrough
	case "SYNCSTATE":
		*dt = DataTypeSyncState
	case "MACHINE_CONFIG":
		fallthrough
	case "MACHINECONFIG":
		*dt = DataTypeGlobalConfig
	case "GLOBAL_CONFIG":
		fallthrough
	case "GLOBALCONFIG":
		*dt = DataTypeGlobalConfig
	default:
		return fmt.Errorf("unknown data_type value %q", mode)
	}
	return nil
}

// MarshalText
func (dt DataType) MarshalText() ([]byte, error) {
	switch dt {
	case DataTypeSensorData:
		return []byte("SENSORDATA"), nil
	case DataTypeSyncState:
		return []byte("SYNCSTATE"), nil
	case DataTypeMachineConfig:
		return []byte("MACHINECONFIG"), nil
	case DataTypeGlobalConfig:
		return []byte("GLOBALCONFIG"), nil
	case DataTypeRulesFeed:
		return []byte("RULESFEED"), nil
	default:
		return nil, fmt.Errorf("unknown data_type %s", dt)
	}
}

// MarshalDynamoDBAttributeValue implements the Marshal interface
func (dt DataType) MarshalDynamoDBAttributeValue() (awstypes.AttributeValue, error) {
	var s string
	switch dt {
	case DataTypeSensorData:
		s = "SensorData"
	case DataTypeSyncState:
		s = "SyncState"
	case DataTypeMachineConfig:
		s = "MachineConfig"
	case DataTypeGlobalConfig:
		s = "GlobalConfig"
	case DataTypeRulesFeed:
		s = "RulesFeed"
	default:
		return nil, fmt.Errorf("unknown data_type value %q", dt)
	}
	return &awstypes.AttributeValueMemberS{Value: s}, nil
}

// UnmarshalDynamoDBAttributeValue implements the Unmarshaler interface
func (dt *DataType) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	switch t := aws.StringValue(av.N); t {
	case "1":
		fallthrough
	case "SENSOR_DATA":
		fallthrough
	case "SENSORDATA":
		*dt = DataTypeSensorData
	case "2":
		fallthrough
	case "SYNC_STATE":
		fallthrough
	case "SYNCSTATE":
		*dt = DataTypeSyncState
	case "3":
		fallthrough
	case "MACHINE_CONFIG":
		*dt = DataTypeMachineConfig
	default:
		return fmt.Errorf("unknown data_type value %q", t)
	}
	return nil
}
