package machineconfiguration

import (
	"log"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/pkg/errors"
)

func GetConfigurationFetcher(client dynamodb.GetItemAPI, timeProvider clock.TimeProvider) ConcreteConfigurationFetcher {
	return ConcreteConfigurationFetcher{
		global: CachedConcreteGlobalConfigurationFetcher{
			getter: client,
			cache:  GetCache(timeProvider),
		},
		machine: ConcreteMachineConfigurationFetcher{
			getter: client,
		},
		universal: ConcreteUniversalConfigurationProvider{},
	}
}

func GetUncachedConfigurationFetcher(client dynamodb.GetItemAPI, timeProvider clock.TimeProvider) ConcreteConfigurationFetcher {
	return ConcreteConfigurationFetcher{
		global: UncachedConcreteGlobalConfigurationFetcher{
			getter: client,
			cache:  GetCache(timeProvider),
		},
		machine: ConcreteMachineConfigurationFetcher{
			getter: client,
		},
		universal: ConcreteUniversalConfigurationProvider{},
	}
}

//
// ConfigurationFetcher is intended to jump through all of the hoops to get the configuration that is
// intended to be deployed to a particular machine, resolving through defaults, overrides, and fallbacks
// and other logic.
//
// type ConfigurationFetcher interface {
// 	GetIntendedConfig(machineID string) (intendedConfig MachineConfiguration, err error)
// 	GetIntendedGlobalConfig() (intendedConfig MachineConfiguration, err error)
// }

type ConcreteConfigurationFetcher struct {
	global    GlobalConfigurationFetcher
	machine   MachineConfigurationFetcher
	universal UniversalConfigurationProvider
}

// GetIntendedConfig returns the machineconfiguration intended for the given machineID. It searches in the following order:
// - First, looks for a machine-specific config in DynamoDB
// - Second, looks for a global config for all machines in DynamoDB
// - Defaults to returning a hardcoded "universal config"
func (f ConcreteConfigurationFetcher) getIntendedConfig(machineID string) (intendedConfig MachineConfiguration, err error) {
	config, err := f.machine.GetMachineSpecificConfig(machineID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get machine config")
		return
	}

	if config == nil {
		config, err = f.global.GetGlobalConfig()
		if err != nil {
			err = errors.Wrapf(err, "failed to get fallback global config")
			return
		}
	}

	// Check if the result.Item is not nil. If this is nil, no valid global config exists within the DynamoDB and return a generic hardcoded value
	if config == nil {
		intendedConfig = f.universal.GetUniversalDefaultConfig()
	} else {
		intendedConfig = *config
	}

	return
}

// GetIntendedGlobalConfig returns the global configuration. It searches in the following order:
// - First, looks for a global config for all machines in DynamoDB
// - Defaults to returning a hardcoded "universal config"
func (f ConcreteConfigurationFetcher) getIntendedGlobalConfig() (intendedConfig MachineConfiguration, isDefaultConfig bool, err error) {
	config, err := f.global.GetGlobalConfig()
	if err != nil {
		err = errors.Wrapf(err, "failed to get fallback global config")
		return
	}

	// Check if the result.Item is not nil. If this is nil, no valid global config exists within the DynamoDB and return a generic hardcoded value
	if config == nil {
		isDefaultConfig = true
		intendedConfig = f.universal.GetUniversalDefaultConfig()
	} else {
		intendedConfig = *config
	}

	return
}

//
// GlobalConfigurationFetcher is intended to fetch the global configuration that all machines fallback onto when
// they are missing their machine-specific overrides.
//
// Because the global config is not intended to change frequently, we cache this item to save on dynamodb getitem calls.
//

func GetGlobalConfigurationFetcher(client dynamodb.GetItemAPI, timeProvider clock.TimeProvider) GlobalConfigurationFetcher {
	return CachedConcreteGlobalConfigurationFetcher{
		getter: client,
		cache:  GetCache(timeProvider),
	}
}

func GetUncachedGlobalConfigurationFetcher(client dynamodb.GetItemAPI, timeProvider clock.TimeProvider) GlobalConfigurationFetcher {
	return UncachedConcreteGlobalConfigurationFetcher{
		getter: client,
		cache:  GetCache(timeProvider),
	}
}

type GlobalConfigurationFetcher interface {
	GetGlobalConfig() (*MachineConfiguration, error)
}

type CachedConcreteGlobalConfigurationFetcher struct {
	getter dynamodb.GetItemAPI
	cache  Cache
}

type UncachedConcreteGlobalConfigurationFetcher struct {
	getter dynamodb.GetItemAPI
	cache  Cache
}

