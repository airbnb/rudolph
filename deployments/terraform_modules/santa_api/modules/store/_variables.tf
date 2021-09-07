variable "prefix" {
  type        = string
  description = "Prefix to all resource names"
}

variable "dynamodb_table_name" {
  type = string
  description = "OPTIONAL: A DynamoDB table name that overrides the default. Used to import a restored backup."
  default = ""
}

variable "read_lambda_role_names" {
  type        = list(string)
  description = "List of IAM Role names that should be allowed to read from rule store tables"
}

variable "aws_account_id" {
  type        = string
  description = "AWS Account Id"
}

variable "region" {
  type        = string
  description = "AWS Region"
}

variable "kms_key_administrators_arns" {
  type = list(string)
  description = "List of KMS Key Administrator ARNs to allow access to Rudolph KMS key operations"
  default = []
}
