package machinerules

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

const MachineRuleDefaultExpirationHours = 24
const machineRulesPKPrefix = "MachineRules#"

type MachineRuleRow struct {
	dynamodb.PrimaryKey
	rules.SantaRule
	Description      string `dynamodbav:"Description,omitempty"`
	DeleteOnNextSync bool   `dynamodbav:"DeleteOnNextSync,omitempty"`
	ExpiresAfter     int64  `dynamodbav:"ExpiresAfter,omitempty"`
	MachineID        string `dynamodbav:"MachineID,omitempty"` // Broken; don't use this for now
}

// Fragments
type ruleRemovalRequest struct {
	Policy           types.Policy `dynamodbav:"Policy"`
	DeleteOnNextSync bool         `dynamodbav:"DeleteOnNextSync"`
}

type updateRulePolicyRequest struct {
	Policy       types.Policy `dynamodbav:"Policy,omitempty"`
	ExpiresAfter int64        `dynamodbav:"ExpiresAfter,omitempty"`
	Description  string       `dynamodbav:"Description,omitempty"`
}

func machineRulePK(machineID string) string {
	return fmt.Sprintf("%s%s", machineRulesPKPrefix, machineID)
}
func machineRuleSK(sha256 string, ruleType types.RuleType) string {
	return rules.RuleSortKeyFromTypeSHA(sha256, ruleType)
}
