package repair

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/scan"
)

var (
	RepairCmd *cobra.Command
)

func init() {
	RepairCmd = &cobra.Command{
		Use:   "repair",
		Short: "Scans the entire DynamoDB database and repairs broken records",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)
			scanService := scan.GetScanService(dynamodbClient)

			return repair(scanService, dynamodbClient)
		},
	}
}

func repair(scanService scan.ScanService, deleter dynamodb.DeleteItemAPI) error {

	// Configuration
	pageSize := int32(10)
	delete := false

	input := awsdynamodb.ScanInput{
		ConsistentRead: aws.Bool(false),
		Limit:          aws.Int32(pageSize),
	}

	page := 0
	callback := func(out *awsdynamodb.ScanOutput) (err error) {
		page += 1
		fmt.Println("==========")
		fmt.Printf("Page #%d\n", page)
		fmt.Printf("  Scanned %d items\n", out.ScannedCount)
		fmt.Printf("  Discovered %d items\n", out.Count)

		var items []dynamodbItem
		err = attributevalue.UnmarshalListOfMaps(out.Items, &items)
		if err != nil {
			err = fmt.Errorf("failed to unmarshal output items: %w", err)
			return
		}

		for _, item := range items {
			if strings.HasPrefix(item.PartitionKey, "MachineInfo#") || strings.HasPrefix(item.PartitionKey, "MachineConfig#") {
				fmt.Printf("    ++ LEGACY ITEM LOCATED: (%+v | %+v) ++\n", item.PartitionKey, item.SortKey)
				if delete {
					_, err := deleter.DeleteItem(item.PrimaryKey)
					if err != nil {
						return fmt.Errorf("failed to delete item: %w", err)
					}
					fmt.Printf("      Item Deleted\n")
				}
			}
		}

		fmt.Println("")
		return nil
	}

	stop := func(out *awsdynamodb.ScanOutput) (bool, error) {
		// Don't stop until done.
		return false, nil
	}

	err := scanService.ScanAll(input, callback, stop)
	if err != nil {
		return err
	}
	return nil
}

type dynamodbItem struct {
	dynamodb.PrimaryKey
}
