package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

// RuleType represents a Santa rule type.
type RuleType int

const (
	// @deprecated
	Binary = RuleTypeBinary

	// @deprecated
	Certificate = RuleTypeCertificate
)

const (
	// Binary rules use the SHA-256 hash of the entire binary as an identifier.

	RuleTypeBinary RuleType = iota + 1
	// Certificate rules are formed from the SHA-256 fingerprint of an X.509 leaf signing certificate.
	// This is a powerful rule type that has a much broader reach than an individual binary rule .
	// A signing certificate can sign any number of binaries.
	RuleTypeCertificate
)

// UnmarshalText for JSON marshalling interface
func (r *RuleType) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BINARY":
		*r = RuleTypeBinary
	case "CERTIFICATE":
		*r = RuleTypeCertificate
	default:
		return errors.Errorf("unknown rule_type value %q", t)
	}
	return nil
}

// MarshalText for JSON marshalling interface
func (r RuleType) MarshalText() ([]byte, error) {
	switch r {
	case RuleTypeBinary:
		return []byte("BINARY"), nil
	case RuleTypeCertificate:
		return []byte("CERTIFICATE"), nil
	default:
		return nil, errors.Errorf("unknown rule_type %d", r)
	}
}

// MarshalDynamoDBAttributeValue for ddb
func (r RuleType) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	var s string
	switch r {
	case RuleTypeBinary:
		// s = "BINARY"
		s = "1"
	case RuleTypeCertificate:
		// s = "CERTIFICATE"
		s = "2"
	default:
		return errors.Errorf("unknown rule_type value %q", r)
	}
	// av.S = &s
	av.N = &s
	return nil
}

// UnmarshalDynamoDBAttributeValue implements the Unmarshaler interface
func (r *RuleType) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	// switch t := aws.StringValue(av.S); t {
	switch t := aws.StringValue(av.N); t {
	case "1":
		fallthrough
	case "BINARY":
		*r = RuleTypeBinary
	case "2":
		fallthrough
	case "CERTIFICATE":
		*r = RuleTypeCertificate
	default:
		return errors.Errorf("unknown rule_type value %q", t)
	}

	return nil
}
