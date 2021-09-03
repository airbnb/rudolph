package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/airbnb/rudolph/cmd/cli/flags"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machineconfiguration"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	var (
		clientModeArg            flags.ClientMode
		blockedPathRegexArg      string
		allowedPathRegexArg      string
		batchSizeArg             int
		isEnableBundles          bool
		isEnabledTransitiveRules bool
		isCleanSync              bool
		fullSyncIntervalArg      int
		uploadLogsUrlArgs        string
	)

	tf := flags.TargetFlags{}

	var configSetCmd = &cobra.Command{
		Use:   "set [-m <machine-id>|--global] [-c <ClientMode - 'monitor' or 'lockdown'>|--client-mode]",
		Short: "Create a configuration and set globally or a specific machine UUID",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			timeProvider := clock.ConcreteTimeProvider{}

			service := machineconfiguration.GetMachineConfigurationService(dynamodbClient, timeProvider)

			return applyConfig(
				service,
				tf,
				clientModeArg,
				blockedPathRegexArg,
				allowedPathRegexArg,
				batchSizeArg,
				isEnableBundles,
				isEnabledTransitiveRules,
				isCleanSync,
				uploadLogsUrlArgs,
				fullSyncIntervalArg,
			)
		},
	}

	tf.AddTargetFlags(configSetCmd)

	// client-mode should be one of "monitor" or "lockdown"
	configSetCmd.Flags().VarP(&clientModeArg, "client-mode", "c", `type of client mode being applied. valid options are: "monitor" or "lockdown"`)
	// Mark as a required flag
	_ = configSetCmd.MarkFlagRequired("client-mode")

	// Flags defining blocked and allowed regex paths
	configSetCmd.Flags().StringVarP(&blockedPathRegexArg, "blocked-paths", "b", "", `A comma separated list of regex paths to be blocked`)
	configSetCmd.Flags().StringVarP(&allowedPathRegexArg, "allowed-paths", "a", "", `A comma separated list of regex paths to be allowed`)

	// Flags defining batchSize and events to sync back to the server
	configSetCmd.Flags().IntVarP(&batchSizeArg, "batch-size", "s", 50, "Int value to define the number of rules to download per sync request. Defaults to '50'")
	configSetCmd.Flags().BoolVar(&isEnableBundles, "bundles", false, "Enables bundle events to be uploaded back to the sync server")
	configSetCmd.Flags().BoolVar(&isEnabledTransitiveRules, "transitive-rules", false, "Enables the usage of transitive rules")

	// Flags define sync behavior and upload logs url
	configSetCmd.Flags().BoolVar(&isCleanSync, "clean-sync", false, "Enforces that the next sync process will be a clean-sync. Only enforced on machine specific configurations")
	configSetCmd.Flags().StringVarP(&uploadLogsUrlArgs, "upload-logs", "u", "", "Set an upload logs URL link here. If using the sync server, define the API endpoint of the sync server logs upload path")

	configSetCmd.Flags().IntVarP(&fullSyncIntervalArg, "full-sync-interval", "f", machineconfiguration.DefaultFullSyncInterval, "Set full sync interval in seconds (default 600)")

	ConfigCmd.AddCommand(configSetCmd)
}

func applyConfig(
	service machineconfiguration.MachineConfigurationService,
	tf flags.TargetFlags,
	clientModeArg flags.ClientMode,
	blockedPathRegexArg string,
	allowedPathRegexArg string,
	batchSizeArg int,
	isEnableBundles bool,
	isEnabledTransitiveRules bool,
	isCleanSync bool,
	uploadLogsUrlArgs string,
	fullSyncIntervalArg int) (err error) {
	clientMode := clientModeArg.AsClientMode()

	// Get machineID from flags
	var machineID string
	if tf.IsGlobal {
		machineID = "(Global)"
	} else {
		machineID, err = tf.GetMachineID()
		if err != nil {
			return errors.Wrap(err, "Failed to get MachineID!")
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
	fmt.Fprintln(writer, "BlockedPathRegex:\t \"", blockedPathRegexArg, "\"")
	fmt.Fprintln(writer, "AllowedPathRegex:\t \"", allowedPathRegexArg, "\"")
	fmt.Fprintln(writer, "BatchSize:\t", batchSizeArg)
	fmt.Fprintln(writer, "BundlesEnabled:\t", isEnableBundles)
	fmt.Fprintln(writer, "EnabledTransitiveRules:\t", isEnabledTransitiveRules)
	fmt.Fprintln(writer, "CleanSync:\t", isCleanSync)
	fmt.Fprintln(writer, "FullSyncInterval:\t", fullSyncIntervalArg)
	fmt.Fprintln(writer, "UploadLogUrl:\t \"", uploadLogsUrlArgs, "\"")
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

	newConfig := machineconfiguration.MachineConfiguration{
		ClientMode:             clientMode,
		AllowedPathRegex:       allowedPathRegexArg,
		BlockedPathRegex:       blockedPathRegexArg,
		BatchSize:              batchSizeArg,
		EnableBundles:          isEnableBundles,
		EnabledTransitiveRules: isEnabledTransitiveRules,
		CleanSync:              isCleanSync,
		FullSyncInterval:       fullSyncIntervalArg,
		UploadLogsURL:          uploadLogsUrlArgs,
	}

	if tf.IsGlobal {
		err = service.SetGlobalConfig(newConfig)
	} else {
		err = service.SetMachineConfig(machineID, newConfig)
	}

	if err != nil {
		log.Print(errors.Wrapf(err, " error writing the configuration to the sync server..."))
	} else {
		fmt.Println("Success! Configuration was sent properly to DynamoDB...")
	}
	return

}
