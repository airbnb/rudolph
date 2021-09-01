module "santa_api" {
  source = "../../terraform_modules/santa_api"

  prefix                             = var.prefix
  region                             = var.region
  aws_account_id                     = var.aws_account_id
  stage_name                         = var.stage_name
  allowed_inbound_ips                = var.allowed_inbound_ips

  # StreamAlert integration
  eventupload_handler             = var.eventupload_handler
  eventupload_kinesis_name        = var.eventupload_kinesis_name
  eventupload_autocreate_policies = var.eventupload_autocreate_policies
  eventupload_output_lambda_name  = var.eventupload_output_lambda_name

  # The route53 zone id
  route53_zone_name   = var.route53_zone_name
  use_existing_route53_zone = var.use_existing_route53_zone

  lambda_zip      = var.zip_file_path
  package_version = var.package_version
}
