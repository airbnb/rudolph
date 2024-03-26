package types

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	// 	Most Specific                                  Least Specific
	// Binary   -->   Signing ID   -->   Certificate   -->   Team ID

	// Binary rules use the SHA-256 hash of the entire binary as an identifier.
	RuleTypeBinary RuleType = iota + 1

	// Certificate rules are formed from the SHA-256 fingerprint of an X.509 leaf signing certificate.
	// This is a powerful rule type that has a much broader reach than an individual binary rule .
	// A signing certificate can sign any number of binaries.
	RuleTypeCertificate

	// SigningID rules are arbitrary identifiers under developer control that are given to a binary at signing time.
	// Typically, these use reverse domain name notation and include the name of the binary (e.g. com.google.Chrome).
	// Because the signing IDs are arbitrary, the Santa rule identifier must be prefixed with the Team ID associated with the Apple developer certificate used to sign the application.
	// For example, a signing ID rule for Google Chrome would be: EQHXZ8M8AV:com.google.Chrome.
	//For platform binaries (i.e. those binaries shipped by Apple with the OS) which do not have a Team ID, the string platform must be used (e.g. platform:com.apple.curl).
	RuleTypeSigningID

	// TeamID rules are formed from the Apple Developer Program Team ID is a 10-character identifier issued by Apple and tied to developer accounts/organizations.
	// This is distinct from Certificates, as a single developer account can and frequently will request/rotate between multiple different signing certificates and entitlements.
	// This is an even more powerful rule with broader reach than individual certificate rules.
	RuleTypeTeamID
)

// UnmarshalText for JSON marshalling interface
func (r *RuleType) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BINARY":
		*r = RuleTypeBinary
	case "CERTIFICATE":
		*r = RuleTypeCertificate
	case "SIGNINGID":
		*r = RuleTypeSigningID
	case "TEAMID":
		*r = RuleTypeTeamID
	default:
		return fmt.Errorf("unknown rule_type value %q", t)
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
	case RuleTypeSigningID:
		return []byte("SIGNINGID"), nil
	case RuleTypeTeamID:
		return []byte("TEAMID"), nil
	default:
		return nil, fmt.Errorf("unknown rule_type %d", r)
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
	case RuleTypeSigningID:
		s = "3"
	case RuleTypeTeamID:
		s = "4"
	default:
		return fmt.Errorf("unknown rule_type value %q", r)
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
	case "3":
		fallthrough
	case "SIGNINGID":
		*r = RuleTypeSigningID
	case "4":
		fallthrough
	case "TEAMID":
		*r = RuleTypeTeamID
	default:
		return fmt.Errorf("unknown rule_type value %q", t)
	}

	return nil
}
