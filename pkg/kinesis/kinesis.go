package kinesis

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	awskinesis "github.com/aws/aws-sdk-go/service/kinesis"
)

func (c client) Send(machineID string, events KinesisEvents) (err error) {
	putRecordsRecords := []*awskinesis.PutRecordsRequestEntry{}
	for _, event := range events.Items {
		item, err := json.Marshal(event)
		if err != nil {
			return err
		}

		record := &awskinesis.PutRecordsRequestEntry{
			Data: item,

			// 2021-05-13 - Using the machineID as a partition key
			//   It's not a "bad idea" per se. The alternative of generating a random uuid is not strictly
			//   "better" as the machineID is just a uuid anyway (sometimes. depends on the user).
			//   The implication of this is any data records from a specific machine get put deterministically
			//   onto the same partition. If we have a partition failure, we lose all data from one set of
			//   specific machines, but data for other machines remains intact. It also means that a single
			//   machine with a lot of logs uploaded simultaneously will "burst" onto a single shard.
			PartitionKey: aws.String(machineID),
		}

		putRecordsRecords = append(putRecordsRecords, record)
	}

	putRecordsInput := &awskinesis.PutRecordsInput{
		Records:    putRecordsRecords,
		StreamName: aws.String(c.kinesisStream),
	}

	_, err = c.kinesisService.PutRecords(putRecordsInput)

	return
}
