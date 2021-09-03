package lookup

import (
	"fmt"
	"strings"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/sensordata"
	"github.com/spf13/cobra"
)

func init() {

	var lookupSerialNumberCmd = &cobra.Command{
		Use:   "serial-number",
		Short: "Attempts to search for a machine ID given a serial number and returns the corresponding machine ID",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := cmd.Flags().GetString("region")
			table, _ := cmd.Flags().GetString("dynamodb_table")

			dynamodbClient := dynamodb.GetClient(table, region)

			serialNumber := strings.Join(args, "")

			sensorDataFinderService := sensordata.GetSensorDataFinder(dynamodbClient)

			machineIDs, err := sensorDataFinderService.GetMachineIDsFromSerialNumber(serialNumber, 5)
			if err != nil {
				return err
			}
			fmt.Println(machineIDs)
			return nil
		},
	}

	LookupCmd.AddCommand(lookupSerialNumberCmd)
}
