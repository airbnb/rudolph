#
# REST API Authorizer
#
# This module covers the REST API Gateway Authorizer
#
resource "aws_api_gateway_authorizer" "api_authorizer" {
  name                   = "${var.prefix}_rudolph_lambda_authorizer" # Authorizers are scoped per gateway, so names don't conflict
  rest_api_id            = var.api_gateway_id
  authorizer_uri         = module.authorizer_function.lambda_invoke_arn
  authorizer_credentials = aws_iam_role.invocation_role.arn

  # Options defining how this authorizer operates
  type                             = "REQUEST" # Validate entire request

  # https://github.com/hashicorp/terraform-provider-aws/issues/5845#issuecomment-517998604
  identity_source                  = ""
  authorizer_result_ttl_in_seconds = 0 # Only GET methods are cached by default so just set to zero to reduce confusion
}

#
# Lambda for authorizer
#
module "authorizer_function" {
  source = "../api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.alias_name
  lambda_source_bucket      = var.lambda_source_bucket
  lambda_source_key         = var.lambda_source_key
  lambda_source_hash        = var.lambda_source_hash
  endpoint                  = "authorizer"
  api_gateway_execution_arn = var.api_gateway_execution_arn

  env_vars = var.env_vars

}
