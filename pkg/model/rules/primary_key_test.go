package rules

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/types"
)

func Test_RuleSortKeyFromTypeIdentifier(t *testing.T) {
	type test struct {
		identifier string
		ruleType   types.RuleType
		sortKey    string
	}
	tests := []test{
		{
			identifier: "61977d6006459c4cefe9b988a453589946224957bfc07b262cd7ca1b7a61e04e",
			ruleType:   types.RuleTypeBinary,
			sortKey: RuleSortKeyFromTypeIdentifier(
				"61977d6006459c4cefe9b988a453589946224957bfc07b262cd7ca1b7a61e04e",
				types.RuleTypeBinary,
			),
		},
		{
			identifier: "61977d6006459c4cefe9b988a453589946224957bfc07b262cd7ca1b7a61e04e",
			ruleType:   types.RuleTypeCertificate,
			sortKey: RuleSortKeyFromTypeIdentifier(
				"61977d6006459c4cefe9b988a453589946224957bfc07b262cd7ca1b7a61e04e",
				types.RuleTypeCertificate,
			),
		},
		{
			identifier: "EQHXZ8M8AV",
			ruleType:   types.RuleTypeTeamID,
			sortKey: RuleSortKeyFromTypeIdentifier(
				"EQHXZ8M8AV",
				types.RuleTypeTeamID,
			),
		},
		{
			identifier: "EQHXZ8M8AV:com.google.Chrome",
			ruleType:   types.RuleTypeSigningID,
			sortKey: RuleSortKeyFromTypeIdentifier(
				"EQHXZ8M8AV:com.google.Chrome",
				types.RuleTypeSigningID,
			),
		},
		{
			identifier: "EQHXZ8M8AV:com.google.Chrome",
			ruleType:   0,
			sortKey:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.identifier, func(t *testing.T) {
			got := RuleSortKeyFromTypeIdentifier(
				tt.identifier,
				tt.ruleType,
			)
			if got != tt.sortKey {
				t.Errorf("RuleSortKeyFromTypeIdentifier() got = %v, want %v", got, tt.sortKey)
				return
			}
			if tt.ruleType == types.RuleTypeBinary || tt.ruleType == types.RuleTypeCertificate {
				got = RuleSortKeyFromTypeSHA(
					tt.identifier,
					tt.ruleType,
				)
				if got != tt.sortKey {
					t.Errorf("RuleSortKeyFromTypeSHA() got = %v, want %v", got, tt.sortKey)
					return
				}
			}
		})
	}
}
