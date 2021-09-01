resource "aws_api_gateway_model" "machine_config" {
  rest_api_id  = aws_api_gateway_rest_api.api_gateway.id
  name         = "MachineConfig"
  description  = "A configuration returned by /preflight"
  content_type = "application/json"

  schema = <<EOF
{
  "$schema" : "http://json-schema.org/draft-04/schema#",
  "title" : "MachineConfig schema",
  "type" : "object",
  "properties" : {
    "client_mode" : { "type" : "number" },
    "blocked_path_regex": { "type": "string" },
    "allowed_path_regex": { "type": "string" },
    "batch_size": { "type": "number" },
    "enable_bundles": { "type": "boolean" },
    "enable_transitive_rules": { "type": "boolean" },
    "clean_sync": { "type": "boolean" },
    "upload_logs_url": { "type": "string" }
  },
  "required": ["client_mode", "batch_size"]
}
EOF
}

resource "aws_api_gateway_model" "ruledownload" {
  rest_api_id  = aws_api_gateway_rest_api.api_gateway.id
  name         = "RuledownloadResponse"
  description  = "A response returned by /ruledownload"
  content_type = "application/json"

  schema = <<EOF
{
  "$schema" : "http://json-schema.org/draft-04/schema#",
  "title" : "Ruledownload Schema",
  "type" : "object",
  "properties" : {
    "rules": {
      "type" : "array",
      "items": {
        "type": "object",
        "properties": {
          "rule_type": { "type": "string" },
          "policy": { "type": "string" },
          "sha256": { "type": "string" },
          "custom_msg": { "type": "string" }
        },
        "required": ["rule_type", "policy", "sha256"]
      }
    },
    "cursor":  {
      "type": "object",
      "properties": {
        "is_last_page": { "type": "boolean" },
        "pk": { "type": "string" },
        "sk": { "type": "string" }
      }
    }
  }
}
EOF
}

resource "aws_api_gateway_model" "preflight_request" {
  rest_api_id  = aws_api_gateway_rest_api.api_gateway.id
  name         = "PreflightRequest"
  description  = "A request sent to /preflight"
  content_type = "application/json"

  schema = <<EOF
{
  "$schema" : "http://json-schema.org/draft-04/schema#",
  "title" : "PreflightRequest schema",
  "type" : "object",
  "properties" : {
    "os_build" : { "type" : "string" },
    "santa_version": { "type": "string" },
    "hostname": { "type": "string" },
    "os_version": { "type": "string" },
    "certificate_rule_count": { "type": "number" },
    "binary_rule_count": { "type": "number" },
    "client_mode": { "type": "string" },
    "serial_num": { "type": "string" },
    "primary_user": { "type": "string" },
    "compiler_rule_count": { "type": "number" },
    "transitive_rule_count": { "type": "number" },
    "request_clean_sync": { "type": "boolean" }
  },
  "required": ["santa_version", "serial_num", "primary_user"]
}
EOF
}
