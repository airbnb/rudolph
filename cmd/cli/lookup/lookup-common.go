package lookup

import (
	"github.com/spf13/cobra"
)

var (
	LookupCmd = &cobra.Command{
		Use:   "lookup",
		Short: "Perform various lookup/search operations on sensordata",
	}
)
