# Deployment
Ready to build Rudolph? Let's go at it!

---

## Getting Ready
You will only need to do this stuff once.

### Golang
You will need golang version 1.15+ installed. Go get it from [the golang website](https://golang.org/dl/).

### Amazon Web Services
You will need an AWS account handy, as well as an IAM role with sufficient privileges.

### Terraform
You will need Terraform v0.14+ installed. We recommend using [tfenv](https://github.com/tfutils/tfenv).

### Registered Domain
You will need a registered DNS domain name for Rudolph's API server.

Purchasing a domain is very easy, simply use use AWS's Route 53 service to register a new domain. The exact URL
is rarely seen by end users of Santa-Rudolph, so we advise picking any random cheap one to start things up quickly.

---

## Initial Setup
First, you will need to set up a new environment.

### Pick an ENV
Pick a new environment name for your deployment.

An environment is a unique instance of Rudolph is associated with a filesystem directory under [deployments/environments](deployments/environments/).

### EXPORT your ENV
When working with a specific deployment environment, always make sure your `ENV` environment variable is set
properly.

For example, to direct commands at the `deployments/environment/example` environment, you would use:

```
export ENV=example
```

### Create your Terraform Files
You will now create all necessary Terraform configuration files. For convenience, you can set up everything via:

```
make init-env
```

You will end up creating:

* A new directory under `deployments/environments`
* A `main.tf` symlink
* A `variables.tf` symlink
* A `versions.tf` file
* A `_backend.tf` file
* A `config.auto.tfvars.json` file
* (Optional) A `_outputs.tf` file

You will also need to edit some of these files. A good resource to understanding these configuration files
is the example ENV's [`README.md`](deployments/environments/example/README.md).


### Install Dependencies
Download and install all golang dependencies with:

```
make deps
```

### First Build and Deploy
You can deploy your entire application now. Luckily, this step is super easy! Do not forget to authenticate with aws-cli.

```
make deploy
```

### Test it Out
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

---

## Future Deployments
Deploying code and making any infrastructure changes is easy:

```
export ENV=YOURENV
make deploy
```
