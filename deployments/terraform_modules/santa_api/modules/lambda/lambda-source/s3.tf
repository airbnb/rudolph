#
# S3 Bucket for Lambda source
#

locals {
  source_bucket_name       = "${var.prefix}-${var.org}-rudolph-source"
  s3_logging_bucket_name   = var.existing_logging_bucket_name != "" ? var.existing_logging_bucket_name : "${local.source_bucket_name}-logging"
  create_s3_logging_bucket = var.enable_logging && var.existing_logging_bucket_name == ""
}

resource "aws_s3_bucket" "santa_api_source" {
  bucket = local.source_bucket_name
  acl    = "private"
  policy = format(
    data.aws_iam_policy_document.default_bucket_policy_template.json,
    local.source_bucket_name,
    local.source_bucket_name
  )

  force_destroy = true

  versioning {
    enabled = true
  }

  dynamic "logging" {
    for_each = var.enable_logging ? [1] : []
    content {
      target_bucket = local.s3_logging_bucket_name
      target_prefix = "${local.source_bucket_name}/"
    }
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
        kms_master_key_id = aws_kms_key.santa_api_source.key_id
      }
    }
  }

  # tags = {
  #   Name = "Rudolph"
  # }
}

// KMS Key for S3 server-side encryption
resource "aws_kms_key" "santa_api_source" {
  enable_key_rotation = true
  description         = "Rudolph Source S3 Server-Side Encryption"

  # tags = {
  #   Name = "Rudolph"
  # }
}

// KMS Alias for S3 server-side encryption
resource "aws_kms_alias" "santa_api_source" {
  name          = "alias/${local.source_bucket_name}-s3-sse"
  target_key_id = aws_kms_key.santa_api_source.key_id
}
