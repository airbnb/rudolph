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
	RuleType      types.RuleType `dynamodbav:"RuleType"`
	Policy        types.Policy   `dynamodbav:"Policy"`
	SHA256        string         `dynamodbav:"SHA256"`
	CustomMessage string         `dynamodbav:"CustomMessage,omitempty"`
}
