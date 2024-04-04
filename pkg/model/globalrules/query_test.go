package globalrules

// generate a unit test for the GetPaginatedGlobalRules function

import (
	"testing"

	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type mockQuery func(input *awsdynamodb.QueryInput) (*awsdynamodb.QueryOutput, error)

func (m mockQuery) Query(input *awsdynamodb.QueryInput) (*awsdynamodb.QueryOutput, error) {
	return m(input)
}

func TestGetPaginatedGlobalRules(t *testing.T) {
	rules := []rules.SantaRule{
		{
			RuleType:   types.RuleTypeBinary,
			Policy:     types.RulePolicyAllowlist,
			Identifier: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9820",
		},
		{
			RuleType: types.RuleTypeBinary,
			Policy:   types.RulePolicyAllowlist,
			SHA256:   "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9821",
		},
		{
			RuleType:   types.RuleTypeCertificate,
			Policy:     types.RulePolicyAllowlist,
			Identifier: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9822",
		},
		{
			RuleType: types.RuleTypeCertificate,
			Policy:   types.RulePolicyAllowlist,
			SHA256:   "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9823",
		},
		{
			RuleType:   types.RuleTypeTeamID,
			Policy:     types.RulePolicyAllowlist,
			Identifier: "EQHXZ8M8AV",
		},
		{
			RuleType:   types.RuleTypeSigningID,
			Policy:     types.RulePolicyAllowlist,
			Identifier: "platform:com.apple.Safari",
		},
		{
			RuleType:   types.RuleTypeSigningID,
			Policy:     types.RulePolicyAllowlist,
			Identifier: "EQHXZ8M8AV:com.google.Chrome",
		},
		{
			RuleType:   types.RuleTypeSigningID,
			Policy:     types.RulePolicySilentBlocklist,
			Identifier: "BJ4HAAB9B3:us.zoom.xos",
		},
	}

	mockQueryClient := mockQuery(
		func(input *awsdynamodb.QueryInput) (*awsdynamodb.QueryOutput, error) {
			items := make([]map[string]awstypes.AttributeValue, len(rules))

			for i, rule := range rules {
				ruleType, err := rule.RuleType.MarshalDynamoDBAttributeValue()
				if err != nil {
					return nil, err
				}

				policy, err := rule.Policy.MarshalDynamoDBAttributeValue()
				if err != nil {
					return nil, err
				}
				items[i] = map[string]awstypes.AttributeValue{
					"PK": &awstypes.AttributeValueMemberS{
						Value: globalRulesPK,
					},
					"SK": &awstypes.AttributeValueMemberS{
						Value: globalRulesSK(rule.Identifier, rule.RuleType),
					},
					"Identifier": &awstypes.AttributeValueMemberS{
						Value: rule.Identifier,
					},
					"RuleType": ruleType,
					"Policy":   policy,
				}
			}

			return &awsdynamodb.QueryOutput{
				Items: items,
				Count: int32(len(items)),
			}, nil
		},
	)

	globalRules, lastEvaluatedKey, err := GetPaginatedGlobalRules(
		mockQueryClient,
		1000,
		nil,
	)
	assert.Nil(t, err)
	assert.Nil(t, lastEvaluatedKey)
	if len(rules) != len(globalRules) {
		t.Fatalf("expected %d rules, got %d", len(rules), len(globalRules))
	}

	for i, globalRule := range globalRules {
		assert.Equal(t, globalRulesPK, globalRule.PrimaryKey.PartitionKey)
		assert.Equal(t, globalRulesSK(rules[i].Identifier, rules[i].RuleType), globalRule.PrimaryKey.SortKey)
		assert.Equal(t, rules[i].Identifier, globalRule.SantaRule.Identifier)
		assert.Equal(t, rules[i].RuleType, globalRule.RuleType)
		assert.Equal(t, rules[i].Policy, globalRule.SantaRule.Policy)
	}
}
