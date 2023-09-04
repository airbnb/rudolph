#
# Stuff down here is related to Route53 and ACM
#
resource "aws_acm_certificate" "api_ssl_certificate" {
  domain_name       = local.api_domain_name
  validation_method = "DNS"
}

resource "aws_acm_certificate_validation" "api_certificate_validation" {
  certificate_arn = aws_acm_certificate.api_ssl_certificate.arn

  validation_record_fqdns = [
    aws_route53_record.api_cert_validation_record.fqdn
  ]
}

# Import an existing Hosted Zone if desired
data "aws_route53_zone" "existing" {
  count = var.use_existing_route53_zone ? 1 : 0
  name  = var.route53_zone_name
#  zone_id = var.route53_zone_id
}

resource "aws_route53_zone" "new" {
  count = var.use_existing_route53_zone ? 0 : 1
  name  = var.route53_zone_name
}

/*
Hack for this error:

Error: Invalid index

  on terraform/modules/santa_api/route53.tf line 22, in resource "aws_route53_record" "api_cert_validation_record":
  22:   name    = aws_acm_certificate.api_ssl_certificate.domain_validation_options.0.resource_record_name

This value does not have any indices.
*/
locals {
  api_domain_name           = "${var.prefix}-rudolph.${var.route53_zone_name}"
  domain_validation_options = tomap(tolist(aws_acm_certificate.api_ssl_certificate.domain_validation_options)[0])
  route53_zone_id           = var.use_existing_route53_zone ? data.aws_route53_zone.existing[0].zone_id : aws_route53_zone.new[0].zone_id
}

resource "aws_route53_record" "api_cert_validation_record" {
  name    = lookup(local.domain_validation_options, "resource_record_name", "")
  type    = lookup(local.domain_validation_options, "resource_record_type", "")
  zone_id = local.route53_zone_id
  records = [lookup(local.domain_validation_options, "resource_record_value", "")]
  ttl     = 300
}

# @doc https://www.terraform.io/docs/providers/aws/r/api_gateway_domain_name.html
resource "aws_api_gateway_domain_name" "api_custom_domain" {
  domain_name              = local.api_domain_name
  regional_certificate_arn = aws_acm_certificate_validation.api_certificate_validation.certificate_arn
  security_policy          = "TLS_1_2"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

# Example DNS record using Route53.
# Route53 is not specifically required; any DNS host can be used.
resource "aws_route53_record" "rest_api_gateway_record" {
  zone_id = local.route53_zone_id
  name    = aws_api_gateway_domain_name.api_custom_domain.domain_name
  type    = "A"

  alias {
    evaluate_target_health = true
    name                   = aws_api_gateway_domain_name.api_custom_domain.regional_domain_name
    zone_id                = aws_api_gateway_domain_name.api_custom_domain.regional_zone_id
  }
}
