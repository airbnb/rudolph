package lambda

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
)

type mockInvokeAPI func(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)

func (m mockInvokeAPI) Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	return m(ctx, params, optFns...)
}

func TestInvokeLambda(t *testing.T) {
	err := invokeLambda(
		context.Background(),
		mockInvokeAPI(func(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
			return &lambda.InvokeOutput{}, nil
		}),
		"lambda-test",
		"lambda-test-abcdef",
		"machineID",
		LambdaEvents{
			Source: "test",
			Items:  []interface{}{"item0", "item1"},
		},
	)

	assert.Nil(t, err)
}
