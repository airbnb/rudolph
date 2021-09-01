package globalrules

import (
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/pkg/errors"
)

func GetGlobalRuleBySortKey(client dynamodb.GetItemAPI, ruleSortKey string) (*GlobalRuleRow, error) {
	return getItemAsGlobalRule(client, globalRulesPK, ruleSortKey)
}

func GetGlobalRuleByShaType(client dynamodb.GetItemAPI, sha256 string, ruleType types.RuleType) (*GlobalRuleRow, error) {
	pk := globalRulesPK
	sk := globalRulesSK(sha256, ruleType)
	return getItemAsGlobalRule(client, pk, sk)
}

func getItemAsGlobalRule(client dynamodb.GetItemAPI, partitionKey string, sortKey string) (rule *GlobalRuleRow, err error) {
	output, err := client.GetItem(
		dynamodb.PrimaryKey{
			PartitionKey: partitionKey,
			SortKey:      sortKey,
		},
		false,
	)

	if err != nil {
		return
	}

	if len(output.Item) == 0 {
		return
	}

	err = attributevalue.UnmarshalMap(output.Item, &rule)

	if err != nil {
		err = errors.Wrap(err, "succeeded GetItem but failed to unmarshalMap into GlobalRuleRow")
		return
	}

	return
}
