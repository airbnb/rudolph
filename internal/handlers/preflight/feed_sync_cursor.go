package preflight

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
)

func (c concreteStateTrackingService) getFeedSyncStateCursor(syncState *syncstate.SyncStateRow) (feedSyncCursor string, performCleanSync bool) {
	// Check when the last preflight request took place
	if syncState != nil && syncState.FeedSyncCursor != "" {
		// Inherit the feed feed sync cursor from the previous sync state to kind of "pick up where it left off"
		feedSyncCursor = syncState.FeedSyncCursor
	} else {
		// If there is no previous sync state, or if no cursor exists.
		// Always force it to clean sync and just set the feed sync cursor to "now"
		feedSyncCursor = clock.RFC3339(c.timeProvider.Now())
		performCleanSync = true
	}

	return
}
