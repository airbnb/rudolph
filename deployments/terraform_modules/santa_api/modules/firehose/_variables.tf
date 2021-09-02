variable "prefix" {
  type        = string
  description = "Prefix to all resource names"
}

variable "aws_account_id" {
  type        = string
  description = "AWS Account Id"
}

variable "region" {
  type        = string
  description = "AWS Region"
}

variable "enable_logging" {
  type        = bool
  description = "Whether or not S3 access logging should be enabled for the source bucket"
}

variable "existing_logging_bucket_name" {
  type        = string
  description = "Name of existing S3 bucket to use for storing access logs. The default of an empty string will result in a bucket being created"
  default     = ""
}

variable "eventupload_firehose_name" {
  type        = string
  description = "When eventupload_handler is FIREHOSE, specify the Firehose's name"
  default     = ""
}

variable "firehose_upload_lambda_role_names" {
  type        = list(string)
  description = "List of Lambda role names that should be allowed to write and upload to Firehose"
}

variable "org" {
  type        = string
  description = "Organization name"
}
