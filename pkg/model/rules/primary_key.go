package rules

import (
	"fmt"
	"log"

	"github.com/airbnb/rudolph/pkg/types"
)

// @deprecated
func RuleSortKeyFromTypeSHA(sha256 string, ruleType types.RuleType) string {
	return RuleSortKeyFromTypeIdentifier(sha256, ruleType)
}

func RuleSortKeyFromTypeIdentifier(identifier string, ruleType types.RuleType) string {
	switch ruleType {
	case types.RuleTypeBinary:
		return fmt.Sprintf("%s%s", binaryRuleSKPrefix, identifier)
	case types.RuleTypeCertificate:
		return fmt.Sprintf("%s%s", certificateRuleSKPrefix, identifier)
	case types.RuleTypeTeamID:
		return fmt.Sprintf("%s%s", teamIDRuleSKPrefix, identifier)
	case types.RuleTypeSigningID:
		return fmt.Sprintf("%s%s", signingIDRuleSKPrefix, identifier)
	default:
		log.Printf("error (recovered): encountered unknown ruleType: (%+v)", ruleType)
		return ""
	}
}
