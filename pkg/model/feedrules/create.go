package feedrules

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/rules"
)

func ConstructFeedRuleFromBaseRule(timeProvider clock.TimeProvider, rule rules.SantaRule) FeedRuleRow {
	return FeedRuleRow{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: feedRulesPK,
			// With this sort key, all rules will be ordered by the date they are created,
			// but also uniqified by the specific binary. This means that you can seek all rules
			// over time, kind of like picking up diffs through VCS changes.
			SortKey: fmt.Sprintf(
				"%s#%s",
				clock.RFC3339(timeProvider.Now()),
				rules.RuleSortKeyFromTypeSHA(rule.SHA256, rule.RuleType),
			),
		},
		SantaRule:    rule,
		ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
		DataType:     GetDataType(),
	}
}

func ReconstructFeedSyncLastEvaluatedKeyFromDate(feedSyncCursor string) dynamodb.PrimaryKey {
	return dynamodb.PrimaryKey{
		PartitionKey: feedRulesPK,
		SortKey:      feedSyncCursor,
	}
}
