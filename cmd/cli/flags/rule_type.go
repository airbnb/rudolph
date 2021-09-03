package flags

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/types"
)

const (
	binType           = "binary"
	binTypeShort      = "bin"
	certType          = "certificate"
	certTypeShort     = "cert"
	monitorMode       = "monitor"
	lockdownMode      = "lockdown"
	defaultClientMode = "monitor"
)

// ruleType is a custom type for use as a CLI flag representing the type of rule being applied
type RuleType types.RuleType

func newRuleTypeValue(val string, p *RuleType) *RuleType {
	err := p.Set(val)
	if err != nil {
		fmt.Println(`Warning: invalid default value for rule-type, using "binary"`)
		*p = RuleType(types.Binary)
	}

	return p
}

func (i *RuleType) AsRuleType() types.RuleType {
	return types.RuleType(*i)
}

func (i *RuleType) Set(s string) error {
	switch s {
	case binType, binTypeShort:
		*i = RuleType(types.RuleTypeBinary)
	case certType, certTypeShort:
		*i = RuleType(types.RuleTypeCertificate)
	default:
		return fmt.Errorf(`invalid rule type; must be "binary" or "cert"`)
	}
	return nil
}

func (i *RuleType) Type() string {
	return "string"
}

func (i *RuleType) String() string {
	v := (types.RuleType)(*i)
	switch v {
	case types.RuleTypeBinary:
		return binType
	case types.RuleTypeCertificate:
		return certType
	}

	// No default
	return ""
}
