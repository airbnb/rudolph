package handlers

import (
	"log"
	"net/http"

	"github.com/airbnb/rudolph/internal/handlers/eventupload"
	"github.com/airbnb/rudolph/internal/handlers/health"
	"github.com/airbnb/rudolph/internal/handlers/postflight"
	"github.com/airbnb/rudolph/internal/handlers/preflight"
	"github.com/airbnb/rudolph/internal/handlers/ruledownload"
	"github.com/airbnb/rudolph/internal/handlers/xsrf"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

var (
	handlers []HandlerInterface
)

func init() {
	handlers = []HandlerInterface{
		&health.GetHealthHandler{},
		&eventupload.PostEventuploadHandler{},
		&preflight.PostPreflightHandler{},
		&ruledownload.PostRuledownloadHandler{},
		&postflight.PostPostflightHandler{},
		&xsrf.PostXSRFHandler{},
	}
}

func ApiRouter(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Api Request: %+v", request)

	response, err := getResponse(request)

	if err != nil {
		log.Printf("Api ERROR: %+v", err)
	}
	log.Printf("Api Response: %+v", response)

	return response, err
}

func getResponse(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	for _, h := range handlers {
		if h.Handles(request) {
			if h.Boot() != nil {
				return response.APIResponse(http.StatusInternalServerError, nil)
			}
			return h.Handle(request)
		}
	}

	log.Printf("ERROR: ApiRouter failure: unrouteable request: [%+v] %+v", request.HTTPMethod, request.Resource)
	return response.APIResponse(http.StatusMethodNotAllowed, nil)
}

type HandlerInterface interface {
	Handles(request events.APIGatewayProxyRequest) bool
	Boot() error
	Handle(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}
