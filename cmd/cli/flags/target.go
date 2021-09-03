package flags

import (
	"errors"

	"github.com/airbnb/rudolph/cmd/cli/santa_sensor"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/spf13/cobra"
)

type RuleTargetFlags struct {
	RuleTypeArg  types.RuleType
	IsGlobal     bool
	MachineIdArg string
	SHA256Arg    string
	FilepathArg  string
}

type TargetFlags struct {
	MachineID     string
	IsGlobal      bool
	SelfMachineID string
}

func (t *TargetFlags) AddTargetFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&t.IsGlobal, "global", "g", false, "Apply globally. Provide one of [--global|--machine], or exclude both to apply to current machine.")
	cmd.Flags().StringVarP(&t.MachineID, "machine", "m", "", `The uuid of the machine.`)
}

func (t *TargetFlags) AddTargetFlagsRules(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&t.IsGlobal, "global", "g", false, "Retrive rules that apply globally.")
	cmd.Flags().StringVarP(&t.MachineID, "machine", "m", "", "Retrieve rules for a single machine. Omit to apply to the current machine.")

}

func (t TargetFlags) GetMachineID() (string, error) {
	t.getSelfMachineID() // We do some ninja initialization

	if t.IsGlobal {
		return "", errors.New("do not call GetMachineID when IsGlobal is true")
	} else if t.MachineID != "" {
		return t.MachineID, nil
	} else {
		return santa_sensor.GetSelfMachineID()
	}
}

// WARNING: This is best-effort as if you don't call getSelfMachineID first, it'll just always return false.
func (t TargetFlags) IsTargetSelf() bool {
	return t.SelfMachineID == t.MachineID
}

func (t *TargetFlags) getSelfMachineID() (string, error) {
	if t.SelfMachineID != "" {
		return t.SelfMachineID, nil
	}

	selfMachineId, err := santa_sensor.GetSelfMachineID()
	if err != nil {
		return "", err
	}
	t.SelfMachineID = selfMachineId
	return t.SelfMachineID, nil
}
