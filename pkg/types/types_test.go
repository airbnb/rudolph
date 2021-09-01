package types

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

var (
	av dynamodb.AttributeValue
)

func TestTypesClientMode_Marshal(t *testing.T) {
	t.Run("Monitor", func(t *testing.T) {
		clientMode, err := ClientMode(1).MarshalText()
		assert.Empty(t, err)
		assert.Equal(t, clientMode, []byte("MONITOR"))
	})

	t.Run("Lockdown", func(t *testing.T) {
		clientMode, err := ClientMode(2).MarshalText()
		assert.Empty(t, err)
		assert.Equal(t, clientMode, []byte("LOCKDOWN"))
	})

	t.Run("Error", func(t *testing.T) {
		_, err := ClientMode(3).MarshalText()
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, "unknown client_mode 3")
	})
}

func TestTypesClientMode_Unmarshal(t *testing.T) {
	t.Run("Monitor", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("MONITOR"))
		assert.Empty(t, err)
		assert.Equal(t, clientMode, ClientMode(1))
	})

	t.Run("Lockdown", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("LOCKDOWN"))
		assert.Empty(t, err)
		assert.Equal(t, clientMode, ClientMode(2))
	})

	t.Run("ERROR", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("NOPE"))
		assert.NotEmpty(t, err)
	})

	t.Run("ERROR MONITORS", func(t *testing.T) {
		var clientMode ClientMode
		err := clientMode.UnmarshalText([]byte("MONITORS"))
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, `unknown client_mode value "MONITORS"`)
	})
}

func TestTypesRuleType_Marshal(t *testing.T) {
	t.Run("Binary", func(t *testing.T) {
		ruleType, err := RuleType.MarshalText(1)
		assert.Empty(t, err)
		assert.Equal(t, ruleType, []byte("BINARY"))
	})

	t.Run("Certificate", func(t *testing.T) {
		ruleType, err := RuleType.MarshalText(2)
		assert.Empty(t, err)
		assert.Equal(t, ruleType, []byte("CERTIFICATE"))
	})

	t.Run("ERROR", func(t *testing.T) {
		_, err := RuleType.MarshalText(3)
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, "unknown rule_type 3")
	})
}

func TestTypesRuleType_Unmarshal(t *testing.T) {
	t.Run("BINARY", func(t *testing.T) {
		var ruleType RuleType
		err := ruleType.UnmarshalText([]byte("BINARY"))
		assert.Empty(t, err)
		assert.Equal(t, ruleType, RuleType(1))
	})

	t.Run("CERTIFICATE", func(t *testing.T) {
		var ruleType RuleType
		err := ruleType.UnmarshalText([]byte("CERTIFICATE"))
		assert.Empty(t, err)
		assert.Equal(t, ruleType, RuleType(2))
	})

	t.Run("ERROR MISSPELLED", func(t *testing.T) {
		var ruleType RuleType
		err := ruleType.UnmarshalText([]byte("CERTIFICATES"))
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, `unknown rule_type value "CERTIFICATES"`)
	})
}

func TestTypesPolicy_Marshall(t *testing.T) {
	t.Run("ALLOWLIST", func(t *testing.T) {
		policy, err := Policy.MarshalText(1)
		assert.Empty(t, err)
		assert.Equal(t, policy, []byte("ALLOWLIST"))
	})

	t.Run("BLOCKLIST", func(t *testing.T) {
		policy, err := Policy.MarshalText(2)
		assert.Empty(t, err)
		assert.Equal(t, policy, []byte("BLOCKLIST"))
	})

	t.Run("SILENT_BLOCKLIST", func(t *testing.T) {
		policy, err := Policy.MarshalText(3)
		assert.Empty(t, err)
		assert.Equal(t, policy, []byte("SILENT_BLOCKLIST"))
	})

	t.Run("REMOVE", func(t *testing.T) {
		policy, err := Policy.MarshalText(4)
		assert.Empty(t, err)
		assert.Equal(t, policy, []byte("REMOVE"))
	})

	t.Run("ALLOWLIST_COMPILER", func(t *testing.T) {
		policy, err := Policy.MarshalText(5)
		assert.Empty(t, err)
		assert.Equal(t, policy, []byte("ALLOWLIST_COMPILER"))
	})

	t.Run("ALLOWLIST_TRANSITIVE", func(t *testing.T) {
		policy, err := Policy.MarshalText(6)
		assert.Empty(t, err)
		assert.Equal(t, policy, []byte("ALLOWLIST_TRANSITIVE"))
	})

	t.Run("ERROR MISSPELLED", func(t *testing.T) {
		_, err := Policy.MarshalText(7)
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, "unknown policy 7")
	})

}

