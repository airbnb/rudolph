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
First, you will need to [set up a new environment](deployments/environments/example/README.md).

## Terraform Files
Follow the instructions in the provided link and create all necessary Terraform configuration files. You will
end up creating:

* A new directory under `deployments/environments`
* A `main.tf` symlink
* A `variables.tf` symlink
* A `versions.tf` file
* A `_backend.tf` file
* A `config.auto.tfvars.json` file

## EXPORT your ENV
When working with a specific deployment environment, always make sure your `ENV` environment variable is set
properly.

For example, to direct commands at the `deployments/environment/example` environment, you would use:
```
export ENV=example
```

## Install Dependencies
Download and install all golang dependencies with:

```
make deps
```

## Build and Deploy
You can deploy your entire application now. Do not forget to authenticate with aws-cli.

```
make deploy
```

## Test it Out
Make an HTTP POST request to your new deployment. The request should hit a URL that looks something like:

```
POST /preflight/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE
Host: {{PREFIX}}-rudolph.{YOUR_DOMAIN_NAME}
content-type:application/json
accept:*/*
{
  "serial_num": "1234"
}
```


# Deploying & Configuring Santa Agents

