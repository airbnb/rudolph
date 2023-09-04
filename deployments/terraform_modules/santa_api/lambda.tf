locals {
  lambda_source_hash   = filebase64sha256(var.lambda_api_zip)
  lambda_source_key    = "rudolph-source-${filemd5(var.lambda_api_zip)}.zip"
  lambda_authorizer_hash = filebase64sha256(var.lambda_authorizer_zip)
  lambda_authorizer_source_key    = "rudolph-source-authorizer-${filemd5(var.lambda_authorizer_zip)}.zip"
  lambda_source_bucket = length(module.lambda_source) > 0 ? module.lambda_source[0].bucket_name : var.lambda_source_s3_bucket_name
  dynamodb_table_name = format("%s_rudolph_store", var.prefix)
  firehose_name     = var.eventupload_firehose_name == "" ? format("%s_rudolph_eventsupload_firehose", var.prefix) : var.eventupload_firehose_name
}

#
# Source bucket for Lambda code
#
module "lambda_source" {
  count = var.lambda_source_s3_bucket_name != "" ? 0 : 1

  source                       = "./modules/lambda/lambda-source"
  prefix                       = var.prefix
  org                          = var.org
  existing_logging_bucket_name = var.existing_logging_bucket_name
  enable_logging               = var.enable_s3_logging
}

resource "aws_s3_bucket_object" "santa_api_source" {
  bucket = local.lambda_source_bucket
  key    = local.lambda_source_key
  source = var.lambda_api_zip
  etag   = filemd5(var.lambda_api_zip)
}

resource "aws_s3_bucket_object" "santa_api_authorizer_source" {
  bucket = local.lambda_source_bucket
  key    = local.lambda_authorizer_source_key
  source = var.lambda_authorizer_zip
  etag   = filemd5(var.lambda_authorizer_zip)
}


module "health_function" {
  source = "./modules/lambda/api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash
  endpoint                  = "health"
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn

  env_vars = {
    REGION = var.region
    DYNAMODB_NAME = local.dynamodb_table_name
  }
}


module "xsrf_function" {
  source = "./modules/lambda/api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash
  endpoint                  = "xsrf"
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn

  env_vars = {}
}


module "preflight_function" {
  source = "./modules/lambda/api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash
  endpoint                  = "preflight"
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn

  env_vars = {
    REGION = var.region
    DYNAMODB_NAME = local.dynamodb_table_name
  }
}


module "eventupload_function" {
  source = "./modules/lambda/api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash
  endpoint                  = "eventupload"
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn

  env_vars = {
    REGION        = var.region
    HANDLER       = var.eventupload_handler
    FIREHOSE_NAME = local.firehose_name
    KINESIS_NAME  = var.eventupload_kinesis_name
    LAMBDA_NAME   = var.eventupload_output_lambda_name
  }
}


module "ruledownload_function" {
  source = "./modules/lambda/api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash
  endpoint                  = "ruledownload"
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn

  env_vars = {
    REGION = var.region
    DYNAMODB_NAME = local.dynamodb_table_name
  }
}


module "postflight_function" {
  source = "./modules/lambda/api-handler"

  prefix                    = var.prefix
  region                    = var.region
  alias_name                = var.stage_name
  lambda_source_bucket      = aws_s3_bucket_object.santa_api_source.bucket
  lambda_source_key         = aws_s3_bucket_object.santa_api_source.key
  lambda_source_hash        = local.lambda_source_hash
  endpoint                  = "postflight"
  api_gateway_execution_arn = aws_api_gateway_rest_api.api_gateway.execution_arn

  env_vars = {
    REGION = var.region
    DYNAMODB_NAME = local.dynamodb_table_name
  }
}
