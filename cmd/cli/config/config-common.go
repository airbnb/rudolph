package config

import (
	"github.com/spf13/cobra"
)

// config specific variables
// isGlobal, machineIDArg reside in parser.go

var (
	ConfigCmd *cobra.Command
)

func init() {
	ConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Perform various config operations",
	}
}
