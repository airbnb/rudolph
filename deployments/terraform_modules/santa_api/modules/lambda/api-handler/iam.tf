#
# IAM
#

# IAM Role for the API Handler Lambda
resource "aws_iam_role" "api_handler_role" {
  name               = "${var.prefix}_rudolph_${var.endpoint}_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_policy.json
  path               = "/rudolph/"

  tags = {
    Name = "Rudolph"
  }
}

data "aws_iam_policy_document" "lambda_execution_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

# Attach write permissions for CloudWatch logs
data "aws_iam_policy" "basic_execution_role" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "basic_execution_role" {
  role       = aws_iam_role.api_handler_role.id
  policy_arn = data.aws_iam_policy.basic_execution_role.arn
}
