package machinerules

import (
	"time"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
)

// @deprecated
type MachineRulesUpdater interface {
	UpdateMachineRulePolicy(machineID string, sha256 string, ruleType types.RuleType, rulePolicy types.Policy) error
}

// @deprecated
type ConcreteMachineRulesUpdater struct {
	Updater      dynamodb.UpdateItemAPI
	TimeProvider clock.TimeProvider
}

// @deprecated
func (c ConcreteMachineRulesUpdater) UpdateMachineRulePolicy(machineID string, sha256 string, ruleType types.RuleType, rulePolicy types.Policy) error {
	expires := c.TimeProvider.Now().Add(time.Hour * MachineRuleDefaultExpirationHours).UTC()

	return UpdateMachineRule(c.Updater, machineID, sha256, ruleType, rulePolicy, expires)
}

//
// This service exposes all machine rules access methods
//
type MachineRulesService interface {
	Get(machineId string, sha256 string, ruleType types.RuleType) (rule *MachineRuleRow, err error)
	Add(machineId string, sha256 string, ruleType types.RuleType, policy types.Policy, description string, expires time.Time) error
	Update(machineId string, sha256 string, ruleType types.RuleType, rulePolicy types.Policy, expires time.Time) error
	Remove(machineId string, sha256 string, ruleType types.RuleType) error
	RemoveBySortKey(machineId string, ruleSortKey string) error
	GetMachineRules(machineID string) (items *[]MachineRuleRow, err error)
}

type ConcreteMachineRulesService struct {
	dynamodb dynamodb.DynamoDBClient
}

func GetMachineRulesService(dynamodb dynamodb.DynamoDBClient) MachineRulesService {
	return ConcreteMachineRulesService{
		dynamodb: dynamodb,
	}
}

func (s ConcreteMachineRulesService) Get(machineId string, sha256 string, ruleType types.RuleType) (rule *MachineRuleRow, err error) {
	return getItemAsMachineRule(s.dynamodb, machineRulePK(machineId), machineRuleSK(sha256, ruleType))
}
func (s ConcreteMachineRulesService) Add(machineId string, sha256 string, ruleType types.RuleType, policy types.Policy, description string, expires time.Time) error {
	return AddNewMachineRule(s.dynamodb, machineId, sha256, ruleType, policy, description, expires)
}
func (s ConcreteMachineRulesService) Update(machineId string, sha256 string, ruleType types.RuleType, rulePolicy types.Policy, expires time.Time) error {
	return UpdateMachineRule(s.dynamodb, machineId, sha256, ruleType, rulePolicy, expires)
}
func (s ConcreteMachineRulesService) RemoveBySortKey(machineId string, ruleSortKey string) error {
	return RemoveMachineRule(s.dynamodb, s.dynamodb, machineId, ruleSortKey)
}
func (s ConcreteMachineRulesService) Remove(machineId string, sha256 string, ruleType types.RuleType) error {
	return RemoveMachineRule(s.dynamodb, s.dynamodb, machineId, machineRuleSK(sha256, ruleType))
}
func (s ConcreteMachineRulesService) GetMachineRules(machineId string) (items *[]MachineRuleRow, err error) {
	return GetMachineRules(s.dynamodb, machineId)
}
