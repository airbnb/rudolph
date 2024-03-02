package machineconfiguration

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMachineConfigurationUpdater func(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error)

func (m mockMachineConfigurationUpdater) UpdateItem(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
	return m(key, item)
}

const (
	verbose bool = false
)

var (
	machineID                 = "858CBF28-5EAA-58A3-A155-BA5E81D5B5DD"
	expectedMachineClientMode = types.Lockdown
	expectedGlobalClientMode  = types.Monitor
	expectedGlobalSyncMode    = types.SyncTypeNormal

	expectedMachineConfig = MachineConfiguration{
		ClientMode:       expectedMachineClientMode,
		BatchSize:        10,
		AllowedPathRegex: "/usr/bin/local, /opt/bin",
		BlockedPathRegex: "/trash",
		SyncType:         expectedGlobalSyncMode,
	}

	expectedGlobalConfig = MachineConfiguration{
		ClientMode:       expectedGlobalClientMode,
		BatchSize:        20,
		AllowedPathRegex: "/usr/bin/local, /opt/bin",
		BlockedPathRegex: "/trash, /tmp",
		SyncType:         expectedGlobalSyncMode,
	}

	expectedGlobalConfigDefault = GetUniversalDefaultConfig()

	machineConfigurationRequestBatchSize    int    = 11
	machineConfigurationRequestAllowedPaths string = "/usr/bin/local, /opt/bin"
	machineConfigurationRequestBlockedPaths string = "/trash"

	machineConfigurationRequest = MachineConfigurationUpdateRequest{
		ClientMode:       &expectedMachineClientMode,
		BatchSize:        &machineConfigurationRequestBatchSize,
		AllowedPathRegex: &machineConfigurationRequestAllowedPaths,
		BlockedPathRegex: &machineConfigurationRequestBlockedPaths,
		SyncType:         &expectedGlobalSyncMode,
	}

	globalConfigurationRequestBatchSize    int    = 21
	globalConfigurationRequestAllowedPaths string = "/usr/bin/local, /opt/bin"
	globalConfigurationRequestBlockedPaths string = "/trash, /tmp"

	globalConfigurationRequest = MachineConfigurationUpdateRequest{
		ClientMode:       &expectedGlobalClientMode,
		BatchSize:        &globalConfigurationRequestBatchSize,
		AllowedPathRegex: &globalConfigurationRequestAllowedPaths,
		BlockedPathRegex: &globalConfigurationRequestBlockedPaths,
	}

	expectedGlobalConfigLockdown = MachineConfiguration{
		ClientMode:       types.Lockdown,
		BatchSize:        20,
		AllowedPathRegex: "/usr/bin/local, /opt/bin",
		BlockedPathRegex: "/trash, /tmp",
		SyncType:         types.SyncTypeNormal,
	}

	machineConfigPKey = dynamodb.PrimaryKey{
		PartitionKey: machineConfigurationPK(machineID),
		SortKey:      machineConfigurationSK(),
	}

	globalConfigPKey = dynamodb.PrimaryKey{
		PartitionKey: globalConfigurationPK,
		SortKey:      machineConfigurationSK(),
	}

	machineConfigRow = MachineConfigurationRow{
		machineConfigPKey,
		expectedMachineConfig,
	}

	globalConfigRow = MachineConfigurationRow{
		globalConfigPKey,
		expectedGlobalConfig,
	}

	globalConfigRowDefault = MachineConfigurationRow{
		globalConfigPKey,
		MachineConfiguration{
			ClientMode:             types.Monitor,
			BlockedPathRegex:       "",
			AllowedPathRegex:       "",
			BatchSize:              50,
			EnableBundles:          false,
			EnabledTransitiveRules: false,
			CleanSync:              false,
			FullSyncInterval:       DefaultFullSyncInterval,
			UploadLogsURL:          "",
			SyncType:               types.SyncTypeNormal,
			DataType:               types.DataTypeGlobalConfig,
		},
	}

	globalConfigRowLockdown = MachineConfigurationRow{
		globalConfigPKey,
		expectedGlobalConfig,
	}

	timeProvider = clock.ConcreteTimeProvider{}
)
var _ dynamodb.UpdateItemAPI = mockMachineConfigurationUpdater(nil)

