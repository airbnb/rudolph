variable "prefix" {
  type        = string
  description = "Prefix to all resource names"
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
