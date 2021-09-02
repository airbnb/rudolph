# Rudolph
Rudolph is the control server counterpart of [Santa](https://github.com/google/santa), and is used to rapidly deploy configurations to Santa agents.

Rudolph is built in Amazon Web Services, and utilizes exclusively serverless components to reduce operational burden. It is designed to be fast,
easy-to-use, low-maintenance, and cost-conscious.

# Getting Ready

## Golang
You will need golang version 1.15+ installed. Go get it from [the golang website](https://golang.org/dl/).

## Amazon Web Services
You will need an AWS account handy, as well as an IAM role with sufficient privileges.

## Terraform
You will need Terraform v0.14+ installed. We recommend using [tfenv](https://github.com/tfutils/tfenv).

## Registered Domain
You will need a registered DNS domain name for Rudolph's API server.

Purchasing a domain is very easy, simply use use AWS's Route 53 service to register a new domain. The exact URL
is rarely seen by end users of Santa-Rudolph, so we advise picking any random cheap one to start things up quickly.


# Initial Setup

Create a new directory here:

```
deployments/environments/{{ENV}}
```

Note down the value of `{{ENV}}`, this will be referred to as your "environment" and you will need this handy any time
deployments are made.

Under this directory, create two symlinks:

```
ln -s deployments/terraform_modules/default_main.tf deployments/environments/{{ENV}}/main.tf
ln -s deployments/terraform_modules/default_variables.tf deployments/environments/{{ENV}}/variables.tf
```

Create a configuration file named `config.auto.tfvars.json`.

Create a `versions.tf`

Create a `_backend.tf`

TBD
