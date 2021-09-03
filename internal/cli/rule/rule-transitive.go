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

	var ruleTransitiveCmd = &cobra.Command{
		Use:     "transitive",
		Aliases: []string{"allow-transitive"},
		Short:   "Create a rule that applies the AllowlistTransitive policy to the specified file",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] has already been validated as a file before this
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			time := clock.ConcreteTimeProvider{}

			return applyPolicyForPath(time, dynamodbClient, types.AllowlistTransitive, tf, rf)
		},
	}

	tf.AddTargetFlags(ruleTransitiveCmd)
	rf.AddRuleInfoFlags(ruleTransitiveCmd)

	RuleCmd.AddCommand(ruleTransitiveCmd)
}
