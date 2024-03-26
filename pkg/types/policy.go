package types

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Policy represents the Santa Rule Policy.
type Policy int

const (
	// @deprecated
	Allowlist           = RulePolicyAllowlist
	Blocklist           = RulePolicyBlocklist
	SilentBlocklist     = RulePolicySilentBlocklist
	Remove              = RulePolicyRemove
	AllowlistCompiler   = RulePolicyAllowlistCompiler
	AllowlistTransitive = RulePolicyAllowlistTransitive
)

const (
	RulePolicyAllowlist Policy = iota + 1
	RulePolicyBlocklist
	RulePolicySilentBlocklist
	// Remove is a "special" rule in that, when it is sent by the server, it instructs the sensor
	// to delete any associated rule.
	RulePolicyRemove
	// AllowlistCompiler is a Transitive Allowlist policy which allows binaries created by
	// a specific compiler. EnabledTransitiveRules must be set to true in the Preflight first.
	RulePolicyAllowlistCompiler
	// Transitive rules are created by the santa sensor itself; it is never created by the server.
	// Transitive rules are destroyed upon every clean sync.
	RulePolicyAllowlistTransitive
)

// UnmarshalText for JSON marshalling interface
// Use Santa defined constants
// https://github.com/google/santa/blob/main/Source/santactl/Commands/sync/SNTCommandSyncConstants.m#L98-L109
func (p *Policy) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "ALLOWLIST":
		*p = RulePolicyAllowlist
	case "BLOCKLIST":
		*p = RulePolicyBlocklist
	case "SILENT_BLOCKLIST":
		*p = RulePolicySilentBlocklist
	case "REMOVE":
		*p = RulePolicyRemove
	case "ALLOWLIST_COMPILER":
		*p = RulePolicyAllowlistCompiler
	case "ALLOWLIST_TRANSITIVE":
		*p = RulePolicyAllowlistTransitive
	default:
		return fmt.Errorf("unknown policy value %q", t)
	}
	return nil
}

// MarshalText for JSON marshalling interface
func (p Policy) MarshalText() ([]byte, error) {
	switch p {
	case RulePolicyAllowlist:
		return []byte("ALLOWLIST"), nil
	case RulePolicyBlocklist:
		return []byte("BLOCKLIST"), nil
	case RulePolicySilentBlocklist:
		return []byte("SILENT_BLOCKLIST"), nil
	case RulePolicyRemove:
		return []byte("REMOVE"), nil
	case RulePolicyAllowlistCompiler:
		return []byte("ALLOWLIST_COMPILER"), nil
	case RulePolicyAllowlistTransitive:
		return []byte("ALLOWLIST_TRANSITIVE"), nil
	default:
		return nil, fmt.Errorf("unknown policy %d", p)
	}
}

// MarshalDynamoDBAttributeValue for ddb
func (p Policy) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	var s string
	switch p {
	case RulePolicyAllowlist:
		s = "1"
	case RulePolicyBlocklist:
		s = "2"
	case RulePolicySilentBlocklist:
		s = "3"
	case RulePolicyRemove:
		s = "4"
	case RulePolicyAllowlistCompiler:
		s = "5"
	case RulePolicyAllowlistTransitive:
		s = "6"
	default:
		return fmt.Errorf("unknown policy value %q", p)
	}
	// av.S = &s
	av.N = &s
	return nil
}

// UnmarshalDynamoDBAttributeValue implements the Unmarshaler interface
func (p *Policy) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	switch t := aws.StringValue(av.N); t {
	case "1":
		fallthrough
	case "ALLOWLIST":
		*p = RulePolicyAllowlist
	case "2":
		fallthrough
	case "BLOCKLIST":
		*p = RulePolicyBlocklist
	case "3":
		fallthrough
	case "SILENT_BLOCKLIST":
		*p = RulePolicySilentBlocklist
	case "4":
		fallthrough
	case "REMOVE":
		*p = RulePolicyRemove
	case "5":
		fallthrough
	case "ALLOWLIST_COMPILER":
		*p = RulePolicyAllowlistCompiler
	case "6":
		fallthrough
	case "ALLOWLIST_TRANSITIVE":
		*p = RulePolicyAllowlistTransitive
	default:
		return fmt.Errorf("unknown policy value %q", t)
	}

	return nil
}
