package globalrules

import (
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/feedrules"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

func AddNewGlobalRule(time clock.TimeProvider, client dynamodb.TransactWriteItemsAPI, sha256 string, ruleType types.RuleType, policy types.Policy, description string) error {
	// Input Validation
	isValid, err := inputValidation(sha256, ruleType, policy, description)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("no errors occurred during the rule validation check but the provided rule is not valid")
	}

	rule := GlobalRuleRow{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: globalRulesPK,
			SortKey:      globalRulesSK(sha256, ruleType),
		},
		Description: description,
		SantaRule: rules.SantaRule{
			RuleType: ruleType,
			Policy:   policy,
			SHA256:   sha256,
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

func inputValidation(sha256 string, ruleType types.RuleType, policy types.Policy, description string) (bool, error) {
	var err error

	// RuleSha256 validation
	err = types.ValidateSha256(sha256)
	if err != nil {
		return false, err
	}

	// RuleType validation
	_, err = ruleType.MarshalText()
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
