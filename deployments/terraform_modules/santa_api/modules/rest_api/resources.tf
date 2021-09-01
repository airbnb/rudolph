#
# This is the main terraform file of the api_resource module. This module allows you to create a single
# REST API resource, and associate it with zero or more Resource Methods.
#
# Resource Methods allow you to invoke the resource via HTTP methods. Each resource method is composed of
# 4 primary components:
# - Gateway Method (aka Method Request)
# - Gateway Integration (aka Integration Request)
# - Integration Response
# - Method Response
#

# The resource
resource "aws_api_gateway_resource" "resource" {
  rest_api_id = var.gateway_rest_api_id
  parent_id   = var.parent_resource_id
  path_part   = var.resource_path
}


module "api_methods" {
  for_each = toset(var.integration_http_methods)

  source = "./modules/api_resource"

  gateway_resource_id   = aws_api_gateway_resource.resource.id
  gateway_rest_api_id   = var.gateway_rest_api_id
  http_method           = each.value
  lambda_invocation_arn = var.lambda_invocation_arn
  authorizer_id         = var.authorizer_id

  success_response_model = var.success_response_model
  request_model          = var.request_model
}
