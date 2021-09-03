package rule

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/airbnb/rudolph/cmd/cli/flags"
	"github.com/airbnb/rudolph/cmd/cli/santa_sensor"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	"github.com/airbnb/rudolph/pkg/model/machinerules"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ruleHandler struct {
	timeProvider       clock.TimeProvider
	dynamodbClient     dynamodb.DynamoDBClient
	globalRuleUpdater  globalrules.GlobalRulesUpdater
	machineRuleUpdater machinerules.MachineRulesUpdater
}

func init() {
	tf := flags.TargetFlags{}
	rf := flags.RuleInfoFlags{}
	ru := flags.RuleUpdateFlags{}

	ruleHandler := ruleHandler{}

	var ruleUpdateCmd = &cobra.Command{
		Use:     "update <file-path>",
		Aliases: []string{"update"},
		Short:   "Update an existing rule",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] has already been validated as a file before this
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			client := dynamodb.GetClient(table, region)
			timeProvider := clock.ConcreteTimeProvider{}

			ruleHandler.dynamodbClient = client
			ruleHandler.timeProvider = timeProvider
			ruleHandler.globalRuleUpdater = globalrules.ConcreteGlobalRulesUpdater{
				ClockProvider: timeProvider,
				TransactWrite: client,
			}
			ruleHandler.machineRuleUpdater = machinerules.ConcreteMachineRulesUpdater{
				Updater: client,
			}

			return ruleHandler.updateRulePolicy(tf, rf, ru)
		},
	}

	tf.AddTargetFlags(ruleUpdateCmd)
	ru.AddRuleUpdateFlags(ruleUpdateCmd)
	rf.AddRuleInfoFlags(ruleUpdateCmd)

	RuleCmd.AddCommand(ruleUpdateCmd)
}

func (rh *ruleHandler) updateRulePolicy(tf flags.TargetFlags, rf flags.RuleInfoFlags, ru flags.RuleUpdateFlags) (err error) {
	// Determine the ruleType and rulePolicy and return the following types from RuleInfoFlags
	ruleType := rf.RuleType.AsRuleType()
	rulePolicy := ru.RulePolicy.AsRulePolicy()
	var description string
	var sha256 string
	if *rf.FilePath != "" {
		fileInfo, err := santa_sensor.RunSantaFileInfo(*rf.FilePath)
		if err != nil {
			return fmt.Errorf("encountered an error while attempting to get file information for %q", *rf.FilePath)
		}

		sha256 = fileInfo.SHA256
		description = fmt.Sprintf("%s from %s", fileInfo.Path, tf.SelfMachineID) // FIXME (derek.wang) tf.SelfMachineID is Not initialized.
		if ruleType == types.Certificate {
			if len(fileInfo.SigningChain) == 0 {
				return fmt.Errorf("NO SIGNING INFO FOUND FOR GIVEN BINARY")
			}
			if fileInfo.SigningChain[0].SHA256 == "" {
				return fmt.Errorf("NO CERTIFICATE SHA FOUND FOR GIVEN BINARY")
			}
			if fileInfo.SigningChain[0].CommonName == "" {
				return fmt.Errorf("NO CERTIFICATE NAME FOUND FOR GIVEN BINARY")
			}

			sha256 = fileInfo.SigningChain[0].SHA256
			description = fmt.Sprintf("%v, by %v (%v)", fileInfo.SigningChain[0].CommonName, fileInfo.SigningChain[0].Organization, fileInfo.SigningChain[0].OrganizationalUnit)
		}

	} else if *rf.SHA256 != "" {
		sha256 = *rf.SHA256
	}

	rulePolicyDescription, err := rulePolicy.MarshalText()
	if err != nil {
		return
	}

	ruleTypeDescription, err := ruleType.MarshalText()
	if err != nil {
		return
	}

	// First, determine which machine to apply
	machineID := "(Global)"
	suffix := ""
	if !tf.IsGlobal {
		machineID, err = tf.GetMachineID()
		if err != nil {
			return errors.Wrap(err, "Failed to get MachineID!")
		}
		// All args set up; send confirmation message
		if tf.IsTargetSelf() {
			suffix = " (This machine)"
		}
	}

	// Query if there is an existing rule: and show the before/after
	if machineID == "(Global)" {
		existingItem, err := globalrules.GetGlobalRuleByShaType(rh.dynamodbClient, sha256, ruleType)
		if err != nil {
			return err
		}

		// If nil, no rule exists
		if existingItem == nil {
			return errors.New("no global rule exists")
		}

		rulePolicyDescription, err := existingItem.Policy.MarshalText()
		if err != nil {
			return err
		}

		ruleTypeDescription, err := existingItem.RuleType.MarshalText()
		if err != nil {
			return err
		}

		fmt.Println("The current rule is rule:")
		fmt.Println("  MachineID:   ", machineID, suffix)
		fmt.Println("  SHA256:      ", existingItem.SHA256)
		fmt.Println("  Policy:      ", existingItem.Policy, "  (", string(rulePolicyDescription), ")")
		fmt.Println("  RuleType:    ", existingItem.RuleType, "  (", string(ruleTypeDescription), ")")
		fmt.Println("  Description: ", description)
	} else {
		existingItem, err := machinerules.GetMachineRuleByShaType(rh.dynamodbClient, machineID, sha256, ruleType)
		if err != nil {
			return err
		}

		// If nil, no rule exists
		if existingItem == nil {
			return errors.New("no machine rule exists")
		}

		rulePolicyDescription, err := existingItem.Policy.MarshalText()
		if err != nil {
			return err
		}

		ruleTypeDescription, err := existingItem.RuleType.MarshalText()
		if err != nil {
			return err
		}

		fmt.Println("The current rule is rule:")
		fmt.Println("  MachineID:   ", machineID, suffix)
		fmt.Println("  SHA256:      ", existingItem.SHA256)
		fmt.Println("  Policy:      ", existingItem.Policy, "  (", string(rulePolicyDescription), ")")
		fmt.Println("  RuleType:    ", existingItem.RuleType, "  (", string(ruleTypeDescription), ")")
		fmt.Println("  Description: ", description)
	}

	fmt.Println("Updating the rule to the following:")
	fmt.Println("  MachineID:   ", machineID, suffix)
	fmt.Println("  SHA256:      ", sha256)
	fmt.Println("  Policy:      ", rulePolicy, "  (", string(rulePolicyDescription), ")")
	fmt.Println("  RuleType:    ", ruleType, "  (", string(ruleTypeDescription), ")")
	fmt.Println("  Description: ", description)
	fmt.Println("")
	fmt.Println(`Apply changes? (Enter: "yes" or "ok")`)
	fmt.Print("> ")

	// Read confirmation
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if text == "ok" || text == "yes" {
		// Do rule update
		if tf.IsGlobal {
			err = rh.globalRuleUpdater.UpdateGlobalRule(sha256, ruleType, rulePolicy)
		} else {
			err = rh.machineRuleUpdater.UpdateMachineRulePolicy(machineID, sha256, ruleType, rulePolicy)
		}
		if err != nil {
			return errors.Wrap(err, "Could not upload rule to dynamodb")
		}
		fmt.Println("Successfully updated the rule on dynamodb")
	} else {
		fmt.Println("Well ok then")
	}
	fmt.Println("")

	return
}
