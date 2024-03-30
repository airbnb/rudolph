package globalrules

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/feedrules"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func UpdateGlobalRule(
	time clock.TimeProvider,
	client dynamodb.TransactWriteItemsAPI,
	identifier string,
	ruleType types.RuleType,
	rulePolicy types.Policy,
) error {
	// Get the PK/SK values
	pk := globalRulesPK
	sk := globalRulesSK(identifier, ruleType)

	primaryKey := dynamodb.PrimaryKey{
		PartitionKey: pk,
		SortKey:      sk,
	}

	// Updated rulePolicy request
	updateItem := updateRulePolicyRequest{
		Policy: rulePolicy,
	}

	// UpdatedRule for the ruleFeed Update
	updatedRule := rules.SantaRule{
		RuleType:   ruleType,
		Policy:     rulePolicy,
		Identifier: identifier,
	}

	updateFeedRuleItem := feedrules.ConstructFeedRuleFromBaseRule(time, updatedRule)

	// Create the Update the rule by creating a TransactUpdateItem
	updateItem1, err := client.CreateTransactUpdateItem(primaryKey, updateItem)
	if err != nil {
		return err
	}

	// Update the ruleFeed by creating a TransactPutItem
	putItem1, err := client.CreateTransactPutItem(updateFeedRuleItem)
	if err != nil {
		return err
	}

	// Build the transactWriteItems
	transactItems := []awsdynamodbtypes.TransactWriteItem{
		*updateItem1,
		*putItem1,
	}

	// Send the TransactWriteRequest
	_, err = client.TransactWriteItems(transactItems, nil)
	if err != nil {
		return fmt.Errorf("failed to update global rule: %w", err)
	}

	return nil
}
