package machinerules

import (
	"errors"
	"time"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

func AddNewMachineRule(
	client dynamodb.PutItemAPI,
	machineID string,
	identifier string,
	ruleType types.RuleType,
	policy types.Policy,
	description string,
	expires time.Time,
) error {
	// Input Validation
	isValid, err := ruleValidation(
		machineID,
		identifier,
		ruleType,
		policy,
		description,
		expires,
	)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("no errors occurred during the rule validation check but the provided rule is not valid")
	}

	rule := MachineRuleRow{
		PrimaryKey: dynamodb.PrimaryKey{
			PartitionKey: machineRulePK(machineID),
			SortKey:      machineRuleSK(identifier, ruleType),
		},
		Description: description,
		SantaRule: rules.SantaRule{
			RuleType:   ruleType,
			Policy:     policy,
			Identifier: identifier,
		},
		ExpiresAfter: expires.Unix(),
	}

	_, err = client.PutItem(rule)
	if err != nil {
		return err
	}
	return nil
}

func ruleValidation(
	machineID,
	identifier string,
	ruleType types.RuleType,
	policy types.Policy,
	description string,
	expires time.Time,
) (bool, error) {
	// MachineID validation
	err := types.ValidateMachineID(machineID)
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

	// Expires validation - must be a positive time period
	if expires.IsZero() {
		err = errors.New("expires time is not a positive time")
		return false, err
	}

	// All validations have passed
	return true, nil
}
