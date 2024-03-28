locals {
  source_bucket_name       = "${var.prefix}-${var.org}-rudolph-events"
  s3_logging_bucket_name   = var.existing_logging_bucket_name != "" ? var.existing_logging_bucket_name : "${local.source_bucket_name}-logging"
  create_s3_logging_bucket = var.enable_logging && var.existing_logging_bucket_name == ""
}

#
# S3 Bucket for S3 logging, optional
#

resource "aws_s3_bucket" "s3_logging" {
  count = local.create_s3_logging_bucket ? 1 : 0

  bucket = local.s3_logging_bucket_name

  force_destroy = true

}

resource "aws_s3_bucket_policy" "s3_logging" {
  count = local.create_s3_logging_bucket ? 1 : 0

  bucket = aws_s3_bucket.s3_logging[0].id
  policy = format(
    data.aws_iam_policy_document.firehose_bucket_policy_template.json,
    local.s3_logging_bucket_name,
    local.s3_logging_bucket_name
  )
}

resource "aws_s3_bucket_versioning" "s3_logging" {
  count = local.create_s3_logging_bucket ? 1 : 0

  bucket = aws_s3_bucket.s3_logging[0].id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "s3_logging" {
  count = local.create_s3_logging_bucket ? 1 : 0

  bucket = aws_s3_bucket.s3_logging[0].id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_logging[0].key_id
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "s3_logging" {
  count = local.create_s3_logging_bucket ? 1 : 0

  bucket = aws_s3_bucket.s3_logging[0].id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "s3_logging" {
  depends_on = [aws_s3_bucket_ownership_controls.s3_logging]

  bucket = aws_s3_bucket.s3_logging[0].id
  acl    = "log-delivery-write"
}

#
# S3 Bucket for firehose
#

resource "aws_s3_bucket" "rudolph_eventsupload_firehose" {
  bucket = local.source_bucket_name

  force_destroy = true


}

resource "aws_s3_bucket_ownership_controls" "rudolph_eventsupload_firehose" {
  bucket = aws_s3_bucket.rudolph_eventsupload_firehose.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "rudolph_eventsupload_firehose" {
  depends_on = [aws_s3_bucket_ownership_controls.rudolph_eventsupload_firehose]

  bucket = aws_s3_bucket.rudolph_eventsupload_firehose.id
  acl    = "private"
}


resource "aws_s3_bucket_policy" "rudolph_eventsupload_firehose" {
  bucket = aws_s3_bucket.rudolph_eventsupload_firehose.id
  policy = format(
    data.aws_iam_policy_document.firehose_bucket_policy_template.json,
    local.source_bucket_name,
    local.source_bucket_name
  )
}

resource "aws_s3_bucket_versioning" "rudolph_eventsupload_firehose" {
  bucket = aws_s3_bucket.rudolph_eventsupload_firehose.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "rudolph_eventsupload_firehose" {
  bucket = aws_s3_bucket.rudolph_eventsupload_firehose.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.rudolph_eventsupload_kms_key.key_id
    }
  }
}

resource "aws_s3_bucket_logging" "rudolph_eventsupload_firehose" {
  count = var.enable_logging ? 1 : 0
  
  bucket = aws_s3_bucket.rudolph_eventsupload_firehose.id

  target_bucket = local.s3_logging_bucket_name
  target_prefix = "${local.source_bucket_name}/"
}
