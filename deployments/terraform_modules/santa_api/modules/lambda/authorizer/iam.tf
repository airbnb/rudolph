#
# IAM
#

# The role that API Gateway assumes and uses to then invoke the lambda function
# This is NOT the same role as the role that the lambda function invokes AS
resource "aws_iam_role" "invocation_role" {
  name               = "${var.prefix}_rudolph_api_gateway_authorizer"
  path               = "/rudolph/"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  tags = {
    Name = "Rudolph"
  }
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "invocation_policy" {
  name   = "${var.prefix}_AllowApiGatewaytoInvokeAuthorizerLambda"
  role   = aws_iam_role.invocation_role.id
  policy = data.aws_iam_policy_document.invocation_policy.json
}

data "aws_iam_policy_document" "invocation_policy" {
  statement {
    effect  = "Allow"
    actions = ["lambda:InvokeFunction"]

    resources = [
      module.authorizer_function.lambda_function_arn,
    ]
  }
}
