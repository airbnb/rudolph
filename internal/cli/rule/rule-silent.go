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

	var ruleSilentCmd = &cobra.Command{
		Use:     "silent",
		Aliases: []string{"silentblock"},
		Short:   "Create a rule that applies the SilentBlocklist policy to the specified file",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] has already been validated as a file before this
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			time := clock.ConcreteTimeProvider{}

			return applyPolicyForPath(time, dynamodbClient, types.SilentBlocklist, tf, rf)
		},
	}

	tf.AddTargetFlags(ruleSilentCmd)
	rf.AddRuleInfoFlags(ruleSilentCmd)

	RuleCmd.AddCommand(ruleSilentCmd)
}
