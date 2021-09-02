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
}

variable "aws_account_id" {
  type    = string
}

variable "org" {
  type        = string
  description = "A lowercase string unique to your organization. Used to deduplicate certain AWS resources."
}

variable "stage_name" {
  type        = string
  description = "Name of stage to use for this deployment"
}

variable "prefix" {
  type        = string
  description = "Unique prefix to use for resources. This gets passed in from 'make deploy' command"
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