func TestTypesPolicy_Unmarshal(t *testing.T) {
	t.Run("ALLOWLIST", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("ALLOWLIST"))
		assert.Empty(t, err)
		assert.Equal(t, policy, Policy(1))
	})

	t.Run("BLOCKLIST", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("BLOCKLIST"))
		assert.Empty(t, err)
		assert.Equal(t, policy, Policy(2))
	})

	t.Run("SILENT_BLOCKLIST", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("SILENT_BLOCKLIST"))
		assert.Empty(t, err)
		assert.Equal(t, policy, Policy(3))
	})

	t.Run("REMOVE", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("REMOVE"))
		assert.Empty(t, err)
		assert.Equal(t, policy, Policy(4))
	})

	t.Run("ALLOWLIST_COMPILER", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("ALLOWLIST_COMPILER"))
		assert.Empty(t, err)
		assert.Equal(t, policy, Policy(5))
	})

	t.Run("ALLOWLIST_TRANSITIVE", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("ALLOWLIST_TRANSITIVE"))
		assert.Empty(t, err)
		assert.Equal(t, policy, Policy(6))
	})

	t.Run("ERROR MISSPELLED", func(t *testing.T) {
		var policy Policy
		err := policy.UnmarshalText([]byte("ALLOWLISTS"))
		assert.NotEmpty(t, err)
		assert.EqualError(t, err, `unknown policy value "ALLOWLISTS"`)
	})
}

func TestTypesPolicy_MarshalDynamoDBAttributeValue_Allowlist(t *testing.T) {
	var policy Policy = 1

	err := policy.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("1")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesPolicy_MarshalDynamoDBAttributeValue_Blocklist(t *testing.T) {
	var policy Policy = 2
	err := policy.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("2")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesPolicy_MarshalDynamoDBAttributeValue_Silent_Blocklist(t *testing.T) {
	var policy Policy = 3

	err := policy.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("3")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}
func TestTypesPolicy_MarshalDynamoDBAttributeValue_Remove(t *testing.T) {
	var av dynamodb.AttributeValue
	var policy Policy = 4

	err := policy.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("4")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesPolicy_MarshalDynamoDBAttributeValue_Allowlist_Complier(t *testing.T) {
	var policy Policy = 5

	err := policy.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("5")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesPolicy_MarshalDynamoDBAttributeValue_Allowlist_Transitive(t *testing.T) {
	var policy Policy = 6

	err := policy.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("6")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesPolicy_MarshalDynamoDBAttributeValue_ErrorsOut(t *testing.T) {
	var policy Policy = 7

	err := policy.MarshalDynamoDBAttributeValue(&av)
	//expectAv := dynamodb.AttributeValue{N: aws.String("7")}
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, fmt.Sprintf("%s%q", "unknown policy value ", policy))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_Allowlist(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("1")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, policy, Policy(1))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_Blocklist(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("2")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, policy, Policy(2))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_Silent_Blocklist(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("3")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, policy, Policy(3))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_Remove(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("4")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, policy, Policy(4))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_Allowlist_Complier(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("5")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, policy, Policy(5))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_Allowlist_Transitive(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("6")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, policy, Policy(6))
}

func TestTypesPolicy_UnmarshalDynamoDBAttributeValue_ErrorsOut(t *testing.T) {
	var av dynamodb.AttributeValue
	av.SetN("7")
	var policy Policy
	err := policy.UnmarshalDynamoDBAttributeValue(&av)
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, `unknown policy value "7"`)
}

func TestTypesRuleTypes_UnmarshalDynamoDBAttributeValue_Binary(t *testing.T) {
	var ruleType RuleType
	av.SetN("1")

	err := ruleType.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, ruleType, RuleType(1))

	av.SetN("BINARY")

	err = ruleType.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, ruleType, RuleType(1))
}

func TestTypesRuleTypes_UnmarshalDynamoDBAttributeValue_Certificate(t *testing.T) {
	var ruleType RuleType
	av.SetN("2")

	err := ruleType.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, ruleType, RuleType(2))

	av.SetN("CERTIFICATE")

	err = ruleType.UnmarshalDynamoDBAttributeValue(&av)
	assert.Empty(t, err)
	assert.Equal(t, ruleType, RuleType(2))
}

func TestTypesRuleTypes_UnmarshalDynamoDBAttributeValue_ErrorsOut(t *testing.T) {
	var ruleType RuleType
	av.SetN("CERTIFICATESS")

	err := ruleType.UnmarshalDynamoDBAttributeValue(&av)
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, `unknown rule_type value "CERTIFICATESS"`)
}

func TestTypesRuleTypes_MarshalDynamoDBAttributeValue_Binary(t *testing.T) {
	var ruleType RuleType = 1

	err := ruleType.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("1")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesRuleTypes_MarshalDynamoDBAttributeValue_Certificate(t *testing.T) {
	var ruleType RuleType = 2

	err := ruleType.MarshalDynamoDBAttributeValue(&av)
	expectAv := dynamodb.AttributeValue{N: aws.String("2")}
	assert.Empty(t, err)
	assert.Equal(t, av, expectAv)
}

func TestTypesRuleTypes_MarshalDynamoDBAttributeValue_ErrorsOut(t *testing.T) {
	var ruleType RuleType = 3

	err := ruleType.MarshalDynamoDBAttributeValue(&av)
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, fmt.Sprintf("%s%q", "unknown rule_type value ", ruleType))
}
