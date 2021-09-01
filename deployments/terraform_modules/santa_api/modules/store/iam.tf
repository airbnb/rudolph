#
# Policy to be attached to specified Lambda function roles
#
data "aws_iam_policy_document" "lambda_permissions" {
  statement {
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:DescribeTable",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
    ]

    resources = [
      aws_dynamodb_table.store.arn,
      "arn:aws:dynamodb:${var.region}:${var.aws_account_id}:table/*_rudolph_store",
    ]
  }

  # Allow PutItem, but only on certain key prefixes, for the purposes
  # of uploading new config data
  statement {
    actions = [
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
    ]

    resources = [
      aws_dynamodb_table.store.arn,
      "arn:aws:dynamodb:${var.region}:${var.aws_account_id}:table/*_rudolph_store",
    ]

    # As an extra precaution, we prevent write operations to be conducted on GLOBAL rules and configurations
    # This prevents bugginess on the main API from propagating bad changes all around
    condition {
      test = "ForAllValues:StringLike"
      variable = "dynamodb:LeadingKeys"

      values = [
        "Machine#*",      # This needs to be consistent with the machineInfoPKPrefix constant
        "MachineInfo#*",  # This needs to be consistent with the machineInfoPKPrefix constant
        "MachineRules#*", # This needs to be consistent with the MachineRulesPKPrefix constant
      ]
    }
  }
}

resource "aws_iam_policy" "lambda_policy" {
  name   = "${var.prefix}_rudolph_store_policy"
  policy = data.aws_iam_policy_document.lambda_permissions.json
}

# Ideally the below resource would use for_each but
# terraform cannot use for_each with computed properties
resource "aws_iam_role_policy_attachment" "lambda_policy" {
  count      = length(var.read_lambda_role_names)
  role       = element(var.read_lambda_role_names, count.index)
  policy_arn = aws_iam_policy.lambda_policy.arn
}
