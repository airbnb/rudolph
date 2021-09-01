#
# Set up the authorizer
#
module "rudolph_api_authorizer" {
  source = "./modules/lambda/authorizer"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  api_gateway_id            = aws_api_gateway_rest_api.api_gateway.id
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash

  env_vars = {
    REGION      = var.region
    GATEWAY_ID  = aws_api_gateway_rest_api.api_gateway.id
    ACCOUNT_ID  = var.aws_account_id
  }
}
