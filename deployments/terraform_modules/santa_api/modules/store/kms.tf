// KMS Key for DynamoDB server-side encryption
resource "aws_kms_key" "store_sse_key" {
  enable_key_rotation = true
  description         = "Santa Rules Tables Server-Side Encryption"
  policy              = data.aws_iam_policy_document.store_sse_permissions.json
}

data "aws_iam_policy_document" "store_sse_permissions" {
  statement {
    sid       = "Enable IAM User Permissions"
    effect    = "Allow"
    actions   = ["kms:*"]
    resources = ["*"]
    principals {
      type        = "AWS"
      identifiers = [
        "arn:aws:iam::${var.aws_account_id}:root",
      ]
    }
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
      variable = "kms:ViaService"
      values   = ["dynamodb.${var.region}.amazonaws.com"]
    }
    condition {
      test = "StringEquals"
      variable = "kms:CallerAccount"
      values = [var.aws_account_id]
    }

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
    sid    = "Allow DynamoDB to get information about the CMK"
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["dynamodb.amazonaws.com"]
    }

    actions = [
      "kms:Describe*",
      "kms:Get*",
      "kms:List*",
    ]

    resources = ["*"]
  }

  dynamic "statement" {
    for_each = length(var.kms_key_administrators_arns) == 0 ? [] : [1]

    content {
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
        identifiers = var.kms_key_administrators_arns
      }
    }
  }
}

// KMS Alias for DynamoDB server-side encryption
resource "aws_kms_alias" "store_sse_key" {
  name          = "alias/${var.prefix}-santa-rules-tables-sse"
  target_key_id = aws_kms_key.store_sse_key.key_id
}
