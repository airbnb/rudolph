package lookup

import (
	"fmt"
	"strings"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/sensordata"
	"github.com/spf13/cobra"
)

func init() {

	var lookupMachineIDsCmd = &cobra.Command{
		Use:   "machine-ids",
		Short: "Attempts to search for a machine ID given a prefix or entire machine ID",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)

			prefix := strings.Join(args, "")

			sensorDataFinderService := sensordata.GetSensorDataFinder(dynamodbClient)

			machineIDs, err := sensorDataFinderService.GetMachineIDsStartingWith(prefix, 10)
			if err != nil {
				return err
			}
			fmt.Println(machineIDs)
			return nil
		},
	}

	LookupCmd.AddCommand(lookupMachineIDsCmd)
}
