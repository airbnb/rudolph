package machineconfiguration

import (
	"log"
	"strings"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/pkg/errors"
)

// UpdateMachineConfigClientMode updates a machineID specific ClientMode configuration to the DynamoDB table via an UpdateItem API call
func UpdateMachineConfigClientMode(client dynamodb.UpdateItemAPI, machineID string, newClientMode types.ClientMode) (err error) {
	pk := dynamodb.PrimaryKey{
		PartitionKey: machineConfigurationPK(machineID),
		SortKey:      machineConfigurationSK(),
	}

	clientMode := updateClientMode{
		ClientMode: newClientMode,
	}

	//UpdateItem takes a key (pk/sk), item name --> match this to the column/struct tag, and interface to set
	_, err = client.UpdateItem(pk, clientMode)
	if err != nil {
		log.Print(errors.Wrap(err, " setting machine config failed"))
	}

	return
}

// UpdateGlobalConfigClientMode updates a machineID specific ClientMode configuration to the DynamoDB table via an UpdateItem API call
func UpdateGlobalConfigClientMode(client dynamodb.UpdateItemAPI, newClientMode types.ClientMode) (err error) {
	pk := dynamodb.PrimaryKey{
		PartitionKey: globalConfigurationPK,
		SortKey:      machineConfigurationSK(),
	}

	clientMode := updateClientMode{
		ClientMode: newClientMode,
	}

	//UpdateItem takes a key (pk/sk), item name --> match this to the column/struct tag, and interface to set
	_, err = client.UpdateItem(pk, clientMode)
	if err != nil {
		log.Print(errors.Wrap(err, " setting global config failed"))
	}

	return
}

func GetConfigurationUpdater(client dynamodb.UpdateItemAPI, fetcher ConcreteConfigurationFetcher, timeProvider clock.TimeProvider) ConcreteConfigurationUpdater {
	return ConcreteConfigurationUpdater{
		global: ConcreteGlobalConfigurationUpdater{
			updater: client,
			cache:   GetCache(timeProvider),
			fetcher: fetcher,
		},
		machine: ConcreteMachineConfigurationUpdater{
			updater: client,
			fetcher: fetcher,
		},
	}
}

//
// ConfigurationUpdater is an interface to accept a new configuration and update either a global or machine configuration
//
type ConfigurationUpdater interface {
	UpdateGlobalConfig(configRequest MachineConfigurationUpdateRequest) (MachineConfiguration, error)
	UpdateMachineConfig(machineID string, configRequest MachineConfigurationUpdateRequest) (MachineConfiguration, error)
}

type ConcreteConfigurationUpdater struct {
	global  GlobalConfigurationUpdater
	machine MachineConfigurationUpdater
}

func (f ConcreteConfigurationUpdater) UpdateGlobalConfig(configRequest MachineConfigurationUpdateRequest) (*MachineConfiguration, error) {
	return f.global.updateConfig(configRequest)
}

// MachineConfigs //

func (f ConcreteConfigurationUpdater) UpdateMachineConfig(machineID string, configRequest MachineConfigurationUpdateRequest) (*MachineConfiguration, error) {
	return f.machine.updateConfig(machineID, configRequest)
}

// GlobalConfigurationUpdater //
type GlobalConfigurationUpdater interface {
	updateConfig(configRequest MachineConfigurationUpdateRequest) (updatedConfig *MachineConfiguration, err error)
}

type ConcreteGlobalConfigurationUpdater struct {
	updater dynamodb.UpdateItemAPI
	cache   Cache
	fetcher ConcreteConfigurationFetcher
}

