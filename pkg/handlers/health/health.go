package health

import (
	"net/http"
	"os"

	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/model/globalrules"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

type GetHealthHandler struct {
	booted         bool
	dynamodbClient dynamodb.DynamoDBClient
}

func (h *GetHealthHandler) Boot() (err error) {
	if h.booted {
		return
	}

	dynamodbTableName := os.Getenv("DYNAMODB_NAME")
	awsRegion := os.Getenv("REGION")

	h.dynamodbClient = dynamodb.GetClient(dynamodbTableName, awsRegion)

	h.booted = true
	return
}

func (h *GetHealthHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/health" && request.HTTPMethod == "GET"
}

func (h *GetHealthHandler) Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return h.healthHandler(request)
}

// HealthHandler to be used by the health function
func (h GetHealthHandler) healthHandler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Test dynamodb connection by querying for 1 global rule at random
	// we don't care about the response, as long as DDB doesn't fail with a permission error or something,
	// our system can be considered "healthy"
	err := globalrules.PingDatabase(h.dynamodbClient)
	if err != nil {
		return response.APIResponse(http.StatusInternalServerError, map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	// Otherwise we gucci, return 200
	return response.APIResponse(http.StatusOK, map[string]string{"status": "healthy"})
}
