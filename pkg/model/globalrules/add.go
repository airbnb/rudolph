package globalrules

import (
	"errors"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/feedrules"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func AddNewGlobalRule(
	time clock.TimeProvider,
	client dynamodb.TransactWriteItemsAPI,
	identifier string,
	ruleType types.RuleType,
	policy types.Policy,
	description string,
) error {
	// Input Validation
	isValid, err := ruleValidation(ruleType, policy)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("no errors occurred during the rule validation check but the provided rule is not valid")
	}

	rule := GlobalRuleRow{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: globalRulesPK,
			SortKey:      globalRulesSK(identifier, ruleType),
		},
		Description: description,
		SantaRule: rules.SantaRule{
			RuleType:   ruleType,
			Policy:     policy,
			Identifier: identifier,
		},
	}

	feedRule := feedrules.ConstructFeedRuleFromBaseRule(time, rule.SantaRule)

	putItem1, err := client.CreateTransactPutItem(rule)
	if err != nil {
		return err
	}
	putItem2, err := client.CreateTransactPutItem(feedRule)
	if err != nil {
		return err
	}

	putItems := []awsdynamodbtypes.TransactWriteItem{
		*putItem1,
		*putItem2,
	}

	_, err = client.TransactWriteItems(putItems, nil)
	return err
}

func ruleValidation(
	ruleType types.RuleType,
	policy types.Policy,
) (bool, error) {
	// RuleType validation
	_, err := ruleType.MarshalText()
	if err != nil {
		return false, err
	}

	// RulePolicy validation
	_, err = policy.MarshalText()
	if err != nil {
		return false, err
	}

	// All validations have passed
	return true, nil
}
