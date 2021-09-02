package eventupload

import (
	"log"

	"github.com/airbnb/rudolph/pkg/firehose"
	"github.com/pkg/errors"
)

func sendToFirehose(firehoseClient firehose.FirehoseClient, machineID string, events []EventUploadEvent) error {
	var forwardedEvents = convertRequestEventsToUploadEvents(machineID, events)
	err := firehoseClient.Send(machineID, firehose.FirehoseEvents{Items: forwardedEvents})
	if err != nil {
		log.Printf("%s\n%s", err.Error(), "upload to firehose was not successful")
		return errors.Wrap(err, "failed to events to AWS Firehose")
	}

	return nil
}
