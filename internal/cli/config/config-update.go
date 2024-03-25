package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/airbnb/rudolph/internal/cli/flags"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	"github.com/spf13/cobra"
)

func init() {
	var (
		clientModeArg flags.ClientMode
	)

	tf := flags.TargetFlags{}

	var configUpdateClientModeCmd = &cobra.Command{
		Use:   "update [-m <machine-id>|--global] [-c <ClientMode - 'monitor' or 'lockdown'>|--client-mode]",
		Short: "Update the client-mode globally or for a specific machine UUID",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			timeProvider := clock.ConcreteTimeProvider{}

			service := machineconfiguration.GetMachineConfigurationService(dynamodbClient, timeProvider)

			return updateConfig(
				service,
				tf,
				clientModeArg,
			)
		},
	}

	tf.AddTargetFlags(configUpdateClientModeCmd)

	// client-mode should be one of "monitor" or "lockdown"
	configUpdateClientModeCmd.Flags().VarP(&clientModeArg, "client-mode", "c", `type of client mode being applied. valid options are: "monitor" or "lockdown"`)
	// Mark as a required flag
	_ = configUpdateClientModeCmd.MarkFlagRequired("client-mode")

	ConfigCmd.AddCommand(configUpdateClientModeCmd)
}

func updateConfig(
	service machineconfiguration.MachineConfigurationService,
	tf flags.TargetFlags,
	clientModeArg flags.ClientMode) (err error) {
	clientMode := clientModeArg.AsClientMode()

	// Get machineID from flags
	var machineID string
	if tf.IsGlobal {
		machineID = "(Global)"
	} else {
		machineID, err = tf.GetMachineID()
		if err != nil {
			return fmt.Errorf("failed to get MachineID: %w", err)
		}
	}

	suffix := ""
	if !(tf.IsGlobal) && tf.IsTargetSelf() {
		suffix = "-->( This machine )"
	}

	clientModeText, err := clientMode.MarshalText()
	if err != nil {
		return
	}

	// Print the output for visual confirmation via a nicely tab corrected output
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', tabwriter.AlignRight)

	fmt.Println("Setting the following configuration")
	fmt.Println()
	fmt.Fprintln(writer, "Config\t Setting")
	fmt.Fprintln(writer, "MachineID:\t", machineID, suffix)
	fmt.Fprintln(writer, "ClientMode:\t", clientMode, "-->(", string(clientModeText), ")")
	writer.Flush()
	fmt.Println()
	fmt.Println(`Apply changes? (Enter: "yes" or "ok")`)
	fmt.Print("> ")

	// Read confirmation
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if text == "ok" || text == "yes" {
		fmt.Println("Sending the configuration to DynamoDB...")
	} else {
		fmt.Println("Confirmation not successful...")
		return
	}
	updateRequest := machineconfiguration.MachineConfigurationUpdateRequest{
		ClientMode: &clientMode,
	}

	if tf.IsGlobal {
		_, err = service.UpdateGlobalConfig(updateRequest)
	} else {
		_, err = service.UpdateMachineConfig(machineID, updateRequest)
	}

	if err != nil {
		return fmt.Errorf("error writing the configuration to the sync server: %w", err)
	} else {
		fmt.Println("Success! Configuration was sent properly to DynamoDB...")
	}
	return

}
