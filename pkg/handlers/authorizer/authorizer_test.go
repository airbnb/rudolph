package authorizer

// import (
// 	"testing"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGenericHandler(t *testing.T) {
// 	tests := []struct {
// 		request events.APIGatewayProxyRequest
// 		expect  string
// 		err     error
// 	}{
// 		{
// 			// // Test to see if genericHandler responds the provided value
// 			request: events.APIGatewayProxyRequest{Body: "test"},
// 			expect:  "Response Body: test\nFuncName: ; FuncVers: \n",
// 			err:     nil,
// 		},
// 		{
// 			// Test to see if genericHandler responds with a hello message
// 			request: events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{
// 				"name": "ryan",
// 			}},
// 			expect: "Hello, ryan!\nFuncName: ; FuncVers: \n",
// 			err:    nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		response, err := genericHandler(test.request)
// 		assert.IsType(t, test.err, err)
// 		assert.Equal(t, test.expect, response.Body)
// 	}
// }