func (c CachedConcreteGlobalConfigurationFetcher) GetGlobalConfig() (config *MachineConfiguration, err error) {
	if c.cache.Has(CacheKeyGlobal) {
		config = c.cache.Get(CacheKeyGlobal)
		return
	}

	config, err = getItemAsMachineConfiguration(
		c.getter,
		globalConfigurationPK,
		machineConfigurationSK(),
	)
	if err != nil {
		err = errors.Wrapf(err, "failed to get global config")
		return
	}

	// Here, it is perfectly possible to have nil config. In this case, we *should* set nil into the cache
	// as it simply means there are no configurations set.
	c.cache.Set(CacheKeyGlobal, config)

	return config, nil
}

func (c UncachedConcreteGlobalConfigurationFetcher) GetGlobalConfig() (config *MachineConfiguration, err error) {
	config, err = getItemAsMachineConfiguration(
		c.getter,
		globalConfigurationPK,
		machineConfigurationSK(),
	)
	if err != nil {
		err = errors.Wrapf(err, "failed to get global config")
		return
	}
	return
}

//
// MachineConfigurationFetcher is intended to retrieve any machine-specific configuration override, if prevent, for
// the given machineID. This override takes precedence over the global and universal defaults.
//

func GetMachineConfigurationFetcher(client dynamodb.GetItemAPI) MachineConfigurationFetcher {
	return ConcreteMachineConfigurationFetcher{
		getter: client,
	}
}

type MachineConfigurationFetcher interface {
	GetMachineSpecificConfig(machineID string) (*MachineConfiguration, error)
}

type ConcreteMachineConfigurationFetcher struct {
	getter dynamodb.GetItemAPI
}

func (f ConcreteMachineConfigurationFetcher) GetMachineSpecificConfig(machineID string) (config *MachineConfiguration, err error) {
	return getItemAsMachineConfiguration(
		f.getter,
		machineConfigurationPK(machineID),
		machineConfigurationSK(),
	)
}

//
// UniversalConfigurationProvider provides a universal config that all machines can fall back onto, when they neither
// have machine-specific nor global overrides.
//

func GetUniversalConfigurationProvider() UniversalConfigurationProvider {
	return ConcreteUniversalConfigurationProvider{}
}

type UniversalConfigurationProvider interface {
	GetUniversalDefaultConfig() MachineConfiguration
}

type ConcreteUniversalConfigurationProvider struct{}

func (p ConcreteUniversalConfigurationProvider) GetUniversalDefaultConfig() MachineConfiguration {
	return GetUniversalDefaultConfig()
}

// ============================== //
// PRIVATE AND DEPRECATED METHODS //
// ============================== //

// @deprecated
func GetIntendedConfig(client dynamodb.GetItemAPI, machineID string) (intendedConfig MachineConfiguration, err error) {
	config, err := GetMachineSpecificConfig(client, machineID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get machine config")
		return
	}

	if config == nil {
		log.Printf("no config items found for %s", machineID)
		log.Printf("%s", "Retrieving global config")

		config, err = GetGlobalConfig(client)
		if err != nil {
			err = errors.Wrapf(err, "failed to get global config")
			return
		}
	}

	// Check if the result.Item is not nil. If this is nil, no valid global config exists within the DynamoDB and return a generic hardcoded value
	if config == nil {
		log.Printf("WARNING: GLOBAL CONFIG ITEM NOT FOUND; DEFAULTING TO HARDCODED DEFAULT CONFIG.")

		intendedConfig = GetUniversalDefaultConfig()
	} else {
		intendedConfig = *config
	}

	return
}

// @deprecated
func GetMachineSpecificConfig(client dynamodb.GetItemAPI, machineID string) (config *MachineConfiguration, err error) {
	return getItemAsMachineConfiguration(
		client,
		machineConfigurationPK(machineID),
		machineConfigurationSK(),
	)
}

// @deprecated
func GetGlobalConfig(client dynamodb.GetItemAPI) (config *MachineConfiguration, err error) {
	return getItemAsMachineConfiguration(
		client,
		globalConfigurationPK,
		machineConfigurationSK(),
	)
}

func getItemAsMachineConfiguration(client dynamodb.GetItemAPI, partitionKey string, sortKey string) (config *MachineConfiguration, err error) {
	output, err := client.GetItem(
		dynamodb.PrimaryKey{
			PartitionKey: partitionKey,
			SortKey:      sortKey,
		},
		false,
	)

	if err != nil {
		return
	}

	if len(output.Item) == 0 {
		return
	}

	err = attributevalue.UnmarshalMap(output.Item, &config)

	if err != nil {
		err = errors.Wrap(err, "succeeded GetItem but failed to unmarshalMap into MachineConfiguration")
		return
	}

	return
}
