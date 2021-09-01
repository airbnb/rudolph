#
# DynamoDB table to store rules
#
locals {
  dynamodb_table_name = var.dynamodb_table_name == "" ? "${var.prefix}_rudolph_store" : var.dynamodb_table_name
}

resource "aws_dynamodb_table" "store" {
  name         = local.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "PK"
  range_key = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "DataType"
    type = "S"
  }

  attribute {
    name = "SerialNum"
    type = "S"
  }

  ttl {
    attribute_name = "ExpiresAfter"
    enabled        = true
  }

  global_secondary_index {
    name               = "SerialNum_DataType_MachineID"
    hash_key           = "SerialNum"
    range_key          = "DataType"
    projection_type    = "INCLUDE"
    non_key_attributes = ["MachineID"]
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.store_sse_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Name = "Rudolph"
  }
}
