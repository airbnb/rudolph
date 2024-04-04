package preflight

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
)

// Select a MOD that will enable a proper dithering technique such as a proven cyclic group
// https://mathworld.wolfram.com/ModuloMultiplicationGroup.html
var MOD float64 = 11

const (
	daysToElapseUntilRefreshCleanSync = 7
)

// Steps to determine if a Clean Sync must be forced upon a requesting machine
// 1. Determined via the rules counts, ie if the returned number of rules from the machine equals zero, force a clean sync
//    Note: the DB may also have zero rules but this is fine then as the resulting sync time is the same
// 2. Periodically refresh systems via a clean sync --> this is determined by the amount of days that elapse since a machines last clean sync
//    Force this after 1d10 + 7 days.

func (c concreteCleanSyncService) determineCleanSync(machineID string, preflightRequest *PreflightRequest, syncState *syncstate.SyncStateRow) (performCleanSync bool, err error) {
	// Determine clean sync via rule count
	performCleanSync = determineCleanSyncByRuleCount(preflightRequest)
	if performCleanSync {
		return
	}

	// Periodically re-force a clean sync to keep rules fresh
	performCleanSync, err = determineCleanSyncRefresh(
		c.timeProvider,
		machineID,
		syncState,
	)
	return

}

func determineCleanSyncByRuleCount(preflightRequest *PreflightRequest) bool {
	ruleCount := preflightRequest.BinaryRuleCount + preflightRequest.CertificateRuleCount + preflightRequest.CompilerRuleCount + preflightRequest.TransitiveRuleCount
	return ruleCount == 0
}

func determineCleanSyncRefresh(timeProvider clock.TimeProvider, machineID string, syncState *syncstate.SyncStateRow) (performCleanSync bool, err error) {
	daysSinceLastCleanSync := daysSinceLastCleanSync(timeProvider, syncState)

	// To reduce stampeding, we introduce a bit of dithering by using the machineID as the seed to randomize the number of days required to elapse before performing a clean sync.
	// Given the same MachineID, performCleanSync will always provide the same chaos int and evenly space all clients out to require clean sync
	// 7 days + 1d10 (Based on the MachineID input) * 10 minutes
	performCleanSync, err = determinePeriodicRefreshCleanSync(
		machineID,
		daysSinceLastCleanSync,
	)
	return
}

// Returns the number of full 24 hour days since the previous clean sync
// Returns an unusually large number (99999999) when it thinks no sync has successfully happened before.
func daysSinceLastCleanSync(timeProvider clock.TimeProvider, prevSyncState *syncstate.SyncStateRow) int {
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

	currentTime := timeProvider.Now().UTC()
	diff := currentTime.Sub(lastCleanSyncTime)

	return int(diff.Hours() / 24)

}

// shouldPerformCleanSync is a method to perform dithering in an attempt to prevent a surge of forced clean syncs at the same time
// For a given MachineID, the chaos is equal to 1d10 and will remain consistent
func determinePeriodicRefreshCleanSync(machineID string, daysSinceLastSync int) (performCleanSync bool, err error) {
	// Generate a chaos by converting the machineID to an int
	chaos, err := machineIDToInt(machineID)
	if err != nil {
		return false, err
	}

	// Check to make sure that chaos is within 0 to 10
	if chaos > 11 || chaos < 0 {
		err = errors.New("chaos was greater than 11 or less than 0")
		return
	}

	if daysSinceLastSync >= (daysToElapseUntilRefreshCleanSync + chaos) {
		performCleanSync = true
	}
	return
}

// machineIDToInt outputs a consistent integer as provided by a consistent machineID string
// For a given MachineID, the chaos is equal to 1d10 or has modified by MOD variable
func machineIDToInt(machineID string) (int, error) {
	var result int
	machineIDHash := sha256.New()
	_, err := machineIDHash.Write([]byte(machineID))
	if err != nil {
		return result, fmt.Errorf("error generating hash from machineID: %w", err)
	}

	// Convert machineIDHash into a hex which will be converted into a *int
	bigInt := new(big.Int)
	bigInt.SetString(hex.EncodeToString(machineIDHash.Sum(nil)), 16)
	//machineIDFloat is required to perform mod math so convert the big int into a float
	machineIDFloat, err := strconv.ParseFloat(bigInt.String(), 64)
	if err != nil {
		return result, fmt.Errorf("error parsing machineID into a float: %w", err)
	}
	// chaos number is 1d10 based on the machineID to remain consistent
	result = int(math.Mod(machineIDFloat, MOD))
	return result, nil
}
