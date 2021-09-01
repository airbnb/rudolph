package globalrules

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
)

type GlobalRulesUpdater interface {
	UpdateGlobalRule(sha256 string, ruleType types.RuleType, rulePolicy types.Policy) error
}

type ConcreteGlobalRulesUpdater struct {
	ClockProvider clock.TimeProvider
	TransactWrite dynamodb.TransactWriteItemsAPI
}

func (c ConcreteGlobalRulesUpdater) UpdateGlobalRule(sha256 string, ruleType types.RuleType, rulePolicy types.Policy) (err error) {
	return UpdateGlobalRule(c.ClockProvider, c.TransactWrite, sha256, ruleType, rulePolicy)
}
