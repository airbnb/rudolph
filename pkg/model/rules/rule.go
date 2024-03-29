package rules

import (
	"github.com/airbnb/rudolph/pkg/types"
)

// SantaRule is the trimmed down concept of a rule
// These fields are commonly shared across different sub-types of rules:
//
// - global rules:   Source of truth for rules that are intended to be distributed to all endpoints
// - feed rules:     A timeline-based feed of rule "diffs" that endpoints can download
// - machine rules:  Rules intended only to be deployed to specific endpoints
type SantaRule struct {
	RuleType      types.RuleType `dynamodbav:"RuleType" json:"rule_type"`
	Policy        types.Policy   `dynamodbav:"Policy" json:"policy"`
	SHA256        string         `dynamodbav:"SHA256,omitempty" json:"sha256,omitempty"` // @deprecated - Use Identifier instead
	Identifier    string         `dynamodbav:"Identifier" json:"identifier"`
	CustomMessage string         `dynamodbav:"CustomMessage,omitempty" json:"custom_msg,omitempty"`
}
