# This is a convenience IAM policy allowing Rudolph to call kinesis:PutRecords to the desired firehose.
# Enable this using the var.eventupload_autocreate_policies variable.
# Alternatively, specify "false" and create the policies elsewhere.
locals {
  create_lambda_policies = var.eventupload_output_lambda_name != "" && var.eventupload_autocreate_policies
}

data "aws_iam_policy_document" "lambda_eventupload_invoke_function" {
  count = local.create_lambda_policies ? 1 : 0

  statement {
    actions = [
      "lambda:InvokeFunction",
    ]

    resources = [
      "arn:aws:lambda:${var.region}:${var.aws_account_id}:function:${var.eventupload_output_lambda_name}",
      "arn:aws:lambda:${var.region}:${var.aws_account_id}:function:${var.eventupload_output_lambda_name}:*",
    ]
  }
}

resource "aws_iam_policy" "lambda_eventupload_invoke_function" {
  count = local.create_lambda_policies ? 1 : 0

  name        = "${var.prefix}_rudolph_eventupload_lambda_policy"
  description = "Policy allowing Rudolph's eventupload endpoint to invoke another Lambda"
  policy      = data.aws_iam_policy_document.lambda_eventupload_invoke_function[0].json
}

resource "aws_iam_role_policy_attachment" "attach_invoke_lambda" {
  count = local.create_lambda_policies ? 1 : 0

  role       = module.eventupload_function.lambda_role_name
  policy_arn = aws_iam_policy.lambda_eventupload_invoke_function[0].arn
}