type MockDynamodb struct {
	dynamodb.DynamoDBClient
	ConcreteConfigurationFetcher
	ConcreteConfigurationSetter
	ConcreteConfigurationUpdater
	ConcreteConfigurationDeleter
	mock.Mock
}

func (m *MockDynamodb) GetItem(key dynamodb.PrimaryKey, consistentRead bool) (*awsdynamodb.GetItemOutput, error) {
	args := m.Called(key, consistentRead)
	return args.Get(0).(*awsdynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamodb) PutItem(item interface{}) (*awsdynamodb.PutItemOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*awsdynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamodb) UpdateItem(key dynamodb.PrimaryKey, item interface{}) (*awsdynamodb.UpdateItemOutput, error) {
	args := m.Called(key, item)
	return args.Get(0).(*awsdynamodb.UpdateItemOutput), args.Error(1)
}

func (m *MockDynamodb) DeleteItem(key dynamodb.PrimaryKey) (*awsdynamodb.DeleteItemOutput, error) {
	args := m.Called(key)
	return args.Get(0).(*awsdynamodb.DeleteItemOutput), args.Error(1)
}

func Test_Service_Fetcher_Functions(t *testing.T) {
	t.Run("GetItem - GetIntendedConfig(machineID)", func(t *testing.T) {
		returnedItem, err := attributevalue.MarshalMap(machineConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		mocked := &MockDynamodb{}
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
			Item: returnedItem,
		}, nil)

		service := ConcreteMachineConfigurationService{
			dynamodb: mocked,
			fetcher:  GetConfigurationFetcher(mocked, timeProvider),
		}

		resultConfig, err := service.GetIntendedConfig(machineID)
		assert.Empty(t, err)

		assert.Equal(t, resultConfig.BatchSize, expectedMachineConfig.BatchSize)
		assert.Equal(t, resultConfig.ClientMode, expectedMachineConfig.ClientMode)
		assert.Equal(t, resultConfig.AllowedPathRegex, expectedMachineConfig.AllowedPathRegex)
		if verbose {
			t.Logf("\n%+v\n", resultConfig)
		}
	})
	t.Run("GetItem - GetIntendedGlobalConfig()", func(t *testing.T) {
		returnedItem, err := attributevalue.MarshalMap(globalConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		mocked := &MockDynamodb{}
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
			Item: returnedItem,
		}, nil)

		// service := ConcreteMachineConfigurationService{
		// 	dynamodb: mocked,
		// 	fetcher:  GetConfigurationFetcher(mocked, timeProvider),
		// }
		// Test the GetMachineConfigurationService functionality
		service := GetMachineConfigurationService(mocked, timeProvider)

		resultConfig, isDefaultConfig, err := service.GetIntendedGlobalConfig()
		assert.Empty(t, err)

		assert.Equal(t, isDefaultConfig, false)
		assert.Equal(t, resultConfig.BatchSize, expectedGlobalConfig.BatchSize)
		assert.Equal(t, resultConfig.ClientMode, expectedGlobalConfig.ClientMode)
		assert.Equal(t, resultConfig.AllowedPathRegex, expectedGlobalConfig.AllowedPathRegex)
		if verbose {
			t.Logf("\n%+v\n", resultConfig)
		}
	})
	t.Run("GetItem - GetIntendedGlobalConfig() - Matches Default - Universal Global Config", func(t *testing.T) {
		mocked := &MockDynamodb{}
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{}, nil)

		// Test the GetMachineConfigurationService functionality
		service := GetMachineConfigurationService(mocked, timeProvider)

		resultConfig, isDefaultConfig, err := service.GetIntendedGlobalConfig()
		assert.Empty(t, err)

		assert.Equal(t, isDefaultConfig, true)
		assert.Equal(t, expectedGlobalConfigDefault.BatchSize, resultConfig.BatchSize)
		assert.Equal(t, resultConfig.ClientMode, expectedGlobalConfigDefault.ClientMode)
		assert.Equal(t, resultConfig.AllowedPathRegex, expectedGlobalConfigDefault.AllowedPathRegex)
		if verbose {
			t.Logf("\n%+v\n", resultConfig)
		}
	})
	t.Run("GetItem - GetUncachedGlobalConfig() - Expect No Error, isDefaultConfig, and Config Matches Default", func(t *testing.T) {
		mocked := &MockDynamodb{}
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{}, nil)

		// Test the GetMachineConfigurationService functionality
		service := GetUncachedMachineConfigurationService(mocked, timeProvider)

		resultConfig, isDefaultConfig, err := service.GetIntendedGlobalConfig()
		assert.Empty(t, err)

		assert.True(t, isDefaultConfig)
		assert.Equal(t, resultConfig, globalConfigRowDefault.MachineConfiguration)
	})
	t.Run("GetItem - GetUncachedGlobalConfig() - Expect No Error but Empty Configuration", func(t *testing.T) {
		returnedItem, err := attributevalue.MarshalMap(globalConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		mocked := &MockDynamodb{}
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
			Item: returnedItem,
		}, nil)

		// service := ConcreteMachineConfigurationService{
		// 	dynamodb: mocked,
		// 	fetcher:  GetConfigurationFetcher(mocked, timeProvider),
		// }
		// Test the GetMachineConfigurationService functionality
		service := GetUncachedMachineConfigurationService(mocked, timeProvider)

		resultConfig, isDefaultConfig, err := service.GetIntendedGlobalConfig()
		assert.Empty(t, err)

		assert.Equal(t, isDefaultConfig, false)
		assert.Equal(t, resultConfig.BatchSize, expectedGlobalConfig.BatchSize)
		assert.Equal(t, resultConfig.ClientMode, expectedGlobalConfig.ClientMode)
		assert.Equal(t, resultConfig.AllowedPathRegex, expectedGlobalConfig.AllowedPathRegex)
		if verbose {
			t.Logf("\n%+v\n", resultConfig)
		}
	})
}

