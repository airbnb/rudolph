package rule

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/airbnb/rudolph/internal/cli/flags"
	"github.com/airbnb/rudolph/internal/cli/santa_sensor"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	"github.com/airbnb/rudolph/pkg/model/machinerules"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/spf13/cobra"
)

// The `rule` command itself does not take any flags or run anything, it's
// simply a passthrough to other subcommands

var (
	RuleCmd = &cobra.Command{
		Use:   "rule",
		Short: "Perform various rule operations",
	}
)

func applyPolicyForPath(timeProvider clock.TimeProvider, client dynamodb.DynamoDBClient, policy types.Policy, tf flags.TargetFlags, rf flags.RuleInfoFlags) (err error) {
	// Second, determine the rule type and identifier
	ruleType := (*rf.RuleType).AsRuleType()
	var description string
	var identifier string

	if *rf.FilePath != "" {
		fileInfo, err := santa_sensor.RunSantaFileInfo(*rf.FilePath)
		if err != nil {
			return fmt.Errorf("encountered an error while attempting to get file information for %q", *rf.FilePath)
		}
		identifier = fileInfo.SHA256
		description = fmt.Sprintf("%s from %s", fileInfo.Path, tf.SelfMachineID) // FIXME (derek.wang) tf.SelfMachineID is Not initialized.

		switch ruleType {
		case types.RuleTypeBinary:
			break
		case types.RuleTypeCertificate:
			if len(fileInfo.SigningChain) == 0 {
				return fmt.Errorf("NO SIGNING INFO FOUND FOR GIVEN BINARY")
			}
			if fileInfo.SigningChain[0].SHA256 == "" {
				return fmt.Errorf("NO CERTIFICATE SHA FOUND FOR GIVEN BINARY")
			}
			if fileInfo.SigningChain[0].CommonName == "" {
				return fmt.Errorf("NO CERTIFICATE NAME FOUND FOR GIVEN BINARY")
			}

			identifier = fileInfo.SigningChain[0].SHA256
			description = fmt.Sprintf("%v, by %v (%v)", fileInfo.SigningChain[0].CommonName, fileInfo.SigningChain[0].Organization, fileInfo.SigningChain[0].OrganizationalUnit)
		case types.RuleTypeTeamID:
			identifier = fileInfo.TeamID
		case types.RuleTypeSigningID:
			identifier = fileInfo.SigningID
		default:
			log.Printf("error (recovered): encountered unknown ruleType: (%+v)", ruleType)
			return fmt.Errorf("error (recovered): encountered unknown ruleType: (%+v)", ruleType)
		}
	} else if *rf.Identifier != "" {
		identifier = *rf.Identifier
	}

	// TODO
	// Query if there is an existing rule: and show the before/after
	// partitionKey := fmt.Sprintf("%s%s", store.MachineRulesPKPrefix, machineID)
	// ruleName :=
	// existingItem, err := store.GetRuleByPK(&store.DynamoDBPrimaryKey{
	// 	PartitionKey: partitionKey,
	// 	SortKey:      ruleName,
	// })

	policyDescription, err := policy.MarshalText()
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
			return fmt.Errorf("failed to get MachineID: %w", err)
		}
		// All args set up; send confirmation message
		if tf.IsTargetSelf() {
			suffix = " (This machine)"
		}
	}

	fmt.Println("Uploading the following rule:")
	fmt.Println("  MachineID:   ", machineID, suffix)
	fmt.Println("  Identifier/SHA256:      ", identifier)
	fmt.Println("  Policy:      ", policy, "  (", string(policyDescription), ")")
	fmt.Println("  RuleType:    ", ruleType, "  (", string(ruleTypeDescription), ")")
	fmt.Println("  Description: ", description)
	fmt.Println("")
	fmt.Println(`Apply changes? (Enter: "yes" or "ok")`)
	fmt.Print("> ")

	// Read confirmation
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if strings.ToLower(text) == "ok" || strings.ToLower(text) == "yes" {
		// Do rule creation
		if tf.IsGlobal {
			err = globalrules.AddNewGlobalRule(timeProvider, client, identifier, ruleType, policy, description)
		} else {
			expires := timeProvider.Now().Add(time.Hour * machinerules.MachineRuleDefaultExpirationHours).UTC()
			err = machinerules.AddNewMachineRule(client, machineID, identifier, ruleType, policy, description, expires)
		}
		if err != nil {
			return fmt.Errorf("could not upload rule to DynamoDB: %w", err)
		}
		fmt.Println("Successfully sent a rule to dynamodb")
	} else {
		fmt.Println("Well ok then")
	}
	fmt.Println("")

	return
}
