package types

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestRuleType_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		rule    RuleType
		want    []byte
		wantErr bool
	}{
		{"Binary", RuleTypeBinary, []byte("BINARY"), false},
		{"Certificate", RuleTypeCertificate, []byte("CERTIFICATE"), false},
		{"SigningID", RuleTypeSigningID, []byte("SIGNINGID"), false},
		{"TeamID", RuleTypeTeamID, []byte("TEAMID"), false},
		{"Invalid", RuleType(0), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rule.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("RuleType.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRuleType_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    RuleType
		wantErr bool
	}{
		{"Binary", []byte("BINARY"), RuleTypeBinary, false},
		{"Certificate", []byte("CERTIFICATE"), RuleTypeCertificate, false},
		{"SigningID", []byte("SIGNINGID"), RuleTypeSigningID, false},
		{"TeamID", []byte("TEAMID"), RuleTypeTeamID, false},
		{"Invalid", []byte("INVALID"), RuleType(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got RuleType
			err := got.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("RuleType.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRuleType_MarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []struct {
		name    string
		rule    RuleType
		want    *dynamodb.AttributeValue
		wantErr bool
	}{
		{"BINARY", RuleTypeBinary, new(dynamodb.AttributeValue).SetN("1"), false},
		{"CERTIFICATE", RuleTypeCertificate, new(dynamodb.AttributeValue).SetN("2"), false},
		{"SIGNINGID", RuleTypeSigningID, new(dynamodb.AttributeValue).SetN("3"), false},
		{"TEAMID", RuleTypeTeamID, new(dynamodb.AttributeValue).SetN("4"), false},
		{"INVALID", RuleType(0), new(dynamodb.AttributeValue), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := &dynamodb.AttributeValue{}
			err := tt.rule.MarshalDynamoDBAttributeValue(av)
			if (err != nil) != tt.wantErr {
				t.Errorf("RuleType.MarshalDynamoDBAttributeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, av)
		})
	}
}

func TestRuleType_UnmarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []struct {
		name    string
		av      *dynamodb.AttributeValue
		want    RuleType
		wantErr bool
	}{
		{"BINARY", new(dynamodb.AttributeValue).SetN("1"), RuleTypeBinary, false},
		{"CERTIFICATE", new(dynamodb.AttributeValue).SetN("2"), RuleTypeCertificate, false},
		{"SIGNINGID", new(dynamodb.AttributeValue).SetN("3"), RuleTypeSigningID, false},
		{"TEAMID", new(dynamodb.AttributeValue).SetN("4"), RuleTypeTeamID, false},
		{"INVALID", new(dynamodb.AttributeValue).SetN("0"), RuleType(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuleType(0)
			err := got.UnmarshalDynamoDBAttributeValue(tt.av)
			if (err != nil) != tt.wantErr {
				t.Errorf("RuleType.UnmarshalDynamoDBAttributeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
