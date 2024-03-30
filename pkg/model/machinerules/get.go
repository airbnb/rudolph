package machinerules

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func getItemAsMachineRule(client dynamodb.GetItemAPI, partitionKey string, sortKey string) (rule *MachineRuleRow, err error) {
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
		err = fmt.Errorf("succeeded GetItem but failed to unmarshalMap into GlobalRuleRow: %w", err)
		return
	}

	return
}

// @deprecated
func GetMachineRuleByShaType(client dynamodb.GetItemAPI, machineID string, sha256 string, ruleType types.RuleType) (rule *MachineRuleRow, err error) {
	return GetMachineRuleByIdentifierType(client, machineID, sha256, ruleType)
}

func GetMachineRuleByIdentifierType(client dynamodb.GetItemAPI, machineID string, identifier string, ruleType types.RuleType) (rule *MachineRuleRow, err error) {
	pk := machineRulePK(machineID)
	sk := machineRuleSK(identifier, ruleType)
	return getItemAsMachineRule(client, pk, sk)
}
