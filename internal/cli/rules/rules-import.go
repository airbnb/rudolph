package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/spf13/cobra"

	"github.com/airbnb/rudolph/internal/csv"
	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	rudolphrules "github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/types"
)

// consts for "workers" in pool for dynamodb operations
const (
	defaultWorkers = 10
	minWorkers     = 1
	maxWorkers     = 2 << defaultWorkers // 2048 default, relative to defaultWorkers
)

func addRuleImportCommand() {
	var filename string
	var workers int

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

			// Try to prevent stupidity
			if workers < minWorkers || workers > maxWorkers {
				fmt.Printf("[WARNING] invalid worker count (%d); using default: %d\n", workers, defaultWorkers)
				workers = defaultWorkers
			}

			return runImport(dynamodbClient, clock.ConcreteTimeProvider{}, filename, workers)
		},
	}

	ruleImportCmd.Flags().StringVarP(&filename, "filename", "f", "", "The filename")
	ruleImportCmd.Flags().IntVarP(&workers, "workers", "w", defaultWorkers, "Number of workers")
	_ = ruleImportCmd.MarkFlagRequired("filename")

	RulesCmd.AddCommand(ruleImportCmd)
}

func runImport(
	client dynamodb.DynamoDBClient,
	timeProvider clock.TimeProvider,
	filename string,
	numWorkers int,
) error {
	if strings.HasSuffix(filename, ".csv") {
		return runCsvImport(client, timeProvider, filename, numWorkers)
	} else if strings.HasSuffix(filename, ".json") {
		return runJsonImport(client, timeProvider, filename, numWorkers)
	}

	return errors.New("unrecognized file extension")
}

func runJsonImport(
	client dynamodb.DynamoDBClient,
	timeProvider clock.TimeProvider,
	filename string,
	numWorkers int,
) (err error) {
	fp, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fp.Close()
	contents, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	var rules []fileRule
	err = json.Unmarshal(contents, &rules)
	if err != nil {
		return
	}

	// Track a total number of lines processed
	// This gets passed to workers and atomic.Add is
	// used to increment in a thread-safe way
	var total uint64

	// Start the workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // ensure Done is called after this worker is complete
			ddbWriter(client, timeProvider, data, &total)
		}()
	}

}

func runCsvImport(
	client dynamodb.DynamoDBClient,
	timeProvider clock.TimeProvider,
	filename string,
	numWorkers int,
) error {
	// ParseCsvFile returns a data channel and an optional error if any issues
	// occurred while opening the file for reading
	data, err := csv.ParseCsvFile(filename)
	if err != nil {
		return err
	}

	// Track a total number of lines processed
	// This gets passed to workers and atomic.Add is
	// used to increment in a thread-safe way
	var total uint64

	// Start the workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // ensure Done is called after this worker is complete

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

			ddbWriter(
				client,
				timeProvider,
				fileRule{
					SHA256:   sha256,
					RuleType: ruleTypeStr,
				},
				&total,
			)
		}()
	}

	// chill
	wg.Wait()

	fmt.Println("processed lines:", total)

	return nil
}

func ddbWriter(
	client dynamodb.DynamoDBClient,
	timeProvider clock.TimeProvider,
	// lines chan map[string]string,
	rules chan fileRule,
	total *uint64,
) {
	for rule := range rules {
		atomic.AddUint64(total, 1)

		var ruleType types.RuleType
		err := ruleType.UnmarshalText([]byte(rule.RuleType))
		if err != nil {
			panic("invalid ruletype")
		}
		var policy types.Policy
		err = policy.UnmarshalText([]byte(rule.Policy))
		if err != nil {
			panic("invalid policy")
		}

		suffix := ""
		if ruleType == types.Certificate {
			suffix = " (Cert)"
		}

		if policy == types.RulePolicyRemove {
			fmt.Printf("  Removing rule: [%s]\n", rule.SHA256)
			sortkey := rudolphrules.RuleSortKeyFromTypeSHA(rule.SHA256, ruleType)
			err = globalrules.RemoveGlobalRule(
				timeProvider,
				client,
				client,
				sortkey,
				"",
			)

		} else {
			fmt.Printf("  Writing rule: [%s] %s%s\n", rule.Policy, rule.SHA256, suffix)
			err = globalrules.AddNewGlobalRule(
				timeProvider,
				client,
				rule.SHA256,
				ruleType,
				policy,
				rule.Description,
			)
		}

		if err != nil {
			panic(err)
		}
	}
}
