#
# Example tf outputs file
#
#   This file is not strictly necessary in each environment directory, but you can
#   include it to print out your API's hostname at the end of each terraform apply
#   for convenience.
#
output "domain_url" {
  value = module.santa_api.api_domain_name
}

output "raw_url" {
  value = module.santa_api.raw_url
}
