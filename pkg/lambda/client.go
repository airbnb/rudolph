package lambda

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// LambdaEvents is a single event entry that appends the MachineID with the EventUploadEvent details
type LambdaEvents struct {
	Source string        `json:"source"`
	Items  []interface{} `json:"items"`
}

type LambdaClient interface {
	InvokeAPI
}

func GetClient(
	functionName string,
	functionQualifier string,
	awsregion string,
) LambdaClient {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(awsregion),
	)
	if err != nil {
		log.Fatalf("unable to load AWS/Lambda SDK config, %v", err)
	}

	svc := lambda.NewFromConfig(cfg)

	return &client{
		lambdaClient:      svc,
		functionName:      functionName,
		functionQualifier: functionQualifier,
	}
}

type client struct {
	lambdaClient      *lambda.Client
	functionName      string
	functionQualifier string
}
