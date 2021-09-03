package main

import (
	"os"

	"github.com/airbnb/rudolph/cmd/cli/config"
	"github.com/airbnb/rudolph/cmd/cli/info"
	"github.com/airbnb/rudolph/cmd/cli/lookup"
	"github.com/airbnb/rudolph/cmd/cli/repair"
	"github.com/airbnb/rudolph/cmd/cli/rule"
	"github.com/airbnb/rudolph/cmd/cli/rules"
	"github.com/spf13/cobra"
)

/*
Usage:
  First run `% make rudolph` to generate the rudolph binary to build/macos/rudolph

export ENV={YOUR_ENVIRONMENT}
  You will have to supply a ENV environment variable with every execution. This variable corresponds to a
  directory inside of the terraform/deployments/ directory.

  ./rudolph [COMMAND]

     ./rudolph info
       Dump information about your current machine and whatever.

     ./rudolph rule allow --rule-type binary [--global] $PATH
       Creates a new rule based upon the given binary and uploads it to the currently configured
       Rudolph server. Can do both binary and certificate rules.

     ./rudolph rules [--global]
       Queries DDB and returns all rules pertinent to either your machine or available globally.

	 ./rudolph config [--global]
		Creates, modifies, or retrieves the current configuration for either a machine or globally.

*/

func init() {
	RootCmd.PersistentFlags().StringVar(&env, "ENV", "", "Environment of Rudolph deployment")
	RootCmd.PersistentFlags().StringVar(&region, "region", "", ".")
	RootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", ".")
	RootCmd.PersistentFlags().StringVar(&dynamodbTableName, "dynamodb_table", "", ".")

	// Add subcommands
	RootCmd.AddCommand(info.InfoCmd)
	RootCmd.AddCommand(rule.RuleCmd)
	RootCmd.AddCommand(rules.RulesCmd)
	RootCmd.AddCommand(config.ConfigCmd)
	RootCmd.AddCommand(repair.RepairCmd)
	RootCmd.AddCommand(lookup.LookupCmd)
}

var (
	env               string
	region            string
	prefix            string
	dynamodbTableName string
)

// RootCmd is the entry point command for the CLI, exported for use elsewhere
var RootCmd = &cobra.Command{
	Use:          "rudolph",
	Short:        "cli for interacting with Santa server",
	SilenceUsage: true,
	Long:         "cli for interacting with Santa server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		_, err := retrieveConfig(cmd)
		if err != nil {
			return err
		}

		return nil
	},
}

// Execute is the main entry point for the CLI
func Execute(version string) {
	RootCmd.Version = version
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