func Test_Service_Setter_Functions(t *testing.T) {
	t.Run("PutItem - SetGlobalConfig()", func(t *testing.T) {
		mocked := &MockDynamodb{}
		mocked.On("PutItem", mock.MatchedBy(func(item interface{}) bool {
			config := item.(*MachineConfigurationRow)
			return config.ClientMode == expectedGlobalConfig.ClientMode && config.AllowedPathRegex == expectedGlobalConfig.AllowedPathRegex && config.BatchSize == expectedGlobalConfig.BatchSize && config.BlockedPathRegex == expectedGlobalConfig.BlockedPathRegex
		})).Return(&awsdynamodb.PutItemOutput{}, nil)

		service := ConcreteMachineConfigurationService{
			dynamodb: mocked,
			setter:   GetConfigurationSetter(mocked, timeProvider),
		}

		err := service.SetGlobalConfig(expectedGlobalConfig)
		assert.Empty(t, err)
		mocked.AssertCalled(t, "PutItem", mock.Anything)
	})
	t.Run("PutItem - SetMachineConfig(machineID)", func(t *testing.T) {
		mocked := &MockDynamodb{}
		mocked.On("PutItem", mock.MatchedBy(func(item interface{}) bool {
			config := item.(*MachineConfigurationRow)
			return config.ClientMode == expectedMachineConfig.ClientMode && config.AllowedPathRegex == expectedMachineConfig.AllowedPathRegex && config.BatchSize == expectedMachineConfig.BatchSize && config.BlockedPathRegex == expectedMachineConfig.BlockedPathRegex
		})).Return(&awsdynamodb.PutItemOutput{}, nil)

		service := ConcreteMachineConfigurationService{
			dynamodb: mocked,
			setter:   GetConfigurationSetter(mocked, timeProvider),
		}

		err := service.SetMachineConfig(machineID, expectedMachineConfig)
		assert.Empty(t, err)
		mocked.AssertCalled(t, "PutItem", mock.Anything)
	})
}

