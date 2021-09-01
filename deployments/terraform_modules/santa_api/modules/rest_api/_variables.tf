variable "gateway_rest_api_id" {
  description = "The resource id of the AWS REST API Gateway that this resource belongs to"
  type        = string
}

variable "parent_resource_id" {
  description = "The resource id of the parent resource"
  type        = string
}

# If you want to include variables in the resource path, use {VAR_NAME}
variable "resource_path" {
  description = "The path for this resource"
  type        = string
}

variable "integration_http_methods" {
  description = "The HTTP methods that this method supports. Only supports GET/POST/PUT/PATCH/DELETE."
  type        = list(string)
  default     = ["GET"]
}

# This isn't the function arn or qualified function arn. It's an invocation ARN which is actually scoped to API Gate
variable "lambda_invocation_arn" {
  description = "The ARN of the API handler for the Lambda"
  type        = string
}

variable "authorizer_id" {
  description = "(OPTIONAL) Id of REST API Authorizer. Omit for NO authorization."
  type        = string
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
