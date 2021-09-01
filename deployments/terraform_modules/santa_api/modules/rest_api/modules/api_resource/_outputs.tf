output "integration_id" {
  value = aws_api_gateway_integration.resource_integration.id
}

output "integration_sha" {
  value = sha1(jsonencode(aws_api_gateway_integration.resource_integration))
}
