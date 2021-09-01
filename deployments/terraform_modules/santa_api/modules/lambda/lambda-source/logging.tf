#
# S3 Bucket for S3 logging, optional
#

resource "aws_s3_bucket" "s3_logging" {
  count = var.existing_logging_bucket_name == "" ? 1 : 0

  bucket = local.s3_logging_bucket_name
  acl    = "log-delivery-write"
  policy = format(
    data.aws_iam_policy_document.default_bucket_policy_template.json,
    local.s3_logging_bucket_name,
    local.s3_logging_bucket_name
  )

  force_destroy = true

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
        kms_master_key_id = aws_kms_key.s3_logging[0].key_id
      }
    }
  }
}


// KMS Key for S3 server-side encryption
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
