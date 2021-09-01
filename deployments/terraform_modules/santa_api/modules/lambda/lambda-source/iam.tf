
// Default Bucket Policy to use as template; Enforces SSL access
data "aws_iam_policy_document" "default_bucket_policy_template" {
  # Force SSL access only
  statement {
    sid = "ForceSSLOnlyAccess"

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
