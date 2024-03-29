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
	rulesBuffer := make(chan fileRule)
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // ensure Done is called after this worker is complete
			ddbWriter(
				client,
				timeProvider,
				rulesBuffer,
				&total,
			)
		}()
	}

	// Shovel all the json-parsed rules into the worker queue
	for _, rule := range rules {
		rulesBuffer <- rule
	}
	close(rulesBuffer)

	// Chill
	wg.Wait()

	fmt.Println("processed lines:", total)

	return
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

	// Channel for csv parsing and workers to communicate over
	rules := make(chan fileRule)

	// Track a total number of lines processed
	// This gets passed to workers and atomic.Add is
	// used to increment in a thread-safe way
	var total uint64

	// Start the workers
	// Fanning out workers allows us to make multiple HTTP requests concurrently which can
	// improve performance assuming we aren't network I/O bottlenecked or something.
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // ensure Done is called after this worker is complete
			ddbWriter(
				client,
				timeProvider,
				rules,
				&total,
			)
		}()
	}

	// Start taking lines from the csv and shoveling them into the workers
	for line := range data {
		identifier, ok := line["identifier"]
		if !ok {
			panic("no identifier")
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
		customMsg, ok := line["custom_msg"]
		if !ok {
			customMsg = ""
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

		rules <- fileRule{
			Identifier:    identifier,
			RuleType:      ruleType,
			Policy:        policy,
			Description:   description,
			CustomMessage: customMsg,
		}
	}
	close(rules)

	// chill
	wg.Wait()

	fmt.Println("processed lines:", total)

	return nil
}

func ddbWriter(
	client dynamodb.DynamoDBClient,
	timeProvider clock.TimeProvider,
	rules chan fileRule,
	total *uint64,
) {
	for rule := range rules {
		var err error
		atomic.AddUint64(total, 1)

		var suffix string
		switch rule.RuleType {
		case types.RuleTypeCertificate:
			suffix = " (Cert)"
		case types.RuleTypeTeamID:
			suffix = " (TeamID)"
		case types.RuleTypeSigningID:
			suffix = " (SigningID)"
		default:
			suffix = ""
		}

		if rule.Policy == types.RulePolicyRemove {
			fmt.Printf("  Removing rule: [%s]\n", rule.Identifier)
			sortkey := rudolphrules.RuleSortKeyFromTypeIdentifier(rule.Identifier, rule.RuleType)
			err = globalrules.RemoveGlobalRule(
				timeProvider,
				client,
				client,
				sortkey,
				"",
			)

		} else {
			fmt.Printf("  Writing rule: [%+v] %s%s\n", rule.Policy, rule.Identifier, suffix)
			err = globalrules.AddNewGlobalRule(
				timeProvider,
				client,
				rule.Identifier,
				rule.RuleType,
				rule.Policy,
				rule.Description,
			)
		}

		if err != nil {
			panic(err)
		}
	}
}
