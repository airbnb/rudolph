package machineconfiguration

import (
	"log"
	"strings"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/pkg/errors"
)

// SetGlobalConfig writes a global configuration to the DynamoDB table via a PutItem API call
func SetGlobalConfig(client dynamodb.PutItemAPI, clientMode types.ClientMode, blockedPathRegex string, allowedPathRegex string, batchSize int, isEnableBundles bool, isEnabledTransitiveRules bool, fullSyncInterval int, uploadLogsURL string) (err error) {
	// This is a safety check to prevent the accidental setting of lockdown mode globally via a config set operation
	// Nothing is preventing this from being manually performed at the DyanmoDB table itself
	if !allowGlobalLockdown && clientMode == types.Lockdown {
		return errors.New("global lockdown configuration is disabled right now")
	}

	// Construct a MachineConfigRow to represent a GlobalConfig
	globalConfigRow := buildConfig(
		globalConfigurationPK,
		clientMode,
		blockedPathRegex,
		allowedPathRegex,
		batchSize,
		isEnableBundles,
		isEnabledTransitiveRules,
		false,
		fullSyncInterval,
		uploadLogsURL,
	)

	_, err = client.PutItem(globalConfigRow)
	if err != nil {
		log.Print(errors.Wrap(err, " setting global config failed"))
	}

	return
}

// SetMachineConfig writes a machineID specific configuration to the DynamoDB table via a PutItem API call
func SetMachineConfig(client dynamodb.PutItemAPI, machineID string, clientMode types.ClientMode, blockedPathRegex string, allowedPathRegex string, batchSize int, isEnableBundles bool, isEnabledTransitiveRules bool, isCleanSync bool, fullSyncInterval int, uploadLogsURL string) (err error) {
	// Create the machineID specific PK
	machinePK := machineConfigurationPK(machineID)

	// Construct a MachineConfigRow to represent a MachineConfig
	machineConfigRow := buildConfig(
		machinePK,
		clientMode,
		blockedPathRegex,
		allowedPathRegex,
		batchSize,
		isEnableBundles,
		isEnabledTransitiveRules,
		isCleanSync,
		fullSyncInterval,
		uploadLogsURL,
	)

	_, err = client.PutItem(machineConfigRow)
	if err != nil {
		log.Print(errors.Wrap(err, " setting global config failed"))
	}

	return
}

// buildConfig constructs a MachineConfigurationRow which will represent either a global or machineID specific configuration set to be used in a DynamoDB PutItem API call
func buildConfig(pk string, clientMode types.ClientMode, blockedPathRegex string, allowedPathRegex string, batchSize int, isEnableBundles bool, isEnabledTransitiveRules bool, isCleanSync bool, fullSyncInterval int, uploadLogsURL string) (configRow *MachineConfigurationRow) {
	// Check batchsize just to make sure at least a valid value is provided if not a positive int value
	if batchSize <= 0 {
		batchSize = 50
	}

	config := MachineConfiguration{
		ClientMode:             clientMode,
		BlockedPathRegex:       blockedPathRegex,
		AllowedPathRegex:       allowedPathRegex,
		BatchSize:              batchSize,
		EnableBundles:          isEnableBundles,
		EnabledTransitiveRules: isEnabledTransitiveRules,
		CleanSync:              isCleanSync,
		FullSyncInterval:       fullSyncInterval,
		UploadLogsURL:          uploadLogsURL,
	}

	if strings.Compare(pk, globalConfigurationPK) == 0 {
		config.DataType = types.DataTypeGlobalConfig
	} else {
		config.DataType = types.DataTypeMachineConfig
	}

	configRow = &MachineConfigurationRow{
		dynamodb.PrimaryKey{
			PartitionKey: pk,
			SortKey:      machineConfigurationSK(),
		},
		config,
	}
	return
}

func GetConfigurationSetter(client dynamodb.PutItemAPI, timeProvider clock.TimeProvider) ConcreteConfigurationSetter {
	return ConcreteConfigurationSetter{
		global: ConcreteGlobalConfigurationSetter{
			setter: client,
			cache:  GetCache(timeProvider),
		},
		machine: ConcreteMachineConfigurationSetter{
			setter: client,
		},
	}
}

//
// ConfigurationSetter is an interface to set and override the all contents of a Global or MachineConfiguration.
// This can be used to initially create a configuration for either a global or machine configuration.
//
// type ConfigurationSetter interface {
// 	SetGlobalConfig(config MachineConfiguration) error
// 	SetMachineConfig(machineID string, config MachineConfiguration) error
// }

type ConcreteConfigurationSetter struct {
	global  GlobalConfigurationSetter
	machine MachineConfigurationSetter
}

// func (f ConcreteConfigurationSetter) SetGlobalConfig(config MachineConfiguration) error {
// 	err := f.global.setGlobalConfig(config)
// 	return err
// }

// func (f ConcreteConfigurationSetter) SetMachineConfig(machineID string, config MachineConfiguration) error {
// 	err := f.machine.setMachineConfig(machineID, config)
// 	return err
// }

// GlobalConfigurationSetter
type GlobalConfigurationSetter interface {
	setGlobalConfig(config MachineConfiguration) error
}

type ConcreteGlobalConfigurationSetter struct {
	setter dynamodb.PutItemAPI
	cache  Cache
}

func (c ConcreteGlobalConfigurationSetter) setGlobalConfig(config MachineConfiguration) error {
	if !allowGlobalLockdown && config.ClientMode == types.Lockdown {
		return errors.New("global lockdown configuration is disabled right now")
	}

	// Construct a MachineConfigRow to represent a GlobalConfig
	globalConfigRow := buildConfig(
		globalConfigurationPK,
		config.ClientMode,
		config.BlockedPathRegex,
		config.AllowedPathRegex,
		config.BatchSize,
		config.EnableBundles,
		config.EnabledTransitiveRules,
		config.CleanSync,
		config.FullSyncInterval,
		config.UploadLogsURL,
	)

	_, err := c.setter.PutItem(globalConfigRow)
	if err != nil {
		log.Print(errors.Wrap(err, " setting global config failed"))
	}

	// Attempt to set the new global configuration into cache
	_ = c.cache.Set(CacheKeyGlobal, &config)

	return nil
}

// MachineConfigurationSetter
type MachineConfigurationSetter interface {
	setMachineConfig(machineID string, config MachineConfiguration) error
}

type ConcreteMachineConfigurationSetter struct {
	setter dynamodb.PutItemAPI
}

func (c ConcreteMachineConfigurationSetter) setMachineConfig(machineID string, config MachineConfiguration) error {
	// Create the machineID specific PK
	machinePK := machineConfigurationPK(machineID)

	// Construct a MachineConfigRow to represent a MachineConfig
	machineConfigRow := buildConfig(
		machinePK,
		config.ClientMode,
		config.BlockedPathRegex,
		config.AllowedPathRegex,
		config.BatchSize,
		config.EnableBundles,
		config.EnabledTransitiveRules,
		config.CleanSync,
		config.FullSyncInterval,
		config.UploadLogsURL,
	)

	_, err := c.setter.PutItem(machineConfigRow)
	if err != nil {
		log.Print(errors.Wrap(err, " setting machine config failed"))
	}
	return nil
}
