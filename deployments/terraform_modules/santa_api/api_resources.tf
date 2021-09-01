#
# API Resources
#

# /health resources
module "health_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = aws_api_gateway_rest_api.api_gateway.root_resource_id
  resource_path            = "health"
  integration_http_methods = ["GET"]
  lambda_invocation_arn    = module.health_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}


# /preflight resources
module "preflight_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = aws_api_gateway_rest_api.api_gateway.root_resource_id
  resource_path            = "preflight"
  integration_http_methods = []
  lambda_invocation_arn    = module.preflight_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}
module "preflight_resource_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = module.preflight_api.resource_id
  resource_path            = "{machine_id}"
  integration_http_methods = ["POST"]
  lambda_invocation_arn    = module.preflight_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
  success_response_model   = aws_api_gateway_model.machine_config.name
  request_model            = aws_api_gateway_model.preflight_request.name
}


# /ruledownload resources
module "ruledownload_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = aws_api_gateway_rest_api.api_gateway.root_resource_id
  resource_path            = "ruledownload"
  integration_http_methods = []
  lambda_invocation_arn    = module.ruledownload_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}
module "ruledownload_resource_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = module.ruledownload_api.resource_id
  resource_path            = "{machine_id}"
  integration_http_methods = ["POST"]
  lambda_invocation_arn    = module.ruledownload_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
  success_response_model   = aws_api_gateway_model.ruledownload.name
}

# /eventupload resources
module "eventupload_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = aws_api_gateway_rest_api.api_gateway.root_resource_id
  resource_path            = "eventupload"
  integration_http_methods = []
  lambda_invocation_arn    = module.eventupload_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}
module "eventupload_resource_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = module.eventupload_api.resource_id
  resource_path            = "{machine_id}"
  integration_http_methods = ["POST"]
  lambda_invocation_arn    = module.eventupload_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}

# /xsrf resources
module "xsrf_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = aws_api_gateway_rest_api.api_gateway.root_resource_id
  resource_path            = "xsrf"
  integration_http_methods = []
  lambda_invocation_arn    = module.xsrf_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}
module "xsrf_resource_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = module.xsrf_api.resource_id
  resource_path            = "{machine_id}"
  integration_http_methods = ["POST"]
  lambda_invocation_arn    = module.xsrf_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}


# /postflight resources
module "postflight_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = aws_api_gateway_rest_api.api_gateway.root_resource_id
  resource_path            = "postflight"
  integration_http_methods = []
  lambda_invocation_arn    = module.postflight_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}
module "postflight_resource_api" {
  source = "./modules/rest_api"

  gateway_rest_api_id      = aws_api_gateway_rest_api.api_gateway.id
  parent_resource_id       = module.postflight_api.resource_id
  resource_path            = "{machine_id}"
  integration_http_methods = ["POST"]
  lambda_invocation_arn    = module.postflight_function.lambda_alias_invoke_arn
  authorizer_id            = module.rudolph_api_authorizer.api_gateway_authorizer_id
}
