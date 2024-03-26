package eventupload

import (
	"fmt"
	"log"

	"github.com/airbnb/rudolph/pkg/kinesis"
)

func sendToKinesis(kinesisClient kinesis.KinesisClient, machineID string, events []EventUploadEvent) error {
	var forwardedEvents = convertRequestEventsToUploadEvents(machineID, events)
	err := kinesisClient.Send(machineID, kinesis.KinesisEvents{Items: forwardedEvents})
	if err != nil {
		log.Printf("Kinesis Failed: %s", err)
		return fmt.Errorf("failed to events to AWS kinesis: %w", err)
	}

	return nil
}
