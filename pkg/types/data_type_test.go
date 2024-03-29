package types

import (
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestDataTypes_MarshalTest(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		want     []byte
		wantErr  bool
	}{
		{"SensorData", DataTypeSensorData, []byte(DataTypeSensorData), false},
		{"SyncState", DataTypeSyncState, []byte(DataTypeSyncState), false},
		{"RulesFeed", DataTypeRulesFeed, []byte(DataTypeRulesFeed), false},
		{"MachineConfig", DataTypeMachineConfig, []byte(DataTypeMachineConfig), false},
		{"GlobalConfig", DataTypeGlobalConfig, []byte(DataTypeGlobalConfig), false},
		{"MISSPELLED", DataType(""), []byte(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dataType.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataType.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDataTypes_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    DataType
		wantErr bool
	}{
		{"SensorData", []byte(DataTypeSensorData), DataTypeSensorData, false},
		{"SyncState", []byte(DataTypeSyncState), DataTypeSyncState, false},
		{"RulesFeed", []byte(DataTypeRulesFeed), DataTypeRulesFeed, false},
		{"MachineConfig", []byte(DataTypeMachineConfig), DataTypeMachineConfig, false},
		{"GlobalConfig", []byte(DataTypeGlobalConfig), DataTypeGlobalConfig, false},
		{"MISSPELLED", []byte(""), DataType(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DataType
			err := dt.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataType.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, dt)
		})
	}
}

func TestDataTypes_MarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		want     awstypes.AttributeValue
		wantErr  bool
	}{
		{"SensorData", DataTypeSensorData, &awstypes.AttributeValueMemberS{Value: string(DataTypeSensorData)}, false},
		{"SyncState", DataTypeSyncState, &awstypes.AttributeValueMemberS{Value: string(DataTypeSyncState)}, false},
		{"RulesFeed", DataTypeRulesFeed, &awstypes.AttributeValueMemberS{Value: string(DataTypeRulesFeed)}, false},
		{"MachineConfig", DataTypeMachineConfig, &awstypes.AttributeValueMemberS{Value: string(DataTypeMachineConfig)}, false},
		{"GlobalConfig", DataTypeGlobalConfig, &awstypes.AttributeValueMemberS{Value: string(DataTypeGlobalConfig)}, false},
		{"MISSPELLED", DataType(""), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dataType.MarshalDynamoDBAttributeValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataType.MarshalDynamoDBAttributeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDataTypes_UnmarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []struct {
		name    string
		av      awstypes.AttributeValue
		want    DataType
		wantErr bool
	}{
		{"SensorData", &awstypes.AttributeValueMemberS{Value: string(DataTypeSensorData)}, DataTypeSensorData, false},
		{"SyncState", &awstypes.AttributeValueMemberS{Value: string(DataTypeSyncState)}, DataTypeSyncState, false},
		{"RulesFeed", &awstypes.AttributeValueMemberS{Value: string(DataTypeRulesFeed)}, DataTypeRulesFeed, false},
		{"MachineConfig", &awstypes.AttributeValueMemberS{Value: string(DataTypeMachineConfig)}, DataTypeMachineConfig, false},
		{"GlobalConfig", &awstypes.AttributeValueMemberS{Value: string(DataTypeGlobalConfig)}, DataTypeGlobalConfig, false},
		{"MISSPELLED", nil, DataType(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DataType("")
			err := got.UnmarshalDynamoDBAttributeValue(tt.av)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataType.UnmarshalDynamoDBAttributeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
