variable "gateway_rest_api_id" {
  type        = string
  description = "The resource id of the AWS REST API Gateway that this resource belongs to"
}

variable "gateway_resource_id" {
  type        = string
  description = "The resource id of the AWS API Gateway Resource that this resource belongs to"
}

variable "http_method" {
  type        = string
  description = "The HTTP method for this integration"
}

# This isn't the function arn or qualified function arn. It's an invocation ARN which is actually scoped to API Gate
variable "lambda_invocation_arn" {
  type        = string
  description = "The ARN of the API handler for the Lambda"
}

variable "authorizer_id" {
  type        = string
  description = "(OPTIONAL) Id of REST API Authorizer. Omit for NO authorization."
  default     = ""
}

variable "success_response_model" {
  type = string
  default = ""
}

variable "request_model" {
  type = string
  default = ""
}
