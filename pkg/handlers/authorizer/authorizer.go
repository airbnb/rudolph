package authorizer

import (
	"errors"
	"fmt"
	"log"
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
}

var (
	authorizerEnv authorizerEnvironment
)

func init() {
	authorizerEnv = authorizerEnvironment{
		Region:    os.Getenv("REGION"),
		GatewayID: os.Getenv("GATEWAY_ID"),
		AccountID: os.Getenv("ACCOUNT_ID"),
	}
}

// HandleAuthorizerRequest is the handler to be used by the authorizer function
func HandleAuthorizerRequest(request events.APIGatewayProxyRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
	log.Printf("lambda request - HandleAuthorizerRequest:\n%+v\n", request)

	if request.HTTPMethod == "GET" && request.Path == "/health" {
		return allowResponse("HEALTH_CHECK"), nil
	}

	if request.HTTPMethod != "POST" {
		return denyResponse("Incorrect Method"), nil
	}

	machineID, ok := request.PathParameters["machine_id"]
	if !ok {
		return denyResponse("Incorrect Request URI"), nil
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

	return allowResponse(machineID), nil
}

func denyResponse(denyReason string) *events.APIGatewayCustomAuthorizerResponse {
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

func allowResponse(machineID string) *events.APIGatewayCustomAuthorizerResponse {
	context := make(map[string]interface{}, 1)
	context["MachineID"] = machineID

	// Creates a generic resource
	//"arn:aws:execute-api:*:*:*/*/*/*"
	//<RANDOM_KEY>/<STAGE>/<HTTP_METHOD>/<URLPATH>
	resourceArn := fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*/*", authorizerEnv.Region, authorizerEnv.AccountID, authorizerEnv.GatewayID)
	return &events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "ValidSantaEndpoint",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{resourceArn},
				},
			},
		},
		Context:            context,
		UsageIdentifierKey: machineID,
	}
}
