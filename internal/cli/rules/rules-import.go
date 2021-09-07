package rules

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"github.com/airbnb/rudolph/internal/csv"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	rudolphrules "github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

func addRuleImportCommand() {
	var filename string

	var ruleImportCmd = &cobra.Command{
		Use:     "import <file-name>",
		Aliases: []string{"rules-import"},
		Short:   "Imports riles from a csv file",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] has already been validated as a file before this
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)

			return runImport(dynamodbClient, clock.ConcreteTimeProvider{}, filename)
		},
	}

	ruleImportCmd.Flags().StringVarP(&filename, "filename", "f", "", "The filename")
	_ = ruleImportCmd.MarkFlagRequired("filename")

	RulesCmd.AddCommand(ruleImportCmd)
}

func runImport(
	client dynamodb.DynamoDBClient,
	timeProvider clock.TimeProvider,
	filename string,
) (err error) {
	csvLines := make(chan map[string]string)

	wg := new(sync.WaitGroup)

	go func() {
		for {
			line := <-csvLines

			sha256, ok := line["sha256"]
			if !ok {
				panic("no sha256")
			}
			ruleTypeStr, ok := line["type"]
			if !ok {
				panic("no type")
			}
			policyStr, ok := line["policy"]
			if !ok {
				panic("no policy")
			}
			description, ok := line["description"]
			if !ok {
				description = ""
			}
			var ruleType types.RuleType
			err := ruleType.UnmarshalText([]byte(ruleTypeStr))
			if err != nil {
				panic("invalid ruletype")
			}
			var policy types.Policy
			err = policy.UnmarshalText([]byte(policyStr))
			if err != nil {
				panic("invalid policy")
			}

			suffix := ""
			if ruleType == types.Certificate {
				suffix = " (Cert)"
			}

			if policy == types.RulePolicyRemove {
				fmt.Printf("  Removing rule: [%s]\n", sha256)
				sortkey := rudolphrules.RuleSortKeyFromTypeSHA(sha256, ruleType)
				err = globalrules.RemoveGlobalRule(
					timeProvider,
					client,
					client,
					sortkey,
					"",
				)

			} else {
				fmt.Printf("  Writing rule: [%s] %s%s\n", policyStr, sha256, suffix)
				err = globalrules.AddNewGlobalRule(
					timeProvider,
					client,
					sha256,
					ruleType,
					policy,
					description,
				)
			}

			if err != nil {
				panic(err)
			}

			wg.Done()
		}
	}()

	err = csv.ParseCsvFile(filename, csvLines, wg)
	if err != nil {
		panic(err)
	}

	wg.Wait()

	return nil
}
