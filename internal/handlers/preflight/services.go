package preflight

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"math/big"
	"strconv"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	"github.com/airbnb/rudolph/pkg/model/sensordata"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	"github.com/pkg/errors"
)

// Select a MOD that will enable a proper dithering technique such as a proven cyclic group
// https://mathworld.wolfram.com/ModuloMultiplicationGroup.html
var MOD float64 = 11

//
// SensorDataSaver
//
type sensorDataSaver interface {
	saveSensorDataFromRequest(time clock.TimeProvider, machineID string, request *PreflightRequest) error
}
type concreteSensorDataSaver struct {
	putter dynamodb.PutItemAPI
}

func (c concreteSensorDataSaver) saveSensorDataFromRequest(time clock.TimeProvider, machineID string, request *PreflightRequest) error {
	sensorData := sensordata.NewSensorData(
		time,
		machineID,
		request.SerialNumber,
		request.OSVersion,
		request.OSBuild,
		request.RequestCleanSync,
		request.PrimaryUser,
		request.CertificateRuleCount,
		request.BinaryRuleCount,
		request.CompilerRuleCount,
		request.TransitiveRuleCount,
	)
	_, err := c.putter.PutItem(sensorData)
	return err
}

//
// MachineConfigurationGetter
//
type machineConfigurationGetter interface {
	getDesiredConfig(machineID string) (config machineconfiguration.MachineConfiguration, err error)
}
type concreteMachineConfigurationGetter struct {
	fetcher machineconfiguration.MachineConfigurationService
}

func (c concreteMachineConfigurationGetter) getDesiredConfig(machineID string) (config machineconfiguration.MachineConfiguration, err error) {
	return c.fetcher.GetIntendedConfig(machineID)
}

//
// SyncStateGetter
//
type syncStateManager interface {
	getSyncState(machineID string) (syncState *syncstate.SyncStateRow, err error)
	saveNewSyncState(time clock.TimeProvider, machineID string, requestCleanSync bool, lastCleanSync string, batchSize int, feedSyncCursor string) error
}
type concreteSyncStateManager struct {
	getter dynamodb.GetItemAPI
	putter dynamodb.PutItemAPI
}

func (c concreteSyncStateManager) getSyncState(machineID string) (syncState *syncstate.SyncStateRow, err error) {
	return syncstate.GetByMachineID(c.getter, machineID)
}
func (c concreteSyncStateManager) saveNewSyncState(time clock.TimeProvider, machineID string, requestCleanSync bool, lastCleanSync string, batchSize int, feedSyncCursor string) error {
	syncState := syncstate.CreateNewSyncState(time, machineID, requestCleanSync, lastCleanSync, batchSize, feedSyncCursor)
	_, err := c.putter.PutItem(syncState)
	return err
}

// Returns the number of full 24 hour days since the previous clean sync
// Returns an unusually large number (99999999) when it thinks no sync has successfully happened before.
func daysSinceLastSync(time clock.TimeProvider, prevSyncState *syncstate.SyncStateRow) int {
	infinity := 99999999
	if prevSyncState == nil {
		return infinity
	}

	if prevSyncState.LastCleanSync == "" {
		return infinity
	}

	lastCleanSyncTime, err := clock.ParseRFC3339(prevSyncState.LastCleanSync)
	if err != nil {
		// Disregard the error
		log.Printf("failed to determine number of days since last sync from value: (%s); going to clean sync anyway", prevSyncState.LastCleanSync)
		return infinity
	}

	currentTime := time.Now().UTC()
	diff := currentTime.Sub(lastCleanSyncTime)

	return int(diff.Hours() / 24)
}

// shouldPerformCleanSync is a method to perform dithering in an attempt to prevent a surge of forced clean syncs at the same time
// For a given MachineID, the chaos is equal to 1d10 and will remain consistent
func shouldPerformCleanSync(machineID string, daysSinceLastSync int, daysElapseUntilCleanSync int) (bool, error) {
	// Generate a chaos by converting the machineID to an int
	chaos, err := machineIDToInt(machineID)
	if err != nil {
		return false, err
	}

	// Check to make sure that chaos is within 0 to 10
	if chaos > 11 || chaos < 0 {
		return false, errors.New("chaos was greater than 11 or less than 0")
	}
	// log.Printf("MachineID: %s | Days Since Last Sync: %d | Days Elapse Until Clean Sync: %d | chaos: %d", machineID, daysSinceLastSync, daysElapseUntilCleanSync, chaos)
	if daysSinceLastSync >= (daysElapseUntilCleanSync + chaos) {
		return true, nil
	} else {
		return false, nil
	}
}

// machineIDToInt outputs a consistent integer as provided by a consistent machineID string
// For a given MachineID, the chaos is equal to 1d10 or has modified by MOD variable
func machineIDToInt(machineID string) (int, error) {
	var result int
	machineIDHash := sha256.New()
	_, err := machineIDHash.Write([]byte(machineID))
	if err != nil {
		return result, errors.Wrap(err, "error generating hash from machineID")
	}

	// Convert machineIDHash into a hex which will be converted into a *int
	bigInt := new(big.Int)
	bigInt.SetString(hex.EncodeToString(machineIDHash.Sum(nil)), 16)
	//machineIDFloat is required to perform mod math so convert the big int into a float
	machineIDFloat, err := strconv.ParseFloat(bigInt.String(), 64)
	if err != nil {
		return result, errors.Wrap(err, "error parsing machineID into a float")
	}
	// chaos number is 1d10 based on the machineID to remain consistent
	result = int(math.Mod(machineIDFloat, MOD))
	return result, nil
}
