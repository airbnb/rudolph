package ruledownload

import (
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

// RuledownloadResponse is the response body returned by /ruledownload endpoints
type RuledownloadResponse struct {
	Rules []RuledownloadRule `json:"rules"`
	// When a cursor is returned by the server, it is an indicator to the Santa sensor that there are
	// additional rules to be paginated through. This cursor is passed to the next request.
	Cursor *ruledownloadCursor `json:"cursor,omitempty"`
}

// RuledownloadRule is a single rule returned in a RuledownloadResponse
// It duck-types to/from the SantaRule struct type
// Documentation: https://santa.dev/development/sync-protocol.html#rules-objects
type RuledownloadRule struct {
	RuleType      types.RuleType `json:"rule_type"`
	Policy        types.Policy   `json:"policy"`
	SHA256        string         `json:"sha256,omitempty"`
	Identifier    string         `json:"identifier"`
	CustomMessage string         `json:"custom_msg,omitempty"`
}

// DDBRulesToResponseRules type converts the DynamoDB representation of a rule to an API
// representation of a Rule, which is returned in an API response.
func DDBRulesToResponseRules(rulesList []rules.SantaRule) (responseRules []RuledownloadRule) {
	responseRules = make([]RuledownloadRule, len(rulesList))

	for i, rule := range rulesList {
		responseRules[i] = RuledownloadRule(rule)
		// responseRules[i] = RuledownloadRule{
		// 	RuleType:      rule.RuleType,
		//     Policy:        rule.Policy,
		//     Identifier:    rule.Identifier,
		//     CustomMessage: rule.CustomMessage,
		// }
	}
	return
}
