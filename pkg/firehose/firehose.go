package firehose

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsfirehose "github.com/aws/aws-sdk-go/service/firehose"
)

const (
	maxRetries int = 3
	maxRecords int = 500
)

type FirehoseEvents struct {
	Items []interface{}
}

type firehoseEventBatch struct {
	Items []interface{}
}

// Send will upload EventUploadEvents appended with the machineID to Firehose
// This Send function is only for EventUploadEvents and will automatically split the records into equally spaced partitions as defined
func (c client) Send(machineID string, events FirehoseEvents) (err error) {
	batches := eventBatches(events, machineID, maxRecords)
	for _, batch := range batches {
		var response *awsfirehose.PutRecordBatchOutput
		// Loop if not erroring and maxRetries not hit
		for attempt := 1; attempt < maxRetries; attempt++ {
			if response != nil {
				// Filter out all successful records, leaving only the failed ones
				batch, err = c.filterRecords(response, batch)
				if err != nil {
					break // unrecoverable error of some sort
				}
			}
			response, err = c.sendToFirehose(machineID, batch)
			if err != nil {
				break // do not continue to retry if there is an error
			}
			if int(*response.FailedPutCount) == 0 {
				break // all records sent properly, so exit loop
			}
			attempt++
		}
	}
	return
}

func (c client) filterRecords(response *awsfirehose.PutRecordBatchOutput, batch firehoseEventBatch) (failedBatch firehoseEventBatch, err error) {
	var failedEvents []interface{}
	for i, requestResponse := range response.RequestResponses {
		if requestResponse.RecordId != nil {
			continue
		}

		awsErr := *requestResponse.ErrorMessage
		switch awsErr {
		case awsfirehose.ErrCodeInvalidKMSResourceException:
			err = errors.New(awsErr)
			return
		case awsfirehose.ErrCodeServiceUnavailableException:
			time.Sleep(1 * time.Second)
		case awsfirehose.ErrCodeLimitExceededException:
			err = errors.New(awsErr)
			return
		}

		if requestResponse.ErrorCode != nil || requestResponse.ErrorMessage != nil {
			failedEvents = append(failedEvents, batch.Items[i])
		}
	}

	log.Printf("retrying a total of %d failed records\n", len(failedEvents))
	return firehoseEventBatch{Items: failedEvents}, nil
}

// sendToFirehose appends the machineID and EventUploadEvents together and ships them via a PutRecordBatch request
// if errors occur in the batch request, it will trigger attempt to retry them up to the max limit
func (c client) sendToFirehose(machineID string, events firehoseEventBatch) (response *awsfirehose.PutRecordBatchOutput, err error) {
	recordsBatchInput := &awsfirehose.PutRecordBatchInput{
		DeliveryStreamName: aws.String(c.firehoseStreamName),
	}
	records := []*awsfirehose.Record{}

	for _, event := range events.Items {
		batch, err := json.Marshal(event)

		if err != nil {
			log.Printf("marshalling firehose eventupload event encountered an error, %v", err)
			return nil, err
		}

		records = append(
			records,
			&awsfirehose.Record{Data: append(batch, '\n')},
		)
	}

	recordsBatchInput = recordsBatchInput.SetRecords(records)
	response, err = c.firehoseService.PutRecordBatch(recordsBatchInput)
	if err != nil {
		log.Printf("PutRecordBatch err: %v\n", err)
	}

	return // response and error returned into order to see if retry needed
}

// eventBatches takes a slice of EventUploadEvents and splits them to a defined limit
// the limit should be set to the firehose putbatchrecords max entries per request of 500
// This creates our own array of FirehoseEventUploadEvent entries to avoid doing it multiple times later
func eventBatches(events FirehoseEvents, machineID string, limit int) []firehoseEventBatch {
	var batches []firehoseEventBatch
	slice := events.Items
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < limit {
			limit = len(slice)
		}

		var batchItems []interface{}
		for _, item := range slice[0:limit] {
			batchItems = append(batchItems, item)
		}
		batch := firehoseEventBatch{Items: batchItems}

		batches = append(batches, batch)
		slice = slice[limit:]
	}

	return batches
}
