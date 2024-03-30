package types

import (
	"fmt"

	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		fallthrough
	case "SensorData":
		*dt = DataTypeSensorData
	case "RULES_FEED":
		fallthrough
	case "RULESFEED":
		fallthrough
	case "RulesFeed":
		*dt = DataTypeRulesFeed
	case "SYNC_STATE":
		fallthrough
	case "SYNCSTATE":
		fallthrough
	case "SyncState":
		*dt = DataTypeSyncState
	case "MACHINE_CONFIG":
		fallthrough
	case "MACHINECONFIG":
		fallthrough
	case "MachineConfig":
		*dt = DataTypeMachineConfig
	case "GLOBAL_CONFIG":
		fallthrough
	case "GLOBALCONFIG":
		fallthrough
	case "GlobalConfig":
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
		return []byte("SensorData"), nil
	case DataTypeSyncState:
		return []byte("SyncState"), nil
	case DataTypeMachineConfig:
		return []byte("MachineConfig"), nil
	case DataTypeGlobalConfig:
		return []byte("GlobalConfig"), nil
	case DataTypeRulesFeed:
		return []byte("RulesFeed"), nil
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
func (dt *DataType) UnmarshalDynamoDBAttributeValue(av awstypes.AttributeValue) error {
	v, ok := av.(*awstypes.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("unexpected data_type value type: %T", av)
	}

	switch t := v.Value; t {
	case "1":
		fallthrough
	case "SENSOR_DATA":
		fallthrough
	case "SENSORDATA":
		fallthrough
	case "SensorData":
		*dt = DataTypeSensorData
	case "2":
		fallthrough
	case "SYNC_STATE":
		fallthrough
	case "SYNCSTATE":
		fallthrough
	case "SyncState":
		*dt = DataTypeSyncState
	case "3":
		fallthrough
	case "GLOBAL_CONFIG":
		fallthrough
	case "GLOBALCONFIG":
		fallthrough
	case "GlobalConfig":
		*dt = DataTypeGlobalConfig
	case "4":
		fallthrough
	case "MACHINE_CONFIG":
		fallthrough
	case "MACHINECONFIG":
		fallthrough
	case "MachineConfig":
		*dt = DataTypeMachineConfig
	case "5":
		fallthrough
	case "RULES_FEED":
		fallthrough
	case "RULESFEED":
		fallthrough
	case "RulesFeed":
		*dt = DataTypeRulesFeed
	default:
		return fmt.Errorf("unknown data_type value %q", t)
	}

	return nil
}
