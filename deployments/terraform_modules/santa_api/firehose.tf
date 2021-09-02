# This is a convenience IAM policy allowing Rudolph to create the necessary eventupload Firehose resources.
# Enable this using the var.eventupload_autocreate_policies variable.
# Alternatively, specify "false" and create the policies elsewhere.
locals {
  create_firehose = var.eventupload_handler == "FIREHOSE" && var.eventupload_autocreate_policies
}

module "firehose" {
  count = local.create_firehose ? 1 : 0
  source = "./modules/firehose"

  prefix         = var.prefix
  org            = var.org
  aws_account_id = var.aws_account_id
  region         = var.region

  existing_logging_bucket_name = var.existing_logging_bucket_name
  enable_logging               = var.enable_s3_logging

  firehose_upload_lambda_role_names = [
    module.eventupload_function.lambda_role_name,
  ]
}
