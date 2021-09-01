package syncstate

import (
	"time"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/pkg/errors"
)

func UpdatePostflightDate(timeProvider clock.TimeProvider, client dynamodb.UpdateItemAPI, machineID string) (err error) {
	_, err = client.UpdateItem(
		dynamodb.PrimaryKey{
			PartitionKey: syncStatePK(machineID),
			SortKey:      syncStateSK,
		},
		UpdatePostflightItem{
			PostflightAt: clock.RFC3339(timeProvider.Now()),
			// FIXME (derek.wang) ok, hear me out here--
			// Machines sync typically every 10 minutes. Due to DDB's eventual consistency model, there's a small chance for the
			// following to occur:
			//
			// - t=1 Old sensor cursor in database
			// - t=10 new feed rule added
			// - t=11 sensor sync starts
			// - t=12 sensor uses feed cursor "1" to download
			// - t=13 dynamodb:Query scans t>=1 and accidentally misses the new feed rule due to eventual consistency
			// - t=14 ruledownlod returns 0 rules
			// - t=15 postflight saves "15" into the database as the new cursor
			//
			// To mitigate against this, we always save a cursor that's 10 minutes EARLIER than the end of the sync.
			// This forces sensor sync processes to overscan a little. In practice, this will have very minimal impact to
			// performance, as all dynamodb:Query calls consume 0.5 RCUs + more RCUs based upon the number of items returned.
			// Redundantly scanning ranges of sort keys that are almost always empty will rarely increase cost.
			//
			// For FUTURE Derek to figure out:
			// The absolute BEST way of doing this is to only move cursor whenever ruledownload returns rules or has a
			// lastEvaluatedKey. Doing this would require additional writes on /ruledownload which would inofitself increase
			// WCUs which are way more expensive than RCUs. So yeah, my lame excuse of a spaghetti code will stand here.
			FeedSyncCursor: clock.RFC3339(timeProvider.Now().Add(time.Minute * -10)),
			ExpiresAfter:   GetSyncStateExpiresAfter(timeProvider),
		},
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to update item")
		return
	}
	return
}

func UpdateRuledownloadStartedAt(timeProvider clock.TimeProvider, client dynamodb.UpdateItemAPI, machineID string) (err error) {
	_, err = client.UpdateItem(
		dynamodb.PrimaryKey{
			PartitionKey: syncStatePK(machineID),
			SortKey:      syncStateSK,
		},
		updateRuledownloadAtItem{
			RuledownloadStartedAt: clock.RFC3339(timeProvider.Now()),
		},
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to update item")
		return
	}
	return
}

func UpdateRuledownloadFinishedAt(timeProvider clock.TimeProvider, client dynamodb.UpdateItemAPI, machineID string) (err error) {
	_, err = client.UpdateItem(
		dynamodb.PrimaryKey{
			PartitionKey: syncStatePK(machineID),
			SortKey:      syncStateSK,
		},
		updateRuledownloadFinishedAtItem{
			RuledownloadFinishedAt: clock.RFC3339(timeProvider.Now()),
		},
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to update item")
		return
	}
	return
}
