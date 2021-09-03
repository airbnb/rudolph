package rule

import (
	"github.com/airbnb/rudolph/internal/cli/flags"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	"github.com/airbnb/rudolph/pkg/model/machinerules"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// The `rule remove` command takes slightly different arguments than the other
// `rule` commands so we'll handle it a bit differently here
func init() {
	tf := flags.TargetFlags{}

	var removeRuleCmd = &cobra.Command{
		Use:     "remove <rule-name>",
		Aliases: []string{"delete"},
		Short:   "Removes/deletes a rule from the backing store",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)

			globalRemover := globalrules.ConcreteRuleRemovalService{
				TimeProvider: clock.ConcreteTimeProvider{},
				Getter:       dynamodbClient,
				Transacter:   dynamodbClient,
			}
			machineRuleRemover := machinerules.ConcreteRuleRemovalService{
				Getter:  dynamodbClient,
				Updater: dynamodbClient,
			}

			return removeRule(globalRemover, machineRuleRemover, args[0], tf)
		},
	}

	tf.AddTargetFlags(removeRuleCmd)
	RuleCmd.AddCommand(removeRuleCmd)
}

func removeRule(globalRemover globalrules.RuleRemovalService, machineRuleRemover machinerules.RuleRemovalService, ruleName string, tf flags.TargetFlags) error {
	if !tf.IsGlobal {
		machineID, err := tf.GetMachineID()
		if err != nil {
			return errors.Wrap(err, "Failed to get MachineID!")
		}
		return machineRuleRemover.RemoveMachineRule(machineID, ruleName)
	}

	idempotencyKey := uuid.NewString()

	return globalRemover.RemoveGlobalRule(ruleName, idempotencyKey)
}
