package lambda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type InvokeAPI interface {
	Send(ctx context.Context, machineID string, events LambdaEvents) error
}

type invokeAPI interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

func (c *client) Send(
	ctx context.Context,
	machineID string,
	events LambdaEvents,
) error {
	return invokeLambda(
		ctx,
		c.lambdaClient,
		c.functionName,
		c.functionQualifier,
		machineID,
		events,
	)
}

func invokeLambda(
	ctx context.Context,
	api invokeAPI,
	functionName string,
	functionQualifier string,
	machineID string,
	events LambdaEvents,
) error {
	item, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("failed json marshall lambda payload events: %w", err)
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		Qualifier:      aws.String(functionQualifier),
		InvocationType: awsTypes.InvocationTypeEvent,
		LogType:        awsTypes.LogTypeNone,
		Payload:        item,
	}

	_, err = api.Invoke(
		ctx,
		invokeInput,
	)
	if err != nil {
		return fmt.Errorf("lambda:InvokeFunction call failed: %w", err)
	}

	return nil
}
