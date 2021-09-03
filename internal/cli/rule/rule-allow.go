package rule

import (
	"github.com/airbnb/rudolph/internal/cli/flags"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/spf13/cobra"
)

func init() {
	tf := flags.TargetFlags{}
	rf := flags.RuleInfoFlags{}

	var ruleAllowCmd = &cobra.Command{
		Use:   "allow [-f <file-path>|-s <sha>] -t <rule-type> [-m <machine-id>|--global]",
		Short: "Create a rule that applies the Allowlist policy to the specified file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			time := clock.ConcreteTimeProvider{}

			return applyPolicyForPath(time, dynamodbClient, types.Allowlist, tf, rf)
		},
	}

	tf.AddTargetFlags(ruleAllowCmd)
	rf.AddRuleInfoFlags(ruleAllowCmd)

	RuleCmd.AddCommand(ruleAllowCmd)
}
