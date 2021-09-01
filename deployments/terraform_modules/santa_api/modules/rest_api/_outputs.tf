output "resource_id" {
  value = aws_api_gateway_resource.resource.id
}

output "integration_ids" {
  value = toset([for item in module.api_methods : item.integration_id])
}

output "integration_shas" {
  value = join(",", sort([for item in module.api_methods : item.integration_sha]))
}
