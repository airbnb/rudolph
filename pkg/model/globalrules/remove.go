package globalrules

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/feedrules"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// RemoveGlobalRule will remove the rule from the global repository of rules.
// It also creates a new rule entry in the feed that explicitly tells future syncs
// to remove the rule too.
func RemoveGlobalRule(timeProvider clock.TimeProvider, getter dynamodb.GetItemAPI, transacter dynamodb.TransactWriteItemsAPI, ruleSortKey string, txnIdempotencyKey string) (err error) {
	rule, err := GetGlobalRuleBySortKey(getter, ruleSortKey)
	if err != nil {
		return errors.Wrap(err, "query to retrieve existing rule failed")
	}
	if rule == nil {
		return errors.New(fmt.Sprintf("no such rule with sk (%s) exists", ruleSortKey))
	}

	// Delete the global rule
	txnDeleteItem, err := transacter.CreateTransactDeleteItem(rule.PrimaryKey)
	if err != nil {
		return errors.Wrap(err, "failed to create txn item to delete original rule")
	}

	// In order to get non-clean sync clients to pick up the new rule diff, add it to the feed as a "remove"
	// Subsequent non-clean syncs will use the feed cursor to cross the date
	feedrule := feedrules.ConstructFeedRuleFromBaseRule(timeProvider, rule.SantaRule)
	feedrule.Policy = types.Remove
	txnPutItem, err := transacter.CreateTransactPutItem(feedrule)
	if err != nil {
		return errors.Wrap(err, "failed to create txn item to add removal to feed")
	}

	txnItems := []awsdynamodbtypes.TransactWriteItem{
		*txnDeleteItem,
		*txnPutItem,
	}

	// DynamoDB Idempotency Keys may be 1-32 character length
	// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html#API_TransactWriteItems_RequestSyntax
	// If one is not provided, generate one
	if len(txnIdempotencyKey) < 0 || len(txnIdempotencyKey) > 32 {
		txnIdempotencyKey = uuid.NewString()
	}
	_, err = transacter.TransactWriteItems(txnItems, &txnIdempotencyKey)
	if err != nil {
		return errors.Wrap(err, "transaction delete failed")
	}

	return
}

type RuleRemovalService interface {
	RemoveGlobalRule(ruleSortKey string, idempotencyKey string) (err error)
}
type ConcreteRuleRemovalService struct {
	TimeProvider clock.TimeProvider
	Getter       dynamodb.GetItemAPI
	Transacter   dynamodb.TransactWriteItemsAPI
}

func (c ConcreteRuleRemovalService) RemoveGlobalRule(ruleSortKey string, idempotencyKey string) (err error) {
	return RemoveGlobalRule(c.TimeProvider, c.Getter, c.Transacter, ruleSortKey, idempotencyKey)
}
