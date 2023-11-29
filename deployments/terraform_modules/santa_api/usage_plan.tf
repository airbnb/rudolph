# resource "aws_api_gateway_usage_plan" "example" {
#   name         = "${var.prefix}_rudolph_api_usage_plan"
#   description  = "Usage plan for Rudolph API"

#   api_stages {
#     api_id = aws_api_gateway_rest_api.api_gateway.id
#     stage  = var.stage_name
#   }

#   quota_settings {
#     limit  = 24 * 6 * 5
#     offset = 0
#     period = "DAY"
#   }

#   throttle_settings {
#     burst_limit = 5
#     rate_limit  = 1
#   }

# resource "aws_api_gateway_usage_plan_key" "main" {
#   key_id        = aws_api_gateway_api_key.mykey.id
#   key_type      = "API_KEY"
#   usage_plan_id = aws_api_gateway_usage_plan.myusageplan.id
# }