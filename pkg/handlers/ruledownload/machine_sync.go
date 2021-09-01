package ruledownload

import (
	"log"
	"net/http"

	"github.com/airbnb/rudolph/pkg/clock"
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/machinerules"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/model/syncstate"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// On the last page, we always re-send a copy of all machine-specific rules, regardless of whether the client has already received them or not.
// In this way, we ensure that the machine-specific rules take precedence over any other rules in the system.
// FIXME(derek.wang)
//   we don't necessarily have to do it this way; because this makes it really annoying to remove machine rules, as we need complex logic to
//   DELETE them. We could, as an alternative design, simply add machine-specific rules to the feed, and add an extra column to make the rule
//   pertinent only to a single machine. We can use filter expressions to omit them from other machines. However, this design makes it very
//   hard to figure out which rules belong on which machines.
type machineRuleDownloder interface {
	handle(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error)
}

type concreteMachineRuleDownloader struct {
	queryer dynamodb.QueryAPI
	updater dynamodb.UpdateItemAPI
	timer   clock.TimeProvider
}

// On the last page, we always re-send a copy of all machine-specific rules, regardless of whether the client has already received them or not.
// In this way, we ensure that the machine-specific rules take precedence over any other rules in the system.
func (d concreteMachineRuleDownloader) handle(machineID string, ruledownloadRequest *RuledownloadRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("  Ruledownload last page")

	machineRules, err := machinerules.GetMachineRules(d.queryer, machineID)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	// Note that this is the last page
	// Create a sensor sync object to log the FinishedAt time of the rule download process
	err = syncstate.UpdateRuledownloadFinishedAt(d.timer, d.updater, machineID)
	if err != nil {
		log.Printf("Encountered error UpdateItem:")
		log.Print(err.Error())

		return response.APIResponse(http.StatusInternalServerError, err)
	}
	log.Printf("Updated RuledownloadFinishedAt")

	rules := make([]rules.SantaRule, len(*machineRules))
	for i, rule := range *machineRules {
		rules[i] = rule.SantaRule
	}

	// The lack of a cursor in this response signals to the sensor that there is no more stuff to paginate over.
	return response.APIResponse(
		http.StatusOK,
		RuledownloadResponse{
			Rules: DDBRulesToResponseRules(rules),
			// Omit the cursor to signal to the sensor that there are no more pages to paginate through
		},
	)
}
