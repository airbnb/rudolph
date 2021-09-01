locals {
  source_bucket_name     = "${var.prefix}-rudolph-events"
  s3_logging_bucket_name = var.existing_logging_bucket_name != "" ? var.existing_logging_bucket_name : "${local.source_bucket_name}-logging"
}

#
# S3 Bucket for S3 logging, optional
#

resource "aws_s3_bucket" "s3_logging" {
  count = var.existing_logging_bucket_name == "" ? 1 : 0

  bucket = local.s3_logging_bucket_name
  acl    = "log-delivery-write"
  policy = format(
    data.aws_iam_policy_document.firehose_bucket_policy_template.json,
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

#
# S3 Bucket for firehose
#

resource "aws_s3_bucket" "rudolph_eventsupload_firehose" {
  bucket = local.source_bucket_name
  acl    = "private"
  policy = format(
    data.aws_iam_policy_document.firehose_bucket_policy_template.json,
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
        kms_master_key_id = aws_kms_key.rudolph_eventsupload_kms_key.key_id
      }
    }
  }

  tags = {
    Name = "Rudolph"
  }
}
