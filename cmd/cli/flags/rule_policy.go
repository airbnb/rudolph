package flags

import (
	"fmt"

	"strings"

	"github.com/airbnb/rudolph/pkg/types"
	"github.com/spf13/cobra"
)

const (
	allowlist               = "allowlist"
	allowlistShort          = "allow"
	blocklist               = "blocklist"
	blocklistShort          = "block"
	silentBlocklist         = "silent_blocklist"
	silentBlocklistShort    = "silent_block"
	silentBlocklistShortAlt = "silent-block"
	remove                  = "remove"
	removeShort             = "delete"
)

var (
	rulePolicyArg RulePolicy
)

type RuleUpdateFlags struct {
	RulePolicy *RulePolicy
}

func (r *RuleUpdateFlags) AddRuleUpdateFlags(cmd *cobra.Command) {

	// rule-policy is to specify the policy for edit commands
	cmd.Flags().VarP(&rulePolicyArg, "rule-policy", "p", `type of rule being applied. valid options are: "allowlist", "blocklist", "silent_blocklist" or "remove"`)
	_ = cmd.MarkFlagRequired("rule-policy")

	r.RulePolicy = &rulePolicyArg
}

// rulePolicy is a custom type for use as a CLI flag representing the type of rule policy being applied
type RulePolicy types.Policy

func newRulePolicyValue(val string, p *RuleType) *RuleType {
	err := p.Set(val)
	if err != nil {
		fmt.Println(`Warning: invalid default value for rule policy, using "allowlist"`)
		*p = RuleType(types.Allowlist)
	}

	return p
}

func (i *RulePolicy) AsRulePolicy() types.Policy {
	return types.Policy(*i)
}

func (i *RulePolicy) Set(s string) error {
	s = strings.ToLower(s)
	switch s {
	case allowlist, allowlistShort:
		*i = RulePolicy(types.Allowlist)
	case blocklist, blocklistShort:
		*i = RulePolicy(types.Blocklist)
	case silentBlocklist, silentBlocklistShort, silentBlocklistShortAlt:
		*i = RulePolicy(types.SilentBlocklist)
	case remove, removeShort:
		*i = RulePolicy(types.Remove)
	default:
		return fmt.Errorf(`invalid rule policy; must be "allowlist", "blocklist", "silent_blocklist" or "remove"`)
	}
	return nil
}

func (i *RulePolicy) Type() string {
	return "string"
}

func (i *RulePolicy) String() string {
	v := (types.Policy)(*i)
	switch v {
	case types.Allowlist:
		return allowlist
	case types.Blocklist:
		return blocklist
	case types.SilentBlocklist:
		return blocklist
	case types.Remove:
		return remove
	}

	// No default
	return ""
}
