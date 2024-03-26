package rules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/airbnb/rudolph/internal/cli/flags"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	"github.com/airbnb/rudolph/pkg/model/machinerules"
	modelrules "github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

var (
	RulesCmd *cobra.Command
)

func init() {
	tf := flags.TargetFlags{}

	RulesCmd = &cobra.Command{
		Use:   "rules [--global|--machine=XXX]",
		Short: "List rules available on the current machine, a target machine, or globally",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] has already been validated as a file before this
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)

			limit := 40

			return rules(dynamodbClient, tf, limit)
		},
	}

	tf.AddTargetFlagsRules(RulesCmd)

	addRuleExportCommand()
	addRuleImportCommand()
}

func rules(client dynamodb.QueryAPI, tf flags.TargetFlags, limit int) error {
	var machineID string
	var err error

	if tf.IsGlobal {
		fmt.Println("==========================")
		fmt.Println("Retrieving global rules: ")
		fmt.Println("")

		rules, lastEvaluatedKey, err := globalrules.GetPaginatedGlobalRules(client, limit, nil)
		if err != nil {
			return fmt.Errorf("failed to get rules: %w", err)
		}

		ruleCount := len(*rules)
		if lastEvaluatedKey != nil {
			fmt.Printf("Retrieved more than %d rules:\n", ruleCount)
		} else {
			fmt.Printf("Retrieved %d rules:\n", ruleCount)
		}
		for i, rule := range *rules {
			fmt.Println("----- [", i, "] (", rule.SortKey, ")")
			fmt.Printf("%s: %s\n", renderRule(rule.SantaRule), rule.Description)
			fmt.Println("")
		}

	} else {
		fmt.Println("==========================")
		fmt.Printf("Retrieving for ")
		machineID, err = tf.GetMachineID()
		if err != nil {
			return fmt.Errorf("failed to get machineID: %w", err)
		}

		if tf.IsTargetSelf() {
			fmt.Printf("current machine, %s\n", machineID)
		} else {
			fmt.Printf("target machine, %s\n", machineID)
		}
		fmt.Println("")

		rules, err := machinerules.GetMachineRules(client, machineID)
		if err != nil {
			return fmt.Errorf("failed to GetMachineRules: %w", err)
		}

		fmt.Printf("Retrieved %d MachineRules:\n", len(*rules))
		for i, rule := range *rules {
			fmt.Println("----- [", i, "] (", rule.SortKey, ")")
			fmt.Printf("%s: %s\n", renderRule(rule.SantaRule), rule.Description)
			fmt.Println("")
		}
	}

	return nil
}

// String returns the rule in human readable string format
func renderRule(item modelrules.SantaRule) string {
	var predicate string
	switch item.RuleType {
	case types.RuleTypeBinary:
		predicate = "binary"
	case types.RuleTypeCertificate:
		predicate = "certificate"
	default:
		predicate = "?"
	}

	var verb string
	switch item.Policy {
	case types.RulePolicyAllowlist:
		verb = "allows"
	case types.RulePolicyBlocklist:
		verb = "blocks"
	case types.RulePolicySilentBlocklist:
		verb = "silent blocks"
	case types.RulePolicyRemove:
		verb = "removes"
	case types.RulePolicyAllowlistCompiler:
		verb = "allows compiler of"
	case types.RulePolicyAllowlistTransitive:
		verb = "allows transitive of"
	default:
		verb = "?"
	}

	return fmt.Sprintf("Rule %v %v (%s)", verb, predicate, item.SHA256)
}