func (c ConcreteGlobalConfigurationUpdater) updateConfig(configRequest MachineConfigurationUpdateRequest) (updatedConfig *MachineConfiguration, err error) {
	// Create the global configuration PK
	pk := dynamodb.PrimaryKey{
		PartitionKey: globalConfigurationPK,
		SortKey:      machineConfigurationSK(),
	}
	var newGlobalConfig MachineConfiguration
	var changed bool = false

	currentGlobalConfig, _, err := c.fetcher.getIntendedGlobalConfig()
	if err != nil {
		return
	}

	// the new config is a copy of the current config
	newGlobalConfig = currentGlobalConfig

	if configRequest.ClientMode != nil && *configRequest.ClientMode != currentGlobalConfig.ClientMode {
		if *configRequest.ClientMode == types.Lockdown {
			err = errors.New("global lockdown configuration is disabled right now")
			return
		} else {
			newGlobalConfig.ClientMode = *configRequest.ClientMode
			changed = true
		}
	}

	if configRequest.AllowedPathRegex != nil && strings.Compare(*configRequest.AllowedPathRegex, currentGlobalConfig.AllowedPathRegex) != 0 {
		newGlobalConfig.AllowedPathRegex = *configRequest.AllowedPathRegex
		changed = true
	}

	if configRequest.BlockedPathRegex != nil && strings.Compare(*configRequest.BlockedPathRegex, currentGlobalConfig.BlockedPathRegex) != 0 {
		newGlobalConfig.AllowedPathRegex = *configRequest.BlockedPathRegex
		changed = true
	}

	if configRequest.BatchSize != nil && *configRequest.BatchSize != currentGlobalConfig.BatchSize {
		if *configRequest.BatchSize > 0 {
			newGlobalConfig.BatchSize = *configRequest.BatchSize
			changed = true
		}
	}

	if configRequest.EnableBundles != nil && *configRequest.EnableBundles != currentGlobalConfig.EnableBundles {
		newGlobalConfig.EnableBundles = *configRequest.EnableBundles
		changed = true
	}

	if configRequest.EnableTransitiveRules != nil && *configRequest.EnableTransitiveRules != currentGlobalConfig.EnabledTransitiveRules {
		newGlobalConfig.EnabledTransitiveRules = *configRequest.EnableTransitiveRules
		changed = true
	}

	if configRequest.FullSyncInterval != nil && *configRequest.FullSyncInterval != currentGlobalConfig.FullSyncInterval {
		if *configRequest.FullSyncInterval >= 60 {
			newGlobalConfig.FullSyncInterval = *configRequest.FullSyncInterval
			changed = true
		}
	}

	if configRequest.CleanSync != nil {
		newGlobalConfig.CleanSync = *configRequest.CleanSync
		changed = true
	}

	// if no items have been changed, return the same configuration
	if !changed {
		updatedConfig = &newGlobalConfig
		return
	}

	output, err := c.updater.UpdateItem(pk, newGlobalConfig)
	if err != nil {
		err = errors.Wrap(err, " updating global configuration failed")
		return
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &updatedConfig)

	if err != nil {
		err = errors.Wrap(err, "succeeded UpdateItem but failed to unmarshalMap into MachineConfiguration")
		return
	}

	// Attempt to set the new global configuration into cache
	_ = c.cache.Set(CacheKeyGlobal, &newGlobalConfig)

	return
}

// MachineConfig //
type MachineConfigurationUpdater interface {
	updateConfig(machineID string, configRequest MachineConfigurationUpdateRequest) (updatedConfig *MachineConfiguration, err error)
}

type ConcreteMachineConfigurationUpdater struct {
	updater dynamodb.UpdateItemAPI
	fetcher ConcreteConfigurationFetcher
}

func (c ConcreteMachineConfigurationUpdater) updateConfig(machineID string, configRequest MachineConfigurationUpdateRequest) (updatedConfig *MachineConfiguration, err error) {
	// Create the machineID specific PK
	pk := dynamodb.PrimaryKey{
		PartitionKey: machineConfigurationPK(machineID),
		SortKey:      machineConfigurationSK(),
	}

	var newMachineConfig MachineConfiguration
	var changed bool = false

	currentMachineConfig, err := c.fetcher.getIntendedConfig(machineID)
	if err != nil {
		return
	}

	// the new config is a copy of the current config
	newMachineConfig = currentMachineConfig

	if configRequest.ClientMode != nil && *configRequest.ClientMode != currentMachineConfig.ClientMode {
		newMachineConfig.ClientMode = *configRequest.ClientMode
		changed = true
	}

	if configRequest.AllowedPathRegex != nil && strings.Compare(*configRequest.AllowedPathRegex, currentMachineConfig.AllowedPathRegex) != 0 {
		newMachineConfig.AllowedPathRegex = *configRequest.AllowedPathRegex
		changed = true
	}

	if configRequest.BlockedPathRegex != nil && strings.Compare(*configRequest.BlockedPathRegex, currentMachineConfig.BlockedPathRegex) != 0 {
		newMachineConfig.AllowedPathRegex = *configRequest.BlockedPathRegex
		changed = true
	}

	if configRequest.BatchSize != nil && *configRequest.BatchSize != currentMachineConfig.BatchSize {
		if *configRequest.BatchSize > 0 {
			newMachineConfig.BatchSize = *configRequest.BatchSize
			changed = true
		}
	}

	if configRequest.EnableBundles != nil && *configRequest.EnableBundles != currentMachineConfig.EnableBundles {
		newMachineConfig.EnableBundles = *configRequest.EnableBundles
		changed = true
	}

	if configRequest.EnableTransitiveRules != nil && *configRequest.EnableTransitiveRules != currentMachineConfig.EnabledTransitiveRules {
		newMachineConfig.EnabledTransitiveRules = *configRequest.EnableTransitiveRules
		changed = true
	}

	if configRequest.FullSyncInterval != nil && *configRequest.FullSyncInterval != currentMachineConfig.FullSyncInterval {
		if *configRequest.FullSyncInterval >= 60 {
			newMachineConfig.FullSyncInterval = *configRequest.FullSyncInterval
			changed = true
		}
	}

	if configRequest.CleanSync != nil {
		newMachineConfig.CleanSync = *configRequest.CleanSync
		changed = true
	}

	// if no items have been changed, return the same configuration
	if !changed {
		updatedConfig = &newMachineConfig
		return
	}

	output, err := c.updater.UpdateItem(pk, newMachineConfig)
	if err != nil {
		err = errors.Wrap(err, " updating machine configuration failed")
		return
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &updatedConfig)

	if err != nil {
		err = errors.Wrap(err, "succeeded UpdateItem but failed to unmarshalMap into MachineConfiguration")
		return
	}

	return
}
