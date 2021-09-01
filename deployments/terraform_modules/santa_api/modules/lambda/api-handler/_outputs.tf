output "lambda_function_arn" {
  value = aws_lambda_function.api_handler.arn
}

output "lambda_alias_invoke_arn" {
  value = aws_lambda_alias.api_handler.invoke_arn
}

output "lambda_invoke_arn" {
  value = aws_lambda_function.api_handler.invoke_arn
}

output "lambda_role_name" {
  value = aws_iam_role.api_handler_role.id
}
