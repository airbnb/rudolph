
# Gateway Method
# Defines which HTTP method this particular resource action listens on
# It also defines how the request is authorized.
resource "aws_api_gateway_method" "resource_method" {
  rest_api_id   = var.gateway_rest_api_id
  resource_id   = var.gateway_resource_id
  http_method   = var.http_method
  authorization = var.authorizer_id == "" ? "NONE" : "CUSTOM"
  authorizer_id = var.authorizer_id

  request_models = var.request_model == "" ? {} : {
    "application/json" = var.request_model
  }
}

# Gateway Integration
# The Gateway integration specifies whether to proxy the request to AWS Lambda, or just MOCKs a response
resource "aws_api_gateway_integration" "resource_integration" {
  rest_api_id             = var.gateway_rest_api_id
  resource_id             = var.gateway_resource_id
  http_method             = aws_api_gateway_method.resource_method.http_method
  integration_http_method = "POST" # Lambda can only be invoked with POST
  type                    = "AWS_PROXY"

  # This isn't the function arn or qualified function arn. It's an invocation ARN which is actually scoped to API Gate
  # arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:{######}:function:{function_name}/invocations
  uri = var.lambda_invocation_arn
}

# Integration Response
# We don't use this for AWS_PROXY type integrations, since there is no reason to transform the result; the Lambda should
# setup the result exactly.
#

# Method Response 200
resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = var.gateway_rest_api_id
  resource_id = var.gateway_resource_id
  http_method = aws_api_gateway_method.resource_method.http_method
  status_code = "200"

  # Alternative, we use this re-usable module?
  # https://github.com/squidfunk/terraform-aws-api-gateway-enable-cors
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true # FIXME (derek.wang) only send this on cors-enabled endpoints?
    "method.response.header.Content-Type"                = true
  }

  response_models = var.success_response_model == "" ? {} : {
    "application/json" = var.success_response_model
  }
}
