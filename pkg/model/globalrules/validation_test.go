package globalrules

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

func Test_globalRuleValidation(t *testing.T) {
	type test struct {
		name        string
		rule        *GlobalRuleRow
		isValid     bool
		expectError bool
	}
	tests := []test{
		{
			name: "Binary#4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
						types.RuleTypeBinary,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeBinary,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
				},
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "Certificate#1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
						types.RuleTypeCertificate,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeCertificate,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
				},
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "Certificate#1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
						types.RuleTypeCertificate,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeCertificate,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
				},
			},
			isValid:     false,
			expectError: false,
		},
		{
			name: "TeamID#EQHXZ8M8AV",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"EQHXZ8M8AV",
						types.RuleTypeTeamID,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeTeamID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AV",
				},
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "TeamID#EQHXZ8M8AVAAAAA",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"EQHXZ8M8AVAAAAA",
						types.RuleTypeTeamID,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeTeamID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AVAAAAA",
				},
			},
			isValid:     false,
			expectError: false,
		},
		{
			name: "SigningID#EQHXZ8M8AV:com.google.Chrome",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"EQHXZ8M8AV:com.google.Chrome",
						types.RuleTypeSigningID,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AV:com.google.Chrome",
				},
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "SigningID#EQHXZ8M8AVAAAAA:com.google.Chrome",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"EQHXZ8M8AVAAAAA:com.google.Chrome",
						types.RuleTypeSigningID,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AVAAAAA:com.google.Chrome",
				},
			},
			isValid:     false,
			expectError: false,
		},
		{
			name: "SigningID#platform:com.apple.curl",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						"platform:com.apple.curl",
						types.RuleTypeSigningID,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "platform:com.apple.curl",
				},
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "SigningID#:com.apple.curl",
			rule: &GlobalRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: globalRulesPK,
					SortKey: globalRulesSK(
						":com.apple.curl",
						types.RuleTypeSigningID,
					),
				},
				Description: "",
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: ":com.apple.curl",
				},
			},
			isValid:     false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rule.globalRuleValidation()
			if (err != nil) != tt.expectError {
				t.Errorf("globalRuleValidation() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if got != tt.isValid {
				t.Errorf("globalRuleValidation() got = %v, want %v", got, tt.isValid)
				return
			}
		})
	}
}
