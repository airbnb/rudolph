// KMS Key for rudolph_eventsupload-event-logging S3 server-side encryption
resource "aws_kms_key" "s3_logging" {
  count               = var.existing_logging_bucket_name == "" ? 1 : 0
  enable_key_rotation = true
  description         = "S3 Logging S3 Server-Side Encryption"
}

// KMS Alias for S3 server-side encryption
resource "aws_kms_alias" "s3_logging" {
  count         = var.existing_logging_bucket_name == "" ? 1 : 0
  name          = "alias/${local.s3_logging_bucket_name}-s3-sse"
  target_key_id = aws_kms_key.s3_logging[0].key_id
}

// KMS Key for rudolph_eventsupload S3 server-side encryption
resource "aws_kms_key" "rudolph_eventsupload_kms_key" {
  enable_key_rotation = true
  description         = "Rudolph EventsUpload S3 Server-Side Encryption"
  policy              = data.aws_iam_policy_document.rudolph_eventsupload_kms_key_policy.json

  tags = {
    Name = "Rudolph"
  }
}

data "aws_iam_policy_document" "rudolph_eventsupload_kms_key_policy" {
  statement {
    sid    = "Enable IAM User Permissions"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${var.aws_account_id}:root"]
    }

    actions   = ["kms:*"]
    resources = ["*"]
  }

  statement {
    sid    = "Allow principals in the account to use the key"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "StringEquals"
      variable = "kms:CallerAccount"
      values   = [var.aws_account_id]
    }

    # condition {
    #   test     = "StringLike"
    #   variable = "kms:ViaService"
    #   values   = ["firehose.*.amazonaws.com"]
    # }

    actions = [
      "kms:CreateGrant",
      "kms:Decrypt",
      "kms:DescribeKey",
      "kms:Encrypt",
      "kms:GenerateDataKey*",
      "kms:ReEncrypt*",
    ]

    resources = ["*"]
  }

  statement {
    sid       = "Allow access for Key Administrators"
    effect    = "Allow"
    actions   = [
      "kms:Create*",
      "kms:Describe*",
      "kms:Enable*",
      "kms:List*",
      "kms:Put*",
      "kms:Update*",
      "kms:Revoke*",
      "kms:Disable*",
      "kms:Get*",
      "kms:Delete*",
      "kms:TagResource",
      "kms:UntagResource",
      "kms:ScheduleKeyDeletion",
      "kms:CancelKeyDeletion",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]

    principals {
      type        = "AWS"
      identifiers = [
        "arn:aws:iam::${var.aws_account_id}:root",
      ]
    }
  }
}

// KMS Alias for S3 server-side encryption
resource "aws_kms_alias" "rudolph_eventsupload_kms_alias" {
  name          = "alias/${local.source_bucket_name}-s3-sse"
  target_key_id = aws_kms_key.rudolph_eventsupload_kms_key.key_id
}
