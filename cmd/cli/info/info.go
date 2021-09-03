package info

import (
	"fmt"

	"github.com/airbnb/rudolph/cmd/cli/santa_sensor"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	InfoCmd *cobra.Command
)

func init() {
	InfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Get information about your machine",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return info()
		},
	}
}

func info() error {
	// FIXME (derek.wang) allow more machine ids
	machineID, err := santa_sensor.GetSelfMachineID()
	if err != nil {
		return errors.Wrap(err, "Failed to get MachineUUID!")
	}

	fmt.Println("Your machineUUID is: ", machineID)

	// FIXME (derek.wang)
	// Query the remote config and check what the database thinks it should be

	return nil
}
