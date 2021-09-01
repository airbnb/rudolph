module "rule_store" {
  source = "./modules/store"

  prefix          = var.prefix
  aws_account_id  = var.aws_account_id
  region          = var.region

  # Add function role names to this that need access to the
  # rule table(s)
  read_lambda_role_names = [
    module.health_function.lambda_role_name,
    module.ruledownload_function.lambda_role_name,
    module.preflight_function.lambda_role_name,
    module.postflight_function.lambda_role_name,
  ]
}
