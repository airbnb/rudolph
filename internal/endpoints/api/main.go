package main

import (
	"github.com/airbnb/rudolph/pkg/handlers"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.ApiRouter)
}
