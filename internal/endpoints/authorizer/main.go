package main

import (
	"github.com/airbnb/rudolph/pkg/handlers/authorizer"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(authorizer.HandleAuthorizerRequest)
}
