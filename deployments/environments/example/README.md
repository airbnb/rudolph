# Setting Up New Environments
It is easy to deploy and maintain multiple different Rudolph environments, as they are configured in just 5 simple files.
This document will walk through the 5 files necessary.

## `init-env`
The easiest way to create skeletons for all necessary files is via:

```
make init-env
```

## The Directory
You will need a directory under `deployments/environments/{{ENV}}/`. For example, this current environment's `ENV` is "example", as it is located at `deployments/environment/example`.

### `_variables.tf`
This file is a symlink to [deployments/terraform_modules/default_variables.tf](deployments/terraform_modules/default_variables.tf).

You will never need to edit this variables file after creating the symlink. Leave it as-is.

### `main.tf`
This file is a symlink to [deployments/terraform_modules/default_main.tf](deployments/terraform_modules/default_main.tf).

You will never need to edit this main file after creating the symlink. Leave it as-is.

### `versions.tf`
Create a `versions.tf` in the environment directory with the appropriate versions. We recommend the following:
```
terraform {
  required_version = ">= 0.14.11"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.38.0"
    }
  }
}
```

Although you may use Terraform 1.0+. That'll probably be fine?

### `_backend.tf`
The first "complicated" file you'll need to set up is `_backend.tf` which describes how the Terraform state will be stored.

Consider the provided example:
```
terraform {
  backend "s3" {
    bucket     = "your-terraform-tfstate-bucket"
    key        = "rudolph-example-prefix.tfstate"
    region     = "us-east-1"
    acl        = "private"
    encrypt    = true
    kms_key_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    Name = "Rudolph"
  }
}
```

#### `bucket`
This must be a preexisting AWS S3 bucket in your AWS account. It is used for storing the Terraform statefile.
Sharing the same bucket across multiple deployments of Rudolph is OK.

#### `key`
This is the S3 bucket key in which your statefile is stored. It must be writable and must not be in use
by anything else, especially another Rudolph deployment.

We generally recommend including your `{{ENV}}` in some capacity in the key, to properly namespace it.

#### `region`
This is the AWS region into which you deploy Rudolph. Make sure this region is consistent everywhere else.

#### `acl`
Always put "private"!

#### `encrypt`
If you wish to encrypt your TF State with serverside encryption, set this to `true`.

#### `kms_key_id`
Only necessary if `encrypt = true`. Provide the KMS Key Id of an existing KMS key on the AWS account that
will be used for serverside encryption. The key must be usable by your IAM user.

### `config.auto.tfvars.json`
You **must** name this file exactly as specified above.

This file contains configurations that determine how your Rudolph deployment works.

Consider this provided example:
```
{
  "aws_account_id": "001122334455",
  "org": "acme",
  "prefix": "production",
  "stage_name": "prod",
  "region": "us-west-2",
  "eventupload_handler": "NONE",
  "eventupload_kinesis_name": "",
  "eventupload_autocreate_policies": false,
  "route53_zone_name": "mywebsite.com",
  "use_existing_route53_zone": true,
  "eventupload_output_lambda_name": ""
}
```

#### `aws_account_id`
Your AWS account's ID.

#### `org`
This string should be unique to your team or organization. This is used to prevent certain AWS resource Id collisions.
It is highly recomended to be universally unique.

#### `prefix`
This string prefixes all AWS resource names and identifiers so you can distinguish between different deployment
environments within the same AWS account. It must be unique within a given AWS Account, and it is highly recommended
to make this universally unique.

#### `stage_name`
This is used for your Lambda aliases and your API Gateway stage names. You can just use `prod`.

#### `region`
The AWS region, which must match the previous 2 aws regions.

#### `eventupload_handler`
It takes a value of `NONE`, `KINESIS`, or `FIREHOSE`. Use `"NONE"` if you do not intend to send events anywhere.

#### `eventupload_autocreate_policies`
If `true` Rudolph will automatically create sufficient IAM policies to the specified kinesis stream or firehose.

#### `route53_zone_name`
The domain name of the Route53 zone into which Rudolph will be deployed. You did remember to buy a domain name, right?

#### `use_existing_route53_zone`
Generally, the process of purchasing a domain name will automatically create a public hosted zone, so leave this as
`true`.

Specifying `false` will have Rudolph automatically create the Hosted Zone for you... but this can have issues later on
if you share the Hosted Zone between environments, as it becomes managed by Terraform. If you're not careful, Terraform
may accidentally destroy a Hosted Zone that's currently in use by another environment.

#### `eventupload_output_lambda_name`
An AWS Lambda Function name in the same region to forward `/eventupload` events to. Leave as `""` if you do not intend
to use this feature.
