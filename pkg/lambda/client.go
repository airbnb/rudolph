package lambda

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
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
		err = errors.Wrap(err, "failed json marshall lambda payload events")

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
		err = errors.Wrap(err, "lambda:InvokeFunction call failed")
	}

	// inputs := []*awslambda.InvokeInput{}
	// for _, event := range events.Items {
	// 	item, err := json.Marshal(event)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// https://docs.aws.amazon.com/lambda/latest/dg/API_InvokeAsync.html#API_InvokeAsync_RequestSyntax
	// 	input := &awslambda.InvokeInput{
	// 		Payload: item,
	// 		// FunctionName:   aws.String("arn:aws:lambda:us-east-1:009715504418:function:ryxias_comet_streamalert"),
	// 		FunctionName:   aws.String("ryxias_comet_streamalert"),
	// 		Qualifier:      aws.String("$LATEST"),
	// 		InvocationType: aws.String(awslambda.InvocationTypeEvent),
	// 		LogType:        aws.String(awslambda.LogTypeNone),
	// 	}

	// 	inputs = append(inputs, input)

	// 	_, err := c.lambdaService.Invoke(input)
	// }

	// _, err = c.kinesisService.PutRecords(putRecordsInput)

	return
}
