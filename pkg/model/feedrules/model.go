package feedrules

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

const (
	feedRulesPK                     = "RulesFeed"
	feedRulesExpiresAfterInDays int = 90
)

type FeedRuleRow struct {
	dynamodb.PrimaryKey
	rules.SantaRule
	ExpiresAfter int64          `dynamodbav:"ExpiresAfter,omitempty"`
	DataType     types.DataType `dynamodbav:"DataType"`
}

func GetSyncStateExpiresAfter(timeProvider clock.TimeProvider) int64 {
	return clock.Unixtimestamp(timeProvider.Now().UTC().AddDate(0, 0, feedRulesExpiresAfterInDays))
}

func GetDataType() types.DataType {
	return types.DataTypeRulesFeed
}
