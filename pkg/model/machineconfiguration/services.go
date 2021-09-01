package machineconfiguration

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
)

//
// This service exposes all machineconfiguration access methods
//
type MachineConfigurationService interface {
	GetIntendedConfig(machineID string) (intendedConfig MachineConfiguration, err error)
	GetIntendedGlobalConfig() (intendedConfig MachineConfiguration, isDefaultConfig bool, err error)
	SetGlobalConfig(config MachineConfiguration) error
	SetMachineConfig(machineID string, config MachineConfiguration) error
	UpdateGlobalConfig(configRequest MachineConfigurationUpdateRequest) (*MachineConfiguration, error)
	UpdateMachineConfig(machineID string, configRequest MachineConfigurationUpdateRequest) (*MachineConfiguration, error)
	DeleteGlobalConfig() error
	DeleteMachineConfig(machineID string) error
}

type ConcreteMachineConfigurationService struct {
	dynamodb dynamodb.DynamoDBClient
	fetcher  ConcreteConfigurationFetcher
	setter   ConcreteConfigurationSetter
	updater  ConcreteConfigurationUpdater
	deleter  ConcreteConfigurationDeleter
}

func GetMachineConfigurationService(dynamodb dynamodb.DynamoDBClient, timeProvider clock.TimeProvider) MachineConfigurationService {
	return ConcreteMachineConfigurationService{
		dynamodb: dynamodb,
		fetcher:  GetConfigurationFetcher(dynamodb, timeProvider),
		setter:   GetConfigurationSetter(dynamodb, timeProvider),
		updater:  GetConfigurationUpdater(dynamodb, GetConfigurationFetcher(dynamodb, timeProvider), timeProvider),
		deleter:  GetConfigurationDeleter(dynamodb),
	}
}

func GetUncachedMachineConfigurationService(dynamodb dynamodb.DynamoDBClient, timeProvider clock.TimeProvider) MachineConfigurationService {
	return ConcreteMachineConfigurationService{
		dynamodb: dynamodb,
		fetcher:  GetUncachedConfigurationFetcher(dynamodb, timeProvider),
		setter:   GetConfigurationSetter(dynamodb, timeProvider),
		updater:  GetConfigurationUpdater(dynamodb, GetConfigurationFetcher(dynamodb, timeProvider), timeProvider),
		deleter:  GetConfigurationDeleter(dynamodb),
	}
}

// Getter //

// Machine
func (s ConcreteMachineConfigurationService) GetIntendedConfig(machineID string) (intendedConfig MachineConfiguration, err error) {
	return s.fetcher.getIntendedConfig(machineID)
}

// Global
func (s ConcreteMachineConfigurationService) GetIntendedGlobalConfig() (intendedConfig MachineConfiguration, isDefaultConfig bool, err error) {
	return s.fetcher.getIntendedGlobalConfig()
}

// Setter //

// Global
func (s ConcreteMachineConfigurationService) SetGlobalConfig(config MachineConfiguration) error {
	return s.setter.global.setGlobalConfig(config)
}

// Machine
func (s ConcreteMachineConfigurationService) SetMachineConfig(machineID string, config MachineConfiguration) error {
	return s.setter.machine.setMachineConfig(machineID, config)
}

// Updater //

// Global
func (s ConcreteMachineConfigurationService) UpdateGlobalConfig(configRequest MachineConfigurationUpdateRequest) (*MachineConfiguration, error) {
	return s.updater.global.updateConfig(configRequest)
}

// Machine

func (s ConcreteMachineConfigurationService) UpdateMachineConfig(machineID string, configRequest MachineConfigurationUpdateRequest) (*MachineConfiguration, error) {
	return s.updater.machine.updateConfig(machineID, configRequest)
}

// Deleter //

// Global
func (s ConcreteMachineConfigurationService) DeleteGlobalConfig() error {
	return s.deleter.global.deleteGlobalConfig()
}

// Machine
func (s ConcreteMachineConfigurationService) DeleteMachineConfig(machineID string) error {
	return s.deleter.machine.deleteMachineConfig(machineID)
}
