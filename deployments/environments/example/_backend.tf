terraform {
  backend "s3" {
    # The S3 bucket in which the Terraform backend state is kept.
    #
    # This S3 bucket isn't created during a terraform apply and needs to be pre-existing.
    bucket = "your-terraform-tfstate-bucket"

    # The S3 key should be unique within the tfstate bucket
    key = "rudolph-example-prefix.tfstate"

    region  = "us-east-1"
    acl     = "private"
    encrypt = true

    # The KMS key Id used for SSE on the terraform S3 bucket, where applicable.
    kms_key_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    Name = "Rudolph"
  }
}
