package feedrules

import (
	"testing"
	"time"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

func Test_globalRuleValidation(t *testing.T) {
	timeProvider := clock.FrozenTimeProvider{
		Current: time.Now(),
	}
	type test struct {
		name        string
		feedRuleRow *FeedRuleRow
		isValid     bool
		expectError bool
	}
	tests := []test{
		{
			name: "Binary#4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
						types.RuleTypeBinary,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeBinary,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "Certificate#1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
						types.RuleTypeCertificate,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeCertificate,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "Certificate#1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
						types.RuleTypeCertificate,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeCertificate,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     false,
			expectError: false,
		},
		{
			name: "TeamID#EQHXZ8M8AV",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"EQHXZ8M8AV",
						types.RuleTypeTeamID,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeTeamID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AV",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "TeamID#EQHXZ8M8AVAAAAA",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"EQHXZ8M8AVAAAAA",
						types.RuleTypeTeamID,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeTeamID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AVAAAAA",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     false,
			expectError: false,
		},
		{
			name: "SigningID#EQHXZ8M8AV:com.google.Chrome",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"EQHXZ8M8AV:com.google.Chrome",
						types.RuleTypeSigningID,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AV:com.google.Chrome",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "SigningID#EQHXZ8M8AVAAAAA:com.google.Chrome",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"EQHXZ8M8AVAAAAA:com.google.Chrome",
						types.RuleTypeSigningID,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "EQHXZ8M8AVAAAAA:com.google.Chrome",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     false,
			expectError: false,
		},
		{
			name: "SigningID#platform:com.apple.curl",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						"platform:com.apple.curl",
						types.RuleTypeSigningID,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: "platform:com.apple.curl",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     true,
			expectError: false,
		},
		{
			name: "SigningID#:com.apple.curl",
			feedRuleRow: &FeedRuleRow{
				PrimaryKey: dynamodb.PrimaryKey{
					PartitionKey: feedRulesPK,
					SortKey: feedRulesSK(
						timeProvider,
						":com.apple.curl",
						types.RuleTypeSigningID,
					),
				},
				SantaRule: rules.SantaRule{
					RuleType:   types.RuleTypeSigningID,
					Policy:     types.RulePolicyAllowlist,
					Identifier: ":com.apple.curl",
				},
				ExpiresAfter: GetSyncStateExpiresAfter(timeProvider),
				DataType:     GetDataType(),
			},
			isValid:     false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.feedRuleRow.feedRuleRowValidation()
			if (err != nil) != tt.expectError {
				t.Errorf("feedRuleRowValidation() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if got != tt.isValid {
				t.Errorf("feedRuleRowValidation() got = %v, want %v", got, tt.isValid)
				return
			}
		})
	}
}
