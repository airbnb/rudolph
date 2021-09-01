PKG = github.com/airbnb/rudolph
VERSION := $(shell git describe --tags --always)
BUILD_DIR = build
LINUX_BUILD_DIR = $(BUILD_DIR)/linux
MACOS_BUILD_DIR = $(BUILD_DIR)/macos
CLI_NAME = rudolph
PKG_DIR = package
DOCS_DIR ?= ./docs
DEPLOYMENT_ZIP_PATH = $(PKG_DIR)/deployment.zip
TERRAFORM_DEPLOYMENTS_DIR = $(PWD)/deployments/environments
TF_DEFAULT_FLAGS = --var zip_file_path="$(PWD)/$(DEPLOYMENT_ZIP_PATH)" --var package_version=$(VERSION)
LDFLAGS=-ldflags="-X main.version=$(VERSION)"

# HANDLERS should be a list of the name of dirs in internal/endpoints/ to be built
# HANDLERS = authorizer health ruledownload preflight
#
# HANDLERS are stored under the defined HANDLERS_DIR folder
# HANDLERS_DIR contains all of the lambda handers
HANDLERS_DIR = internal/endpoints
HANDLERS = $(shell find $(HANDLERS_DIR) -type d -mindepth 1 -maxdepth 1 -exec basename {} \;)

# Check to ensure the prefix is being passed in as an arg like `ENV=<YOUR_ENVIRONMENT>`
# or is set using environment variables like `export ENV=<YOUR_ENVIRONMENT>`
.check-args:
ifndef ENV
	$(error ENV is not set. please use `make deploy ENV=<YOUR_PREFIX>. The ENV must correspond to a directory under terraform/deployments/`)
endif

TERRAFORM_DEPLOYMENTS_ENV_DIR = $(TERRAFORM_DEPLOYMENTS_DIR)/$(ENV)

# Default target for `make`
deploy: .check-args test build package tf-init justdeploy

# Performs deploy without running tests or rebuilding the application. This is useful if you just
# recently already did a build but did not make any code changes, and don't want to repeat the whole
# process, but you should generally AVOID using this make
justdeploy: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) apply $(TF_DEFAULT_FLAGS)

plan: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) plan $(TF_DEFAULT_FLAGS)

# Equivalent to a `terraform init` in the desired deployment directory.
# Required upon first deployment and when the tf state drifts too far from the local cached state.
tf-init: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) init

# Destroys the current deployment.
destroy: tf-init
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) destroy $(TF_DEFAULT_FLAGS)

# When the API Gateway deployment fails to re-deploy (for whatever reason), you can use this
# command to force a redeployment
force-api-redeploy: .check-args
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) taint module.santa_api.aws_api_gateway_deployment.api_deployment
	terraform -chdir=$(TERRAFORM_DEPLOYMENTS_ENV_DIR) apply --target=module.santa_api.aws_api_gateway_deployment.api_deployment $(TF_DEFAULT_FLAGS)

clean:
	rm -rf $(BUILD_DIR) $(PKG_DIR)

docs:
	rm -rf $(DOCS_DIR)
	mkdir -p $(DOCS_DIR)
	go run -ldflags="-X main.docsDir=$(DOCS_DIR)" ./cmd/docs
	cd $(DOCS_DIR) && ln -s $(CLI_NAME).md README.md

deps:
	go mod download
	go mod tidy

# Builds the source for all of the Lambda functions + the Rudolph CLI
build: rudolph
	$(info *** building endpoints)
	@for handler in $(HANDLERS); do \
		echo "building $$handler function: $(LINUX_BUILD_DIR)/$$handler" ; \
		GOOS=linux GOARCH=amd64 go build -o $(LINUX_BUILD_DIR)/$$handler ./internal/endpoints/$$handler ; \
	done

package: build
	mkdir $(PKG_DIR)
	zip $(DEPLOYMENT_ZIP_PATH) $(LINUX_BUILD_DIR)/*

#	find pkg -maxdepth 5 -type d -exec bash -c "cd '{}' && go test -cover -v -o /tmp/temp.test" \;
test:
	go test -cover -v ./...

.PHONY: .check-args deploy tf-init destroy clean docs deps build package test

# Builds the rudolph CLI binary
rudolph: clean deps
	$(info *** building rudolph CLI and symlinking to current directory)
	GOOS=darwin go build -v -o $(MACOS_BUILD_DIR)/$(CLI_NAME) $(LDFLAGS)
	ln -sf $(MACOS_BUILD_DIR)/$(CLI_NAME) ./$(CLI_NAME)

# Builds the rudolph CLI binary for Darwin platforms using Intel (amd64) and ARM (arm64) archs
rudolph_darwin_universal: clean deps
	$(info *** building rudolph CLI and symlinking to current directory)
	GOOS=darwin GOARCH=amd64 go build -v -o $(MACOS_BUILD_DIR)/$(CLI_NAME)_amd64 $(LDFLAGS)
	GOOS=darwin GOARCH=arm64 go build -v -o $(MACOS_BUILD_DIR)/$(CLI_NAME)_arm64 $(LDFLAGS)
	lipo -create -output $(MACOS_BUILD_DIR)/$(CLI_NAME) $(MACOS_BUILD_DIR)/$(CLI_NAME)_amd64 $(MACOS_BUILD_DIR)/$(CLI_NAME)_arm64
	ln -sf $(MACOS_BUILD_DIR)/$(CLI_NAME) ./$(CLI_NAME)
	echo ""
	echo "Darwin systems using ARM64 (M1+) archs must have codesigned binaries..."
	echo ""
	security find-identity -v -p codesigning | grep -i "developer id application"
	echo ""
	echo "use codesign -s {UUID} $(MACOS_BUILD_DIR)/$(CLI_NAME)"

