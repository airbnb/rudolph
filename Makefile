PKG = github.com/airbnb/rudolph
VERSION := $(shell git describe --tags --always)
DOCS_DIR ?= ./docs
DEPLOYMENT_ZIP_PATH = $(PWD)/build/package/deployment.zip
TERRAFORM_DEPLOYMENTS_DIR = $(PWD)/deployments/environments
TF_DEFAULT_FLAGS = --var zip_file_path="$(DEPLOYMENT_ZIP_PATH)" --var package_version=$(VERSION)
LDFLAGS=-ldflags="-X main.version=$(VERSION)"

# Check to ensure the prefix is being passed in as an arg like `ENV=<YOUR_ENVIRONMENT>`
# or is set using environment variables like `export ENV=<YOUR_ENVIRONMENT>`
#
.check-args:
ifndef ENV
	$(error ENV is not set. please use `ENV=<YOUR_ENV> make deploy` or `export ENV=<YOUR_ENV>.`)
endif

TERRAFORM_DEPLOYMENTS_ENV_DIR = $(TERRAFORM_DEPLOYMENTS_DIR)/$(ENV)

# Default target for `make`; Builds all Rudolph infrastructure and deploys the source code.
#
#
deploy: .check-args test build tf-init justdeploy

# Performs deploy without running tests or rebuilding the application. This is useful if you just
# recently already did a build but did not make any code changes, and don't want to repeat the whole
# process, but you should generally AVOID using this make target unless you know what you're doing.
justdeploy: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) apply $(TF_DEFAULT_FLAGS)

# Equivalent of `terraform plan` in the desired environment directory.
#
#
plan: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) plan $(TF_DEFAULT_FLAGS)

# Equivalent to a `terraform init` in the desired environment directory.
# Required upon first deployment and when the tf state drifts too far from the local cached state.
#
tf-init: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) init

# Destroys the current deployment
#
#
destroy: tf-init
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) destroy $(TF_DEFAULT_FLAGS)

# When the API Gateway deployment fails to re-deploy (for whatever reason), you can use this
# command to force a redeployment
#
force-api-redeploy: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) taint module.santa_api.aws_api_gateway_deployment.api_deployment
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) apply --target=module.santa_api.aws_api_gateway_deployment.api_deployment $(TF_DEFAULT_FLAGS)

# Install dependencies
#
#
deps:
	go mod download
	go mod tidy

# Compiles golang source into binaries
# We use linux amd64 binaries so AWS Lambda can run them
#
build: deps
	@sh -c "'$(CURDIR)/scripts/build.sh'"

# run all golang unit tests
#
#
test:
	go test -cover -v ./...

# Initializes the environment directory and files
#
#
init-env: .check-args
	@sh -c "'$(CURDIR)/scripts/new_env.sh' $(ENV)"

#
#
#
.PHONY: .check-args deploy tf-init plan destroy deps build test force-api-redeploy
