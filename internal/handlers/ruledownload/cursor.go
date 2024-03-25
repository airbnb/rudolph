package ruledownload

import (
	"fmt"
	"log"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/feedrules"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
)

// ruledownloadCursor is passed between the server and client between successive API calls to /ruledownload
// This object not only represents a pagination cursor, but is also a way for the server and client to exchange
// metadata about the sync process without having to resort to serverside storage.
type ruledownloadCursor struct {
	// The Strategy determines HOW the server handles the
	Strategy ruledownloadStrategy `json:"strategy,omitempty"`
	ruledownloadCursorDDBLastEvaluatedKey
	PageNumber int `json:"page,omitempty"`
	BatchSize  int `json:"batch_size,omitempty"`
}

type ruledownloadCursorDDBLastEvaluatedKey struct {
	PartitionKey string `json:"pk,omitempty"`
	SortKey      string `json:"sk,omitempty"`
}

func (r *ruledownloadCursor) SetDynamodbLastEvaluatedKey(pk *dynamodb.PrimaryKey) {
	r.ruledownloadCursorDDBLastEvaluatedKey = ruledownloadCursorDDBLastEvaluatedKey(*pk)
}

// ToDynamoDBPrimaryKey converts the cursor into a DDB Primary Key
// When the cursor does not actually contain any values for the partition or sort key,
// returns nil to indicate that it does not have a last evaluated key
func (r ruledownloadCursor) GetLastEvaluatedKey() *dynamodb.PrimaryKey {
	if r.PartitionKey != "" || r.SortKey != "" {
		return (*dynamodb.PrimaryKey)(&r.ruledownloadCursorDDBLastEvaluatedKey)
	}
	return nil
}

func (r ruledownloadCursor) CloneForNextPage() ruledownloadCursor {
	return ruledownloadCursor{
		Strategy:   r.Strategy,
		BatchSize:  r.BatchSize,
		PageNumber: r.PageNumber + 1,
	}
}

func (r *ruledownloadCursor) SetStrategy(strategy ruledownloadStrategy) {
	r.Strategy = strategy
}

type ruledownloadStrategy int

const (
	// The clean download strategy downloads from the GlobalRules
	// The LastEvaluatedKey references the sort key where pk = "GlobalRules"
	// All clean downloads will paginate over all global rules
	ruledownloadStrategyClean ruledownloadStrategy = iota + 1

	// The incremental download strategy downloads from the rules feed
	// The LastEvaluatedKey references the sort key where pk = "FeedRules"
	// Incremental downloads will read the last incremental download cursor from the previous
	// sync state in the database and paginate from there, never re-reading rules that have already
	// been read from the feed.
	ruledownloadStrategyIncremental

	// The machine download strategy
	// Currently the LastEvaluatedKey is ignored as this strategy will always download
	// the entire set of machine-specific rules, regardless of how large it is.
	ruledownloadStrategyMachine
)

type ruledownloadCursorService interface {
	ConstructCursor(ruledownloadRequest RuledownloadRequest, machineID string) (cursor ruledownloadCursor, err error)
}

type concreteRuledownloadCursorService struct {
	timer   clock.TimeProvider
	updater dynamodb.UpdateItemAPI
	getter  dynamodb.GetItemAPI
}

// ConstructCursor takes the given ruledownload request and constructs a cursor object. This cursor is sent to
// the ruledownload services for processing.
//
// In general, if a cursor already exists in the request, it will be accepted. Else, the cursor is generated anew,
// depending on what the sync state is in the database.
func (c concreteRuledownloadCursorService) ConstructCursor(ruledownloadRequest RuledownloadRequest, machineID string) (cursor ruledownloadCursor, err error) {
	if ruledownloadRequest.Cursor == nil {
		// With a nil cursor, this is the very first ruledownload request in a sync.
		// Couple things we need to do:
		//  * note down that ruledownload has started
		//  * determine the ruledownload strategy
		//  * embed any context, including the desired strategy, into the next cursor
		log.Printf("  Ruledownload first page")

		// Get the sync state as set by /preflight. The sync state will contain the final result of the previous
		syncState, eerr := syncstate.GetByMachineID(c.getter, machineID)
		if eerr != nil {
			err = fmt.Errorf("failed to get previous syncstate: %w", eerr)
			return
		}

		if syncState.CleanSync {
			// Clean syncs never start with a LastEvaluatedKey as it always begins from the beginning
			cursor = ruledownloadCursor{
				Strategy:   ruledownloadStrategyClean,
				BatchSize:  syncState.BatchSize,
				PageNumber: 1,
			}
		} else {
			// Clean syncs
			lastEvalKey := feedrules.ReconstructFeedSyncLastEvaluatedKeyFromDate(syncState.FeedSyncCursor)
			cursor = ruledownloadCursor{
				Strategy:                              ruledownloadStrategyIncremental,
				ruledownloadCursorDDBLastEvaluatedKey: ruledownloadCursorDDBLastEvaluatedKey(lastEvalKey),
				BatchSize:                             syncState.BatchSize,
				PageNumber:                            1,
			}
		}

		// Create a sensor sync object to log the StartedAt time of the rule download process
		err = syncstate.UpdateRuledownloadStartedAt(c.timer, c.updater, machineID)
		if err != nil {
			return
		}

	} else {
		cursor = *ruledownloadRequest.Cursor
	}
	return
}
