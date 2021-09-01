# Custom Gateway responses for different error types.
#
# Here I use Gateway Responses to set custom responses whenever an error is encountered.
#
# For a list of all valid types, go here:
# https://docs.aws.amazon.com/apigateway/api-reference/resource/gateway-response/
# and they are explained here:
# https://docs.aws.amazon.com/apigateway/latest/developerguide/supported-gateway-response-types.html
#
# Remember that the exact error will always show up in the "x-amzn-ErrorType" response header.

# UNAUTHORIZED occurs when the authorizer fails to authenticate; usually when the Token is missing
resource "aws_api_gateway_gateway_response" "unauthorized" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  response_type = "UNAUTHORIZED"
  status_code   = "401"

  response_templates = {
    "application/json" = "{\"message\":$context.error.messageString}"
  }
}

# AUTHORIZER_CONFIGURATION_ERROR is when the Authorizer raises an exception
resource "aws_api_gateway_gateway_response" "authorization_configuration_error" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  response_type = "AUTHORIZER_CONFIGURATION_ERROR"
  status_code   = "500"

  response_templates = {
    "application/json" = "{\"message\":$context.error.messageString}"
  }
}

# ACCESS_DENIED is encountered comes when the Authorizer succeeds and explicitly denies access
#   Note: Because this error only occurs after authorizer is reached, we can bubble up authorization errors
#   that come from the authorizer using this $context.authorizer.error variable.
#
#   For more information on these: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
resource "aws_api_gateway_gateway_response" "access_denied" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  response_type = "ACCESS_DENIED"
  status_code   = "403"

  response_templates = {
    "application/json" = "{\"message\":$context.error.messageString,\"error\":\"$context.authorizer.error\"}"
  }
}

# DEFAULT_5XX is just for any generic internal server error generated from the Lambda. According to the documentation
#   This overrides all other 500 type errors, so this may override the "AUTHORIZER_CONFIGURATION_ERROR" although I may
#   be misunderstanding...
resource "aws_api_gateway_gateway_response" "default_5xx" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  response_type = "DEFAULT_5XX"
  status_code   = "500"

  response_templates = {
    "application/json" = "{\"message\":$context.error.messageString}"
  }
}


# MISSING_AUTHENTICATION_TOKEN is The gateway response for a missing authentication token error, including the cases when
# the client attempts to invoke an unsupported API method or resource. If the response type is unspecified, this response
# defaults to the DEFAULT_4XX type.
resource "aws_api_gateway_gateway_response" "missing_authentication_token" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  response_type = "MISSING_AUTHENTICATION_TOKEN"
  status_code   = "400"

  response_templates = {
    "application/json" = "{\"message\":\"Something went wrong\"}"
  }
}