func Test_Service_Updater_Global_Services(t *testing.T) {
	t.Run("UpdateItem - UpdateGlobalConfig()", func(t *testing.T) {
		modifiedGlobalConfigRow := globalConfigRow
		modifiedGlobalConfigRow.BatchSize = 21
		updatedItem, err := attributevalue.MarshalMap(modifiedGlobalConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		returnedItem, err := attributevalue.MarshalMap(globalConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		mocked := &MockDynamodb{}
		mocked.On("UpdateItem", mock.Anything, mock.Anything).Return(&awsdynamodb.UpdateItemOutput{
			Attributes: updatedItem,
		}, nil)
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
			Item: returnedItem,
		}, nil)

		service := GetMachineConfigurationService(mocked, timeProvider)

		// Modify expectedGlobalConfig.BatchSize to 21
		modifiedGlobalConfig := expectedGlobalConfig
		modifiedGlobalConfig.BatchSize = 21

		globalConfig, err := service.UpdateGlobalConfig(globalConfigurationRequest)
		assert.Empty(t, err)
		mocked.AssertCalled(t, "UpdateItem", mock.Anything, mock.Anything)
		mocked.AssertCalled(t, "GetItem", mock.Anything, mock.Anything)
		assert.NotEmpty(t, globalConfig)
		assert.Equal(t, globalConfig.BatchSize, modifiedGlobalConfig.BatchSize)
		assert.Equal(t, globalConfig.ClientMode, modifiedGlobalConfig.ClientMode)
		assert.Equal(t, globalConfig.AllowedPathRegex, modifiedGlobalConfig.AllowedPathRegex)
		if verbose {
			t.Logf("%+v", globalConfig)
		}
	})
}

func Test_Service_Updater_MachineConfiguration_Services(t *testing.T) {
	t.Run("UpdateItem - UpdateGlobalConfig()", func(t *testing.T) {
		modifiedMachineConfigRow := machineConfigRow
		modifiedMachineConfigRow.BatchSize = 11
		updatedItem, err := attributevalue.MarshalMap(modifiedMachineConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		returnedItem, err := attributevalue.MarshalMap(machineConfigRow)
		if err != nil {
			t.Fatal(err)
		}
		mocked := &MockDynamodb{}
		mocked.On("GetItem", mock.Anything, mock.Anything).Return(&awsdynamodb.GetItemOutput{
			Item: returnedItem,
		}, nil)
		mocked.On("UpdateItem", mock.Anything, mock.Anything).Return(&awsdynamodb.UpdateItemOutput{
			Attributes: updatedItem,
		}, nil)

		service := GetMachineConfigurationService(mocked, timeProvider)

		// Modify expectedGlobalConfig.BatchSize to 11
		modifiedMachineConfig := expectedMachineConfig
		modifiedMachineConfig.BatchSize = 11

		machineConfig, err := service.UpdateMachineConfig(machineID, machineConfigurationRequest)
		assert.Empty(t, err)
		mocked.AssertCalled(t, "UpdateItem", mock.Anything, mock.Anything)
		mocked.AssertCalled(t, "GetItem", mock.Anything, mock.Anything)
		assert.NotEmpty(t, machineConfig)
		assert.Equal(t, machineConfig.BatchSize, modifiedMachineConfig.BatchSize)
		assert.Equal(t, machineConfig.ClientMode, modifiedMachineConfig.ClientMode)
		assert.Equal(t, machineConfig.AllowedPathRegex, modifiedMachineConfig.AllowedPathRegex)
		if verbose {
			t.Logf("%+v", machineConfig)
		}
	})
}

func Test_Service_Deleter_MachineConfiguration_Services(t *testing.T) {
	t.Run("DeleteItem - DeleteMachineConfig(machineID)", func(t *testing.T) {
		mocked := &MockDynamodb{}
		mocked.On("DeleteItem", mock.Anything, mock.Anything).Return(&awsdynamodb.DeleteItemOutput{}, nil)

		service := GetMachineConfigurationService(mocked, timeProvider)

		err := service.DeleteMachineConfig(machineID)
		assert.Empty(t, err)
		mocked.AssertCalled(t, "DeleteItem", mock.Anything, mock.Anything)
	})
}

func Test_Service_Deleter_GlobalConfiguration_Services(t *testing.T) {
	t.Run("DeleteItem - DeleteGlobalConfig()", func(t *testing.T) {
		mocked := &MockDynamodb{}
		mocked.On("DeleteItem", mock.Anything, mock.Anything).Return(&awsdynamodb.DeleteItemOutput{}, nil)

		service := GetMachineConfigurationService(mocked, timeProvider)

		err := service.DeleteGlobalConfig()
		assert.Empty(t, err)
		mocked.AssertCalled(t, "DeleteItem", mock.Anything, mock.Anything)
	})
}
