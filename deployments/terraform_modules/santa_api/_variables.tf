variable "prefix" {
  type        = string
  description = "Prefix to all resource names"
}

# NOTE:
#   By default, Rudolph creates a base table.
#
#   When changing the dynamodb table for a backup, make sure that:
#   The dynamodb has a hash key of "PK", range key of "SK". If the table schema doesn't match,
#   after you import it Terraform will try to destroy the table and recreate it with the correct attributes,
#   which will destroy your table data. Not desirable.
#
#   When rotating in a new table that has been created from a backup snapshot under a different table
#   name, make sure to use "terraform import" to bring the resource in.
variable "dynamodb_table_name" {
  type = string
  description = "OPTIONAL: A DynamoDB table name that overrides the default. Used to import a restored backup."
  default = ""
}

variable "region" {
  type        = string
  description = "AWS Region"
}

variable "aws_account_id" {
  type        = string
  description = "AWS Account Id"
}

variable "org" {
  type        = string
  description = "Organization name"
}

variable "stage_name" {
  type        = string
  description = "Name of stage to use for this deployment"
}

variable "lambda_api_zip" {
  type        = string
  description = "Full path to zip with go binary for Lambda to be uploaded to S3"
}

variable "lambda_authorizer_zip" {
  type        = string
  description = "Full path to zip with go binary for Lambda to be uploaded to S3"
}

variable "use_existing_route53_zone" {
  type        = bool
  description = "Whether or not to import an existing Route 53 Hosted Zone"
  default     = false
}

variable "route53_zone_name" {
  type        = string
  description = "The name of a Route 53 Hosted Zone to use. Omit this to create a Route 53 Hosted Zone"
  default     = ""
}

variable "lambda_source_s3_bucket_name" {
  type        = string
  description = "Name of S3 bucket to use for storing Lambda source code. If no name is supplied, a bucket will be created,"
  default     = ""
}

variable "enable_s3_logging" {
  type        = bool
  description = "Whether or not S3 access logging should be enabled for the source bucket"
  default     = true
}

variable "existing_logging_bucket_name" {
  type        = string
  description = "Name of S3 bucket to use for S3 access logs"
  default     = ""
}

variable "allowed_inbound_ips" {
  type        = list(string)
  description = "Restricts all inbound access to given ip addresses. By default, allows all inbound access."
  default     = []
}

variable "eventupload_handler" {
  type        = string
  description = "One of: [KINESIS, FIREHOSE]. Specifies where to send uploaded events. By default sends to /dev/null"
  default     = "NONE"

  validation {
    condition     = contains(["NONE", "KINESIS", "FIREHOSE"], var.eventupload_handler)
    error_message = "Valid values for eventupload_handler: (KINESIS, FIREHOSE, NONE)."
  }
}

variable "eventupload_firehose_name" {
  type        = string
  description = "When eventupload_handler is FIREHOSE, specify the Firehose's name"
  default     = ""
}

variable "eventupload_kinesis_name" {
  type        = string
  description = "When eventupload_handler is KINESIS, specify the Kinesis Data Stream's name"
  default     = ""
}

variable "eventupload_autocreate_policies" {
  type        = bool
  description = "When specifying an eventupload handler, will automatically create IAM policies to call the desired resources"
  default     = false
}

variable "eventupload_output_lambda_name" {
  type        = string
  description = "A lambda function name. When provided, the eventupload endpoint will also send events to this lambda, in addition to Kinesis or Firehose."
  default     = ""
}

variable "kms_key_administrators_arns" {
  type = list(string)
  description = "List of KMS Key Administrator ARNs to allow access to Rudolph KMS key operations"
  default = []
}
