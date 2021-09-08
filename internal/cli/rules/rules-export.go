package rules

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/airbnb/rudolph/internal/csv"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
)

func addRuleExportCommand() {
	var filename string
	var ruleExportCmd = &cobra.Command{
		Use:     "export  <file-name>",
		Aliases: []string{"rules-export"},
		Short:   "Export all rules into a csv file",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] has already been validated as a file before this
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)

			return runExport(dynamodbClient, filename)
		},
	}

	ruleExportCmd.Flags().StringVarP(&filename, "filename", "f", "", "The filename")
	_ = ruleExportCmd.MarkFlagRequired("filename")

	RulesCmd.AddCommand(ruleExportCmd)
}

func runExport(
	client dynamodb.QueryAPI,
	filename string,
) (err error) {
	csvRules := make(chan []string)

	header := []string{
		"sha256",
		"type",
		"policy",
		"custom_msg",
		"description",
	}

	wg, err := csv.WriteCsvFile(filename, header, csvRules)
	if err != nil {
		panic(err)
	}

	var key *dynamodb.PrimaryKey
	for {
		rules, nextkey, inerr := globalrules.GetPaginatedGlobalRules(client, 50, key)
		if inerr != nil {
			err = errors.Wrap(inerr, "something went wrong querying global rules")
			return
		}
		if len(*rules) == 0 {
			break
		}

		for _, rule := range *rules {
			ruleType, _ := rule.RuleType.MarshalText()
			policy, _ := rule.Policy.MarshalText()
			record := []string{
				rule.SHA256,
				string(ruleType),
				string(policy),
				rule.CustomMessage,
				rule.Description,
			}

			csvRules <- record
		}

		if nextkey == nil {
			break
		}
		key = nextkey
	}
	close(csvRules)

	wg.Wait()

	return nil
}
