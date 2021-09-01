# This is a convenience IAM policy allowing Rudolph to call kinesis:PutRecords to the desired firehose.
# Enable this using the var.eventupload_autocreate_policies variable.
# Alternatively, specify "false" and create the policies elsewhere.
locals {
  create_kinesis_policies = var.eventupload_handler == "KINESIS" && var.eventupload_autocreate_policies
}

data "aws_iam_policy_document" "kinesis_put_records" {
  count = local.create_kinesis_policies ? 1 : 0

  statement {
    actions = [
      "kinesis:PutRecords",
    ]

    resources = [
      "arn:aws:kinesis:${var.region}:${var.aws_account_id}:stream/${var.eventupload_kinesis_name}",
    ]
  }
}

resource "aws_iam_policy" "kinesis_put_records" {
  count = local.create_kinesis_policies ? 1 : 0

  name        = "${var.prefix}_rudolph_kinesis_put_records_policy"
  description = "Policy allowing Rudolph's eventupload endpoint to put records into Kinesis"
  policy      = data.aws_iam_policy_document.kinesis_put_records[0].json
}

resource "aws_iam_role_policy_attachment" "attach_kinesis_write" {
  count = local.create_kinesis_policies ? 1 : 0

  role       = module.eventupload_function.lambda_role_name
  policy_arn = aws_iam_policy.kinesis_put_records[0].arn
}