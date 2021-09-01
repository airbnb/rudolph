package eventupload

import (
	"log"

	"github.com/airbnb/rudolph/pkg/lambda"
	"github.com/pkg/errors"
)

const RUDOLPH_DIRECT_SOURCE = "rudolph-direct"

func sendToLambda(kinesisClient lambda.LambdaClient, machineID string, events []EventUploadEvent) error {
	var forwardedEvents = convertRequestEventsToUploadEvents(machineID, events)
	err := kinesisClient.Send(
		machineID,
		lambda.LambdaEvents{
			Source: RUDOLPH_DIRECT_SOURCE,
			Items:  forwardedEvents,
		},
	)
	if err != nil {
		log.Printf("Lambda Failed: %s", err)
		return errors.Wrap(err, "failed to events to AWS Lambda")
	}

	return nil
}
