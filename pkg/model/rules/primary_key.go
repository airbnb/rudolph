package rules

import (
	"fmt"
	"log"

	"github.com/airbnb/rudolph/pkg/types"
)

func RuleSortKeyFromTypeSHA(sha256 string, ruleType types.RuleType) string {
	if len(sha256) != 64 {
		log.Printf("error (recovered): invalid sha256: (%s)", sha256)
		return ""
	}

	switch ruleType {
	case types.RuleTypeBinary:
		return fmt.Sprintf("%s%s", binaryRuleSKPrefix, sha256)
	case types.RuleTypeCertificate:
		return fmt.Sprintf("%s%s", certificateRuleSKPrefix, sha256)
	default:
		log.Printf("error (recovered): encountered unknown ruleType: (%+v)", ruleType)
		return ""
	}
}
