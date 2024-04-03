package globalrules

import (
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

func (g *GlobalRuleRow) globalRuleValidation() (bool, error) {
	// RuleType validation
	_, err := g.RuleType.MarshalText()
	if err != nil {
		return false, err
	}

	// RulePolicy validation
	_, err = g.Policy.MarshalText()
	if err != nil {
		return false, err
	}

	var validRuleIdentifier bool
	switch g.RuleType {
	case types.RuleTypeBinary:
		fallthrough
	case types.RuleTypeCertificate:
		validRuleIdentifier = rules.ValidSha256(g.Identifier)
	case types.RuleTypeTeamID:
		validRuleIdentifier = rules.ValidTeamID(g.Identifier)
	case types.RuleTypeSigningID:
		validRuleIdentifier = rules.ValidSigningID(g.Identifier)
	}

	if !validRuleIdentifier {
		return false, nil
	}

	// All validations have passed
	return true, nil
}
