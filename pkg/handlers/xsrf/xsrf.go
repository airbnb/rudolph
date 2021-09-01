package xsrf

import (
	"log"
	"net/http"

	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// XSRFHandler handles requests to the /xsrf and /xsrf/{machine_id} API endpoints
type PostXSRFHandler struct {
}

func (h *PostXSRFHandler) Boot() (err error) {
	return
}

func (h *PostXSRFHandler) Handles(request events.APIGatewayProxyRequest) bool {
	return request.Resource == "/xsrf/{machine_id}" && request.HTTPMethod == "POST"
}

func (h *PostXSRFHandler) Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("XSRFHandler request:\n%+v\n", request)

	// FIXME (derek.wang) just returning stub for now
	return response.APIResponse(http.StatusOK, map[string]string{"status": "ok"})
}
