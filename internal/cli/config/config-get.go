package config

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/airbnb/rudolph/internal/cli/flags"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"

	"github.com/spf13/cobra"
)

func init() {
	tf := flags.TargetFlags{}

	var configGetCmd = &cobra.Command{
		Use:   "get [-m <machine-id>|--global]",
		Short: "Get the current global or specific machine UUID specific configuration from the sync server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			timeProvider := clock.ConcreteTimeProvider{}

			service := machineconfiguration.GetMachineConfigurationService(dynamodbClient, timeProvider)

			return getConfig(service, tf)
		},
	}

	tf.AddTargetFlags(configGetCmd)

	ConfigCmd.AddCommand(configGetCmd)
}

func getConfig(service machineconfiguration.MachineConfigurationService, tf flags.TargetFlags) (err error) {
	// Get machineID from flags
	var subject string
	var config machineconfiguration.MachineConfiguration

	// Prompt the user if the global configuration is being retrieved or the machine UUID specific configuration from the get go...
	if tf.IsGlobal {
		fmt.Println("Retrieving the global configuration...")
		subject = "All Machines"

		tmpconfig, _, eerr := service.GetIntendedGlobalConfig()
		if eerr != nil {
			err = fmt.Errorf("failed to do DynamoDB get: %w", eerr)
			return
		}
		config = tmpconfig
	} else {
		machineID, eerr := tf.GetMachineID()
		if eerr != nil {
			err = fmt.Errorf("failed to get MachineUUID: %w", eerr)
			return
		}
		if tf.IsTargetSelf() {
			fmt.Printf("Retreiving the machine configuration for (%s) (Current Machine)\n", machineID)
			subject = fmt.Sprintf("Machine (%s) (This Machine)", machineID)
		} else {
			fmt.Printf("Retreiving the machine specific configuration for machine UUID: %s\n", machineID)
			subject = fmt.Sprintf("Machine (%s)", machineID)
		}
		tmpconfig, eerr := service.GetIntendedConfig(machineID)
		if eerr != nil {
			err = fmt.Errorf("failed to do dynamodb get: %w", eerr)
			return
		}
		config = tmpconfig
	}

	// Marshal the ClientMode out of the config to easily display the config clientmode content
	clientModeText, err := config.ClientMode.MarshalText()
	if err != nil {
		return
	}

	// Alert the user if this configuration is for the same machine requesting the call

	// Print the output for visual confirmation via a nicely tab corrected output
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', tabwriter.AlignRight)

	fmt.Println("Sync server returned the following configuration")
	fmt.Println()
	fmt.Fprintln(writer, "Config\t Setting")
	fmt.Fprintln(writer, "Target:\t", subject)
	fmt.Fprintln(writer, "ClientMode:\t", config.ClientMode, "--> (", string(clientModeText), ")")
	fmt.Fprintln(writer, "BlockedPathRegex:\t \"", config.BlockedPathRegex, "\"")
	fmt.Fprintln(writer, "AllowedPathRegex:\t \"", config.AllowedPathRegex, "\"")
	fmt.Fprintln(writer, "BatchSize:\t", config.BatchSize)
	fmt.Fprintln(writer, "BundlesEnabled:\t", config.EnableBundles)
	fmt.Fprintln(writer, "EnabledTransitiveRules:\t", config.EnabledTransitiveRules)
	fmt.Fprintln(writer, "CleanSync:\t", config.CleanSync)
	fmt.Fprintln(writer, "FullSyncInterval:\t", config.FullSyncInterval)
	fmt.Fprintln(writer, "UploadLogUrl:\t \"", config.UploadLogsURL, "\"")
	writer.Flush()

	return
}
