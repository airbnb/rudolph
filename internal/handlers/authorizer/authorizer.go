package authorizer

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

var (
	errBlankRequestBody = errors.New("Request body is blank")
)

type authorizerEnvironment struct {
	Region    string
	GatewayID string
	AccountID string
	StageName string
}

var (
	authorizerEnv authorizerEnvironment
)

func init() {
	authorizerEnv = authorizerEnvironment{
		Region:    os.Getenv("REGION"),
		GatewayID: os.Getenv("GATEWAY_ID"),
		AccountID: os.Getenv("ACCOUNT_ID"),
		StageName: os.Getenv("STAGE_NAME"),
	}

}

// HandleAuthorizerRequest is the handler to be used by the authorizer function
func HandleAuthorizerRequest(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
	log.Println("Custom Lambda Authorizer - Method ARN: " + request.MethodArn)

	if request.HTTPMethod == http.MethodGet && request.Path == "/health" {
		return allowResponse(request, "HEALTH_CHECK"), nil
	}

	if request.HTTPMethod != http.MethodPost {
		return denyResponse(request, "Incorrect Method"), nil
	}

	machineID, ok := request.PathParameters["machine_id"]
	if !ok {
		return denyResponse(request, "Incorrect Request URI"), nil
	}

	// TODO: FILL ME IN
	//   Here you can bring your own authorization policies. Below are several examples.
	//   There's currently no authentication headers that the santa sensor sends to the server that can be used in
	//   AuthN or AuthZ, so this is more or less a best-effort + BYOB "authentication" system

	// Restrict to only santactl useragent
	// Notably, the useragent can be faked so this isn't a durable security check
	// if request.RequestContext.Identity.UserAgent != "santactl-sync/2021.2" {
	// 	return denyResponse("Invalid agent"), nil
	// }

	// We already can restrict incoming IP using IAM policy, but you can do more fine-grained checks here
	// if request.RequestContext.Identity.SourceIP != "54.54.54.54" {
	// 	return denyResponse("Invalid source"), nil
	// }

	// Can use regexp to validate the format of the machine id, depending on the user's logic
	// _, err := regexp.MatchString(`[A-Z0-9]{8}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{12}`, machineID)
	// if err != nil {
	// 	return denyResponse("Invalid Matching"), nil
	// }
	//
	// if !matched {
	// 	return denyResponse("Incorrect Format"), nil
	// }

	// Or parse the machineID as a uuid
	// _, err := uuid.Parse(machineID)
	// if err != nil {
	// 	return denyResponse("Incorrect Format"), nil
	// }

	return allowResponse(request, machineID), nil
}

func denyResponse(request events.APIGatewayCustomAuthorizerRequestTypeRequest, denyReason string) *events.APIGatewayCustomAuthorizerResponse {
	context := make(map[string]interface{}, 1)
	context["DenyReason"] = denyReason

	return &events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "UNKNOWN_SENSOR",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Deny",
					Resource: []string{"arn:aws:execute-api:*:*:*/*/*/*"},
				},
			},
		},
		Context:            context,
		UsageIdentifierKey: "UNKNOWN_SENSOR",
	}
}

func allowResponse(request events.APIGatewayCustomAuthorizerRequestTypeRequest, principalID string) *events.APIGatewayCustomAuthorizerResponse {
	allowedResourceArns := []string{}
	context := make(map[string]interface{}, 1)
	context["MachineID"] = principalID

	switch request.HTTPMethod {
	case http.MethodGet:
		allowedResourceArns = append(
			allowedResourceArns,
			generateArnResource("GET", "health"),
		)
	case http.MethodPost:
		allowedResourceArns = append(
			allowedResourceArns,
			generateArnResource("POST", "preflight/*"),
			generateArnResource("POST", "eventupload/*"),
			generateArnResource("POST", "ruledownload/*"),
			generateArnResource("POST", "postflight/*"),
			generateArnResource("POST", "xsrf/*"),
		)
	}

	return &events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: allowedResourceArns,
				},
			},
		},
		Context:            context,
		UsageIdentifierKey: principalID,
	}
}

func generateArnResource(verb, resource string) string {
	// Creates a generic resource
	// arn:aws:execute-api:region:account-id:api-id/stage-name/HTTP-VERB/resource-path-specifier
	return fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/%s/%s/%s",
		authorizerEnv.Region,
		authorizerEnv.AccountID,
		authorizerEnv.GatewayID,
		authorizerEnv.StageName,
		verb,
		resource,
	)
}
