package eventupload

import (
	"log"

	"github.com/airbnb/rudolph/pkg/kinesis"
	"github.com/pkg/errors"
)

func sendToKinesis(kinesisClient kinesis.KinesisClient, machineID string, events []EventUploadEvent) error {
	var forwardedEvents = convertRequestEventsToUploadEvents(machineID, events)
	err := kinesisClient.Send(machineID, kinesis.KinesisEvents{Items: forwardedEvents})
	if err != nil {
		log.Printf("Kinesis Failed: %s", err)
		return errors.Wrap(err, "failed to events to AWS kinesis")
	}

	return nil
}
