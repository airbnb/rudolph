package machineconfiguration

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
)

const (
	machineConfigurationPKPrefix        = "Machine#"
	globalConfigurationPK               = "GlobalConfig"
	currentSK                           = "Config"
	allowGlobalLockdown                 = false
	DefaultFullSyncInterval             = 600
	DefaultSyncTypeNormal        string = SyncTypeNormal
	SyncTypeNormal               string = "normal"
	SyncTypeClean                string = "clean"
	SyncTypeCleanAll             string = "clean_all"
)

// MachineConfigurationRow is an encapsulation of a DynamoDB row containing machine configuration data
type MachineConfigurationRow struct {
	dynamodb.PrimaryKey
	MachineConfiguration
}

// MachineConfiguration is the abstract notion, sans DynamoDB magic (e.g. PK/SK)
// Use Santa defined constants
// https://github.com/google/santa/blob/main/Source/santactl/Commands/sync/SNTCommandSyncConstants.m#L32-L35
type MachineConfiguration struct {
	ClientMode             types.ClientMode `dynamodbav:"ClientMode"`
	BlockedPathRegex       string           `dynamodbav:"BlockedPathRegex"`
	AllowedPathRegex       string           `dynamodbav:"AllowedPathRegex"`
	BatchSize              int              `dynamodbav:"BatchSize"`
	EnableBundles          bool             `dynamodbav:"EnableBundles"`
	EnabledTransitiveRules bool             `dynamodbav:"EnableTransitiveRules"`
	CleanSync              bool             `dynamodbav:"CleanSync,omitempty"`
	FullSyncInterval       int              `dynamodbav:"FullSyncInterval,omitempty"`
	UploadLogsURL          string           `dynamodbav:"UploadLogsUrl,omitempty"`
	BlockUsbMount          bool             `dynamodbav:"BlockUsbMount,omitempty"`
	RemountUsbMode         string           `dynamodbav:"RemountUsbMode,omitempty"`
	// SyncType                 types.SyncType   `dynamodbav:"SyncType,omitempty"`
	OverrideFileAccessAction string         `dynamodbav:"OverrideFileAccessAction,omitempty"`
	DataType                 types.DataType `dynamodbav:"DataType,omitempty"`
}

type MachineConfigurationUpdateRequest struct {
	ClientMode            *types.ClientMode
	BlockedPathRegex      *string
	AllowedPathRegex      *string
	BatchSize             *int
	EnableBundles         *bool
	EnableTransitiveRules *bool
	CleanSync             *bool
	FullSyncInterval      *int
	BlockUsbMount         *bool
	RemountUsbMode        *string
	// SyncType                 *types.SyncType
	OverrideFileAccessAction *string
	UploadLogsURL            *string
}

// Fragments for updates
type updateClientMode struct {
	ClientMode types.ClientMode `dynamodbav:"ClientMode"`
}

func machineConfigurationPK(machineID string) string {
	return fmt.Sprintf("%s%s", machineConfigurationPKPrefix, machineID)
}

func machineConfigurationSK() string {
	return currentSK
}

// Default Universal Config
func GetUniversalDefaultConfig() MachineConfiguration {
	return MachineConfiguration{
		ClientMode:             types.Monitor,
		BlockedPathRegex:       "",
		AllowedPathRegex:       "",
		BatchSize:              50,
		EnableBundles:          false,
		EnabledTransitiveRules: false,
		CleanSync:              false,
		FullSyncInterval:       DefaultFullSyncInterval,
		UploadLogsURL:          "",
		BlockUsbMount:          false,
		RemountUsbMode:         "",
		// SyncType:                 types.SyncTypeNormal,
		OverrideFileAccessAction: "",
		DataType:                 types.DataTypeGlobalConfig,
	}
}
