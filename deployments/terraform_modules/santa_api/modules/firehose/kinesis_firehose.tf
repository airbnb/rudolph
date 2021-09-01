locals {
  firehose_name     = var.eventupload_firehose_name == "" ? format("%s_rudolph_eventsupload_firehose", var.prefix) : var.eventupload_firehose_name
}

resource "aws_cloudwatch_log_group" "eventsupload_firehose" {
  name              = "/aws/kinesisfirehose/${local.firehose_name}"

  tags = {
    Name = "Rudolph"
  }
}

resource "aws_kinesis_firehose_delivery_stream" "eventsupload_firehose" {
  name        = local.firehose_name
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.eventsupload_firehose_role.arn
    bucket_arn = aws_s3_bucket.rudolph_eventsupload_firehose.arn
    kms_key_arn        = aws_kms_key.rudolph_eventsupload_kms_key.arn
    cloudwatch_logging_options {
      enabled         = true
      log_group_name  = aws_cloudwatch_log_group.eventsupload_firehose.name
      log_stream_name = "S3Delivery"
    }
    // TODO (@ryan) add logic here to tranform data for athena usage
    # processing_configuration {
    #   enabled = "true"

    #   processors {
    #     type = "Lambda"

    #     parameters {
    #       parameter_name  = "LambdaArn"
    #       parameter_value = "${module.firehose_processor.lambda_alias_invoke_arn.lambda_function_arn}:$LATEST"
    #     }
    #   }
    # }
  }
}
