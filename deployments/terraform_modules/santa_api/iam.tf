#
# IAM Policy for the API gateway
#

data "aws_iam_policy_document" "api_resource_policy" {
  statement {
    effect  = "Allow"
    actions = ["execute-api:Invoke"]
    resources = [
      # The format of these ARNS take the form of
      # <RANDOM_KEY>/<STAGE>/<HTTP_METHOD>/<URLPATH>
      "arn:aws:execute-api:${var.region}:${var.aws_account_id}:*"
    ]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }

  statement {
    effect  = "Allow"
    actions = ["execute-api:Invoke"]
    resources = [
      # The format of these ARNS take the form of
      # <RANDOM_KEY>/<STAGE>/<HTTP_METHOD>/<URLPATH>
      "arn:aws:execute-api:${var.region}:${var.aws_account_id}:*"
    ]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "Bool"
      variable = "aws:ViaAWSService"
      values = [
        "true",
      ]
    }
  }

  # Resource statement to block all incoming requests that aren't from inside the allowed_inbound_ips
  #
  # For reference: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-control-access-policy-language-overview.html
  dynamic "statement" {
    for_each = length(var.allowed_inbound_ips) == 0 ? [] : [1]

    content {
      effect  = "Deny"
      actions = ["execute-api:Invoke"]
      # resources = [
      #   "arn:aws:execute-api:${var.region}:${var.aws_account_id}:*"
      # ]

      # Can use this method to block all access EXCEPT a public health endpoint.
      not_resources = [
        "arn:aws:execute-api:${var.region}:${var.aws_account_id}:*/${var.stage_name}/GET/health",
      ]

      principals {
        type        = "AWS"
        identifiers = ["*"]
      }

      condition {
        test     = "NotIpAddress"
        variable = "aws:SourceIp"
        values = var.allowed_inbound_ips
      }
    }
  }
}
