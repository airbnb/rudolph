package ruledownload

import (
	"log"
	"net/http"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/feedrules"
	"github.com/airbnb/rudolph/pkg/model/rules"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// func handleRuleDownloadFromFeed(machineID string, cursor *ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
// 	ddbCursor := cursor.GetLastEvaluatedKey()
// 	globalRules, lastEvaluatedKey, err := modelRules.GetPaginatedFeedRules(ddbCursor, int64(cursor.BatchSize))
// 	if err != nil {
// 		log.Printf("  GetPaginatedGlobalRules Error %s", err.Error())
// 		return response.APIResponse(http.StatusInternalServerError, err)
// 	}

// 	log.Printf("  lastEvaluatedKey %s", lastEvaluatedKey)

// 	nextCursor := cursor.CloneForNextPage()
// 	if lastEvaluatedKey.Empty() {
// 		log.Printf("     No more stuff to paginate over; returning magic cursor")
// 		nextCursor.SetStrategy(ruledownloadStrategyMachine)
// 	} else {
// 		log.Printf("     More stuff to paginate over")
// 		nextCursor.SetDynamodbLastEvaluatedKey(lastEvaluatedKey)
// 	}

// 	return response.APIResponse(
// 		http.StatusOK,
// 		RuledownloadResponse{Rules: DDBRulesToResponseRules(globalRules), Cursor: &nextCursor},
// 	)
// }

type feedRuleDownloader interface {
	handle(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error)
}

type concreteFeedRuleDownloader struct {
	queryer dynamodb.QueryAPI
}

func (d concreteFeedRuleDownloader) handle(machineID string, cursor ruledownloadCursor) (*events.APIGatewayProxyResponse, error) {
	ddbCursor := cursor.GetLastEvaluatedKey()

	feedRules, lastEvaluatedKey, err := feedrules.GetPaginatedFeedRules(d.queryer, cursor.BatchSize, ddbCursor)
	if err != nil {
		log.Printf("  GetPaginatedFeedRules Error %s", err.Error())
		return response.APIResponse(http.StatusInternalServerError, err)
	}

	log.Printf("  lastEvaluatedKey %s", lastEvaluatedKey)

	nextCursor := cursor.CloneForNextPage()
	if lastEvaluatedKey == nil {
		log.Printf("     No more stuff to paginate over; returning magic cursor")
		nextCursor.SetStrategy(ruledownloadStrategyMachine)
	} else {
		log.Printf("     More stuff to paginate over")
		// Here we inherit the preexisting cursor strategy
		nextCursor.SetDynamodbLastEvaluatedKey(lastEvaluatedKey)
	}

	rules := make([]rules.SantaRule, len(feedRules))
	for i, rule := range feedRules {
		rules[i] = rule.SantaRule
	}

	return response.APIResponse(
		http.StatusOK,
		RuledownloadResponse{
			Rules:  DDBRulesToResponseRules(rules),
			Cursor: &nextCursor,
		},
	)
}
