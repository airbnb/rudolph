variable "zip_file_path" {
  type        = string
  description = "Path to zip on disk to use for deployment of Lambda functions. This gets passed in from 'make deploy' command"
}

variable "package_version" {
  type        = string
  description = "Version of golang binary being used. This value comes from git tags and is used when the Lambda package is uploaded to S3. This gets passed in from 'make deploy' command"
}

// These variables are provided by config.auto.tfvars.json
variable "region" {
  type    = string

  validation {
    condition     = can(regex("[a-z0-9\\-]+", var.region))
    error_message = "The aws region is not valid."
  }
}

variable "aws_account_id" {
  type    = string

  validation {
    condition     = can(regex("[0-9]+", var.aws_account_id))
    error_message = "The aws_account_id should be a number."
  }
}

variable "org" {
  type        = string
  description = "A lowercase string unique to your organization. Used to deduplicate certain AWS resources."

  validation {
    condition     = can(regex("[a-z0-9\\-]+", var.org))
    error_message = "The org value must contain only lowercase alphanumeric characters or dashes."
  }
}

variable "stage_name" {
  type        = string
  description = "Name of stage to use for this deployment"

  validation {
    condition     = can(regex("[a-z\\-]+", var.stage_name))
    error_message = "The stage value must contain only lowercase alphabet characters."
  }
}

variable "prefix" {
  type        = string
  description = "Unique prefix to use for resources. This gets passed in from 'make deploy' command"

  validation {
    condition     = can(regex("[a-z0-9\\-]+", var.prefix))
    error_message = "The prefix value must contain only lowercase alphanumeric characters or dashes."
  }
}

variable "allowed_inbound_ips" {
  type        = list(string)
  description = "Restricts all inbound access to given ip addresses. By default, allows all inbound access."
  default     = []
}

variable "eventupload_handler" {
  type = string
  description = "Type of handler to use for eventupload"
}

variable "eventupload_kinesis_name" {
  type = string
  description = "Name of the kinesis stream"
}

variable "eventupload_autocreate_policies" {
  type = bool
  default = true
}

variable "route53_zone_name" {
  type = string
  description = "Route 53 Zone ID"
}

variable "use_existing_route53_zone" {
  type = bool
  description = "Use exisiting Route 53 zone"
  default = false
}

variable "eventupload_output_lambda_name" {
  type = string
  default = ""
}

variable "enable_s3_logging" {
  type = bool
  default = true
}

variable "kms_key_administrators_arns" {
  type = list(string)
  description = "List of KMS Key Administrator ARNs to allow access to Rudolph KMS key operations"
  default = []
}

variable "enable_mutual_tls_authentication" {
  type = bool
  default = false
}
