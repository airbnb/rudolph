package rules

import (
	"fmt"
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
			ddbWriter(client, timeProvider, data, &total)
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
	lines chan map[string]string,
	total *uint64,
) {
	for line := range lines {
		atomic.AddUint64(total, 1)

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
	}
}
