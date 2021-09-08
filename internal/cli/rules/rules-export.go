package rules

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/airbnb/rudolph/internal/csv"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	"github.com/airbnb/rudolph/pkg/types"
)

func addRuleExportCommand() {
	var filename string
	var format string
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

			return runExport(dynamodbClient, filename, format)
		},
	}

	ruleExportCmd.Flags().StringVarP(&filename, "filename", "f", "", "The filename")
	_ = ruleExportCmd.MarkFlagRequired("filename")

	ruleExportCmd.Flags().StringVarP(&format, "fileformat", "t", "csv", "File format (one of: [json|csv])")

	RulesCmd.AddCommand(ruleExportCmd)
}

func runExport(
	client dynamodb.QueryAPI,
	filename string,
	format string,
) (err error) {
	switch format {
	case "json":
		return runJsonExport(client, filename)
	case "csv":
		return runCsvExport(client, filename)
	}
	return
}

type fileRule struct {
	RuleType      types.RuleType `json:"type"`
	Policy        types.Policy   `json:"policy"`
	SHA256        string         `json:"sha256"`
	CustomMessage string         `json:"custom_msg,omitempty"`
	Description   string         `json:"description"`
}

func runJsonExport(client dynamodb.QueryAPI, filename string) (err error) {
	var jsonRules []fileRule
	err = getRules(client, func(rule globalrules.GlobalRuleRow) (err error) {
		jsonRules = append(jsonRules, fileRule{
			SHA256:        rule.SHA256,
			RuleType:      rule.RuleType,
			Policy:        rule.Policy,
			CustomMessage: rule.CustomMessage,
			Description:   rule.Description,
		})
		return
	})
	if err != nil {
		return
	}

	jsondata, err := json.MarshalIndent(jsonRules, "", "  ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, jsondata, 0644)
	if err != nil {
		return
	}

	return
}

func runCsvExport(
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

	err = getRules(client, func(rule globalrules.GlobalRuleRow) (err error) {
		ruleType, err := rule.RuleType.MarshalText()
		if err != nil {
			return
		}
		policy, err := rule.Policy.MarshalText()
		if err != nil {
			return
		}
		record := []string{
			rule.SHA256,
			string(ruleType),
			string(policy),
			rule.CustomMessage,
			rule.Description,
		}
		if err != nil {
			return
		}

		csvRules <- record
		return
	})
	if err != nil {
		return
	}

	close(csvRules)
	wg.Wait()

	return
}

func getRules(client dynamodb.QueryAPI, callback func(globalrules.GlobalRuleRow) error) (err error) {
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
			err = callback(rule)
			if err != nil {
				return
			}
		}

		if nextkey == nil {
			break
		}
		key = nextkey
	}
	return
}
