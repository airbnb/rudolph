package eventupload

import (
	"context"
	"fmt"
	"log"

	"github.com/airbnb/rudolph/pkg/lambda"
)

const RUDOLPH_DIRECT_SOURCE = "rudolph-direct"

func sendToLambda(
	ctx context.Context,
	lambdaClient lambda.LambdaClient,
	machineID string,
	events []EventUploadEvent,
) error {
	var forwardedEvents = convertRequestEventsToUploadEvents(machineID, events)
	err := lambdaClient.Send(
		ctx,
		machineID,
		lambda.LambdaEvents{
			Source: RUDOLPH_DIRECT_SOURCE,
			Items:  forwardedEvents,
		},
	)
	if err != nil {
		log.Printf("Lambda Failed: %s", err)
		return fmt.Errorf("failed to events to AWS Lambda: %w", err)
	}

	return nil
}
