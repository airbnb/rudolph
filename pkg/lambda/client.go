package lambda

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaEvents is a single event entry that appends the MachineID with the EventUploadEvent details
type LambdaEvents struct {
	Source string        `json:"source"`
	Items  []interface{} `json:"items"`
}

type LambdaClient interface {
	Send(machineID string, events LambdaEvents) (err error)
}

func GetClient(functionName string, functionQualifier string, awsregion string) LambdaClient {
	sess := awssession.Must(awssession.NewSession())

	return client{
		lambdaService:     awslambda.New(sess, aws.NewConfig().WithRegion(awsregion)),
		functionName:      functionName,
		functionQualifier: functionQualifier,
	}
}

type client struct {
	lambdaService     *awslambda.Lambda
	functionName      string
	functionQualifier string
}

func (c client) Send(machineID string, events LambdaEvents) (err error) {
	item, err := json.Marshal(events)
	if err != nil {
		err = fmt.Errorf("failed json marshall lambda payload events: %w", err)

		return
	}

	// https://docs.aws.amazon.com/lambda/latest/dg/API_InvokeAsync.html#API_InvokeAsync_RequestSyntax
	input := &awslambda.InvokeInput{
		FunctionName:   aws.String(c.functionName),
		Qualifier:      aws.String(c.functionQualifier),
		InvocationType: aws.String(awslambda.InvocationTypeEvent),
		LogType:        aws.String(awslambda.LogTypeNone),
		Payload:        item,
	}

	_, err = c.lambdaService.Invoke(input)
	if err != nil {
		err = fmt.Errorf("lambda:InvokeFunction call failed: %w", err)
	}

	return
}
