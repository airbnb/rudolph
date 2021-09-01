package sensordata

import (
	"github.com/airbnb/rudolph/pkg/dynamodb"
	"github.com/airbnb/rudolph/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	MachineID_DataType_GSI           string = "DataType_MachineID"
	SerialNum_DataType_MachineID_GSI string = "SerialNum_DataType_MachineID"
)

func GetSensorDataFinder(api dynamodb.QueryAPI) SensorDataFinder {
	return ConcreteSensorDataFinder{
		queryapi: api,
	}
}

//
// SensorDataFinder is a service that allows you to search the SensorData logs for machineIDs or serial numbers of sensors that have recently
// checked into Rudolph
//
// This service is particularly useful for typeahead and can be re-used for external applications.
//
type SensorDataFinder interface {
	GetMachineIDsStartingWith(prefix string, limit int32) ([]string, error)
	GetMachineIDsFromSerialNumber(serialNumber string, limit int32) ([]string, error)
}

type ConcreteSensorDataFinder struct {
	queryapi dynamodb.QueryAPI
}

func (f ConcreteSensorDataFinder) GetMachineIDsStartingWith(prefix string, limit int32) (machineIDs []string, err error) {
	var keyCond expression.KeyConditionBuilder
	// Build the key conditions
	if prefix != "" {
		keyCond = expression.KeyAnd(
			expression.Key("DataType").Equal(expression.Value(string(GetDataType()))),
			expression.Key("MachineID").BeginsWith(prefix),
		)
	} else {
		keyCond = expression.Key("DataType").Equal(expression.Value(string(GetDataType())))
	}

	proj := expression.NamesList(expression.Name("MachineID"))

	builder := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj)
	expr, err := builder.Build()
	if err != nil {
		return
	}

	input := awsdynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		IndexName:                 aws.String(MachineID_DataType_GSI),
		Limit:                     aws.Int32(limit),
		ConsistentRead:            aws.Bool(false),
	}

	output, err := f.queryapi.Query(&input)
	if err != nil {
		return
	}

	machineIDs = make([]string, len(output.Items))
	gsiItems := make([]dataTypeMachineIDGSIItem, len(output.Items))
	err = attributevalue.UnmarshalListOfMaps(output.Items, &gsiItems)
	if err != nil {
		return
	}

	for index, item := range gsiItems {
		machineIDs[index] = item.MachineID
	}

	return
}

func (f ConcreteSensorDataFinder) GetMachineIDsFromSerialNumber(serialNumber string, limit int32) (machineIDs []string, err error) {
	// Build the key conditions
	keyCond := expression.KeyAnd(
		expression.Key("SerialNum").Equal(expression.Value(serialNumber)),
		expression.Key("DataType").Equal(expression.Value(string(GetDataType()))),
	)

	proj := expression.NamesList(expression.Name("MachineID"))

	builder := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj)
	expr, err := builder.Build()
	if err != nil {
		return
	}

	input := awsdynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		IndexName:                 aws.String(SerialNum_DataType_MachineID_GSI),
		Limit:                     aws.Int32(limit),
		ConsistentRead:            aws.Bool(false),
	}

	output, err := f.queryapi.Query(&input)
	if err != nil {
		return
	}

	machineIDs = make([]string, len(output.Items))
	gsiItems := make([]SerialNumDataTypeMachineIdGSIItem, len(output.Items))
	err = attributevalue.UnmarshalListOfMaps(output.Items, &gsiItems)
	if err != nil {
		return
	}

	for index, item := range gsiItems {
		machineIDs[index] = item.MachineID
	}

	return
}

type SerialNumDataTypeMachineIdGSIItem struct {
	PrimaryKey dynamodb.PrimaryKey
	SerialNumb string         `dynamodbav:"SerialNum"`
	MachineID  string         `dynamodbav:"MachineID"`
	DataType   types.DataType `dynamodbav:"DataType"`
}

type dataTypeMachineIDGSIItem struct {
	PrimaryKey dynamodb.PrimaryKey
	DataType   types.DataType `dynamodbav:"DataType"`
	MachineID  string         `dynamodbav:"MachineID"`
}
