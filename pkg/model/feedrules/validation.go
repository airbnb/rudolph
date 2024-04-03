package feedrules

import (
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

func (f *FeedRuleRow) feedRuleRowValidation() (bool, error) {
	// RuleType validation
	_, err := f.RuleType.MarshalText()
	if err != nil {
		return false, err
	}

	// RulePolicy validation
	_, err = f.Policy.MarshalText()
	if err != nil {
		return false, err
	}

	var validRuleIdentifier bool
	switch f.RuleType {
	case types.RuleTypeBinary:
		fallthrough
	case types.RuleTypeCertificate:
		validRuleIdentifier = rules.ValidSha256(f.Identifier)
	case types.RuleTypeTeamID:
		validRuleIdentifier = rules.ValidTeamID(f.Identifier)
	case types.RuleTypeSigningID:
		validRuleIdentifier = rules.ValidSigningID(f.Identifier)
	}

	if !validRuleIdentifier {
		return false, nil
	}

	// All validations have passed
	return true, nil
}
