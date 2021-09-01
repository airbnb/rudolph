package machineconfiguration

import (
	"github.com/airbnb/rudolph/pkg/dynamodb"
)

func GetConfigurationDeleter(client dynamodb.DeleteItemAPI) ConcreteConfigurationDeleter {
	return ConcreteConfigurationDeleter{
		global: ConcreteGlobalConfigurationDeleter{
			deleter: client,
		},
		machine: ConcreteMachineConfigurationDeleter{
			deleter: client,
		},
	}
}

type ConcreteConfigurationDeleter struct {
	global  GlobalConfigurationDeleter
	machine MachineConfigurationDeleter
}

// GlobalConfigurationDeleter
type GlobalConfigurationDeleter interface {
	deleteGlobalConfig() error
}

type ConcreteGlobalConfigurationDeleter struct {
	deleter dynamodb.DeleteItemAPI
}

func (d ConcreteGlobalConfigurationDeleter) deleteGlobalConfig() error {
	globalPK := dynamodb.PrimaryKey{
		PartitionKey: globalConfigurationPK,
		SortKey:      machineConfigurationSK(),
	}
	_, err := d.deleter.DeleteItem(globalPK)
	return err
}

// GlobalConfigurationDeleter
type MachineConfigurationDeleter interface {
	deleteMachineConfig(machineID string) error
}

type ConcreteMachineConfigurationDeleter struct {
	deleter dynamodb.DeleteItemAPI
}

func (d ConcreteMachineConfigurationDeleter) deleteMachineConfig(machineID string) error {
	machinePK := dynamodb.PrimaryKey{
		PartitionKey: machineConfigurationPK(machineID),
		SortKey:      machineConfigurationSK(),
	}
	_, err := d.deleter.DeleteItem(machinePK)
	return err
}
