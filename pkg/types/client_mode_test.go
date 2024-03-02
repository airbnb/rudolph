package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientModeTypes_MarshallText(t *testing.T) {
	tests := []struct {
		name       string
		clientMode ClientMode
		want       []byte
		wantErr    bool
	}{
		{"MONITOR", Monitor, []byte("MONITOR"), false},
		{"LOCKDOWN", Lockdown, []byte("LOCKDOWN"), false},
		{"MISSPELLED", ClientMode(0), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.clientMode.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientMode.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// func TestPolicyTypes_UnmarshallText(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		text    []byte
// 		want    Policy
// 		wantErr bool
// 	}{
// 		{"ALLOWLIST", []byte("ALLOWLIST"), RulePolicyAllowlist, false},
// 		{"BLOCKLIST", []byte("BLOCKLIST"), RulePolicyBlocklist, false},
// 		{"SILENT_BLOCKLIST", []byte("SILENT_BLOCKLIST"), RulePolicySilentBlocklist, false},
// 		{"REMOVE", []byte("REMOVE"), RulePolicyRemove, false},
// 		{"ALLOWLIST_COMPILER", []byte("ALLOWLIST_COMPILER"), RulePolicyAllowlistCompiler, false},
// 		{"ALLOWLIST_TRANSITIVE", []byte("ALLOWLIST_TRANSITIVE"), RulePolicyAllowlistTransitive, false},
// 		{"MISSPELLED", []byte("MISSPELLED"), Policy(0), true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got Policy
// 			err := got.UnmarshalText(tt.text)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Policy.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestPolicyTypes_MarshalDynamoDBAttributeValue(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		policy  Policy
// 		want    *dynamodb.AttributeValue
// 		wantErr bool
// 	}{
// 		{"ALLOWLIST", RulePolicyAllowlist, new(dynamodb.AttributeValue).SetN("1"), false},
// 		{"BLOCKLIST", RulePolicyBlocklist, new(dynamodb.AttributeValue).SetN("2"), false},
// 		{"SILENT_BLOCKLIST", RulePolicySilentBlocklist, new(dynamodb.AttributeValue).SetN("3"), false},
// 		{"REMOVE", RulePolicyRemove, new(dynamodb.AttributeValue).SetN("4"), false},
// 		{"ALLOWLIST_COMPILER", RulePolicyAllowlistCompiler, new(dynamodb.AttributeValue).SetN("5"), false},
// 		{"ALLOWLIST_TRANSITIVE", RulePolicyAllowlistTransitive, new(dynamodb.AttributeValue).SetN("6"), false},
// 		{"MISSPELLED", Policy(0), new(dynamodb.AttributeValue), true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			av := &dynamodb.AttributeValue{}
// 			err := tt.policy.MarshalDynamoDBAttributeValue(av)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Policy.MarshalDynamoDBAttributeValue() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, av)
// 		})
// 	}
// }

// func TestPolicyType_UnmarshalDynamoDBAttributeValue(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		av      *dynamodb.AttributeValue
// 		want    Policy
// 		wantErr bool
// 	}{
// 		{"ALLOWLIST", new(dynamodb.AttributeValue).SetN("1"), RulePolicyAllowlist, false},
// 		{"BLOCKLIST", new(dynamodb.AttributeValue).SetN("2"), RulePolicyBlocklist, false},
// 		{"SILENT_BLOCKLIST", new(dynamodb.AttributeValue).SetN("3"), RulePolicySilentBlocklist, false},
// 		{"REMOVE", new(dynamodb.AttributeValue).SetN("4"), RulePolicyRemove, false},
// 		{"ALLOWLIST_COMPILER", new(dynamodb.AttributeValue).SetN("5"), RulePolicyAllowlistCompiler, false},
// 		{"ALLOWLIST_TRANSITIVE", new(dynamodb.AttributeValue).SetN("6"), RulePolicyAllowlistTransitive, false},
// 		{"MISSPELLED", new(dynamodb.AttributeValue), Policy(0), true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := Policy(0)
// 			err := got.UnmarshalDynamoDBAttributeValue(tt.av)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("PolicyType.UnmarshalDynamoDBAttributeValue() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}

// }
