package scan

import (
	"github.com/airbnb/rudolph/pkg/dynamodb"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetScanService(scanner dynamodb.ScanAPI) ScanService {
	return ConcreteScanService{
		scanner: scanner,
	}
}

// ScanService provides the functionality to iterate over all entries in DynamoDB given a scan query
// Each output page is passed to the given scan output function. Afterwards, it is passed to the stop
// function, which can return True to immediately halt the scan operation. Otherwise, the scan operation
// will continue through all pages in DynamoDB until it runs out of records to scan.
//
// Returns error or nil. If you need to extract specific records from it, you should put that logic into
// the callback function. If you need to halt the entire operation and resume from a previous state,
// save the lastEvaluatedKey in the callback and pass it to the scan input argument.
type ScanService interface {
	ScanAll(in awsdynamodb.ScanInput, callback func(out *awsdynamodb.ScanOutput) error, stop func(out *awsdynamodb.ScanOutput) (bool, error)) error
}

type ConcreteScanService struct {
	scanner dynamodb.ScanAPI
}

func (c ConcreteScanService) ScanAll(in awsdynamodb.ScanInput, callback func(out *awsdynamodb.ScanOutput) error, stop func(out *awsdynamodb.ScanOutput) (bool, error)) (err error) {
	nextInput := &in
	for {
		shouldStop, lastEvaluatedKey, err := scanIterator(c.scanner, nextInput, callback, stop)
		if err != nil {
			break
		}

		if shouldStop {
			break
		}

		if lastEvaluatedKey != nil {
			nextInput.ExclusiveStartKey = lastEvaluatedKey
			continue
		}

		break
	}

	return
}

func scanIterator(scanner dynamodb.ScanAPI, in *awsdynamodb.ScanInput, callback func(out *awsdynamodb.ScanOutput) error, stop func(out *awsdynamodb.ScanOutput) (bool, error)) (bool, map[string]types.AttributeValue, error) {
	out, err := scanner.Scan(in)
	if err != nil {
		return true, nil, err
	}

	err = callback(out)
	if err != nil {
		return true, nil, err
	}

	if out.LastEvaluatedKey == nil {
		return true, nil, nil
	}

	shouldStop, err := stop(out)
	if err != nil {
		return true, nil, err
	}

	return shouldStop, out.LastEvaluatedKey, nil
}
