#
# Top level REST API resource
#
resource "aws_api_gateway_rest_api" "api_gateway" {
  name        = "${var.prefix}_rudolph"
  description = "Rudolph's REST API Gateway"

  # If allowed_inbound_ips are set, then use the API IAM resource policy else leave this to a null resource
  policy = length(var.allowed_inbound_ips) != 0 ? data.aws_iam_policy_document.api_resource_policy.json : null

  endpoint_configuration {
    types = ["REGIONAL"] # FIXME (derek.wang) Switch to PRIVATE later and attach some VPCs
  }

  # Use the authorizer's UsageIdentifierKey to uniquely identify an endpoint.
  api_key_source = "AUTHORIZER"
}

##########################
# Stages and Deployments #
##########################

# resource "aws_api_gateway_stage" "stage" {
#   stage_name    = "v1"
#   rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
#   deployment_id = aws_api_gateway_deployment.api_deployment.id
# }


resource "aws_api_gateway_deployment" "api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  stage_name  = var.stage_name

  triggers = {
    redeployment = join(",", [
      sha1(jsonencode(aws_api_gateway_rest_api.api_gateway)),
      module.health_api.integration_shas,
      module.ruledownload_api.integration_shas,
      module.ruledownload_resource_api.integration_shas,
      module.preflight_api.integration_shas,
      module.preflight_resource_api.integration_shas,
      module.eventupload_api.integration_shas,
      module.eventupload_resource_api.integration_shas,
      module.xsrf_api.integration_shas,
      module.xsrf_resource_api.integration_shas,
      module.postflight_api.integration_shas,
      module.postflight_resource_api.integration_shas,
    ])
  }

  depends_on = [
    module.health_api.integration_ids,
    module.ruledownload_api.integration_ids,
    module.ruledownload_resource_api.integration_ids,
    module.preflight_api.integration_ids,
    module.preflight_resource_api.integration_ids,
    module.eventupload_api.integration_ids,
    module.eventupload_resource_api.integration_ids,
    module.xsrf_api.integration_ids,
    module.xsrf_resource_api.integration_ids,
    module.postflight_api.integration_ids,
    module.postflight_resource_api.integration_ids,
  ]

  lifecycle {
    create_before_destroy = true
  }
}

# resource "aws_api_gateway_method_settings" "s" {
#   rest_api_id = aws_api_gateway_rest_api.api_gateway.id
#   stage_name  = var.stage_name
#   method_path = "*/*"

#   settings {
#     metrics_enabled = true
#     logging_level   = "INFO"

#   #   # In most cases, each client makes:
#   #   # - 1x preflight
#   #   # - 0 or more eventuploads
#   #   # - 2x ruledownload
#   #   # - 1x preflights
#   #   #
#   #   # Roughly 10 requests every checkin, in a worst case. A good burst rate is equal to the number of
#   #   # (anticipated) clients, and the rate limit should just be:
#   #   #
#   #   #  - number of clients * (10 req/checkin) * (1/10 checkin/min) * (1/60 min/sec) = clients/60 req/sec
#   #   throttling_burst_limit = 7000
#   #   throttling_rate_limit = 7000 / 60
#   }
# }

# @doc https://www.terraform.io/docs/providers/aws/r/api_gateway_base_path_mapping.html
resource "aws_api_gateway_base_path_mapping" "gateway_custom_domain_mapping" {
  api_id      = aws_api_gateway_rest_api.api_gateway.id
  stage_name  = aws_api_gateway_deployment.api_deployment.stage_name
  domain_name = aws_api_gateway_domain_name.api_custom_domain.domain_name
}

