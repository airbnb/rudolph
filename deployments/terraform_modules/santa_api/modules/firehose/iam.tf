// Default firehose bucket policy to use as template
// Enforce Secure Transport
data "aws_iam_policy_document" "firehose_bucket_policy_template" {
  statement {
    sid = ""

    effect = "Deny"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions = ["s3:*"]

    resources = [
      "arn:aws:s3:::%s",
      "arn:aws:s3:::%s/*",
    ]

    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values   = ["false"]
    }
  }
}


data "aws_iam_policy_document" "firehose_assume_role_policy" {
  statement {
    sid = ""

    effect = "Allow"

    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["firehose.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "firehose_s3_policy" {
  statement {
    sid = ""
    effect = "Allow"

    actions = [
      "s3:AbortMultipartUpload",
      "s3:GetBucketLocation",
      "s3:GetObject",
      "s3:ListBucket",
      "s3:ListBucketMultipartUploads",
      "s3:PutObject",
    ]

    resources = [
      aws_s3_bucket.rudolph_eventsupload_firehose.arn,
      "${aws_s3_bucket.rudolph_eventsupload_firehose.arn}/*",
    ]
  }

  statement {
    sid = ""

    effect = "Allow"

    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]

    resources = [
      aws_kms_key.rudolph_eventsupload_kms_key.arn,
    ]

    condition {
      test     = "StringEquals"
      variable = "kms:ViaService"
      values   = ["s3.${var.region}.amazonaws.com"]
    }

  }
}

resource "aws_iam_policy" "eventsupload_firehose_s3_policy" {
  name   = "${var.prefix}_rudolph_eventupload_firehose_s3_policy"
  policy = data.aws_iam_policy_document.firehose_s3_policy.json
}

resource "aws_iam_role_policy_attachment" "eventsupload_firehose_s3_policy_attachment" {
  role = aws_iam_role.eventsupload_firehose_role.name
  policy_arn = aws_iam_policy.eventsupload_firehose_s3_policy.arn
}

resource "aws_cloudwatch_log_stream" "s3_delivery" {
  name           = "S3Delivery"
  log_group_name = aws_cloudwatch_log_group.eventsupload_firehose.name
}

data "aws_iam_policy_document" "firehose_cloudwatch_policy" {
  statement {
    effect = "Allow"

    actions = [
      "logs:DescribeLogStreams",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:${var.region}:${var.aws_account_id}:log-group:/aws/kinesisfirehose/${var.prefix}-rudolph-eventsupload-firehose:*",
    ]
  }
}

resource "aws_iam_policy" "eventsupload_firehose_cloudwatch_policy" {
  name   = "${var.prefix}_rudolph_eventsupload_firehose_cloudwatch_policy"
  policy = data.aws_iam_policy_document.firehose_cloudwatch_policy.json
}

resource "aws_iam_role_policy_attachment" "eventsupload_firehose_cloudwatch_policy_attachment" {
  role = aws_iam_role.eventsupload_firehose_role.name
  policy_arn = aws_iam_policy.eventsupload_firehose_cloudwatch_policy.arn
}

#
# Policy to be attached to specified Lambda function roles
#
data "aws_iam_policy_document" "rudolph_firehose_eventupload" {
  statement {
    actions = [
      "firehose:DeleteDeliveryStream",
      "firehose:PutRecord",
      "firehose:PutRecordBatch",
      "firehose:UpdateDestination"
    ]

    resources = [
      aws_kinesis_firehose_delivery_stream.eventsupload_firehose.arn,
    ]
  }
}

resource "aws_iam_policy" "rudolph_firehose_eventupload" {
  name   = "${var.prefix}_rudolph_firehose_eventupload_policy"
  policy = data.aws_iam_policy_document.rudolph_firehose_eventupload.json
}

resource "aws_iam_role_policy_attachment" "rudolph_firehose_eventupload" {
  count      = length(var.firehose_upload_lambda_role_names)
  role       = element(var.firehose_upload_lambda_role_names, count.index)
  policy_arn = aws_iam_policy.rudolph_firehose_eventupload.arn
}


resource "aws_iam_role" "eventsupload_firehose_role" {
  name   = "${var.prefix}_rudolph_eventsupload_firehose_role"
  path   = "/rudolph/"
  assume_role_policy = data.aws_iam_policy_document.firehose_assume_role_policy.json

  # tags = {
  #   Name = "Rudolph"
  # }
}
