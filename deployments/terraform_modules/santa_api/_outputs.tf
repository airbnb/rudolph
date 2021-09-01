output "api_gateway_id" {
  value = aws_api_gateway_rest_api.api_gateway.id
}

output "api_gateway_execution_arn" {
  value = aws_api_gateway_rest_api.api_gateway.execution_arn
}

output "api_domain_name" {
  value = local.api_domain_name
}

output "raw_url" {
  value = aws_api_gateway_deployment.api_deployment.invoke_url
}